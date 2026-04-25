package goeth9

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"gorm.io/gorm"
)

// 一个最小可运行的"迷你区块浏览器 / ERC-20 监听服务"示例：
// - 后台 goroutine 订阅指定 ERC-20 合约的 Transfer 事件
// - 将最近 N 条事件缓存在内存中
// - 通过 HTTP 接口 GET /events 返回最近事件列表
func Run(db *gorm.DB) {
	fmt.Println("\n open goeth9 success \n")

	rpcURL := os.Getenv("ETH_WS_URL")
	if rpcURL == "" {
		rpcURL = os.Getenv("ETH_RPC_URL")
	}
	if rpcURL == "" {
		log.Fatal("ETH_WS_URL or ETH_RPC_URL must be set")
	}

	contractHex := os.Getenv("ERC20_CONTRACT") // 监听的 ERC-20 合约地址
	if contractHex == "" {
		log.Fatal("ERC20_CONTRACT env is not set")
	}
	contractAddr := common.HexToAddress(contractHex)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	client, err := ethclient.DialContext(ctx, rpcURL)
	if err != nil {
		log.Fatalf("failed to connect to Ethereum node: %v", err)
	}
	defer client.Close()

	parsedABI, err := abi.JSON(strings.NewReader(erc20ABIJSON)) // 解析 ERC-20 ABI
	if err != nil {
		log.Fatalf("failed to parse ABI: %v", err)
	}

	eventStore := NewEventStore(100) // 创建一个事件存储，最多保存 100 条事件

	// 启动后台订阅协程
	go subscribeTransferEvents(ctx, client, parsedABI, contractAddr, eventStore)

	// HTTP 接口
	mux := http.NewServeMux()
	mux.HandleFunc("/events", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		events := eventStore.List()
		_ = json.NewEncoder(w).Encode(events) // 将事件列表编码为 JSON 并写入响应
	})

	server := &http.Server{
		Addr:         ":8080",
		Handler:      mux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	go func() {
		log.Printf("HTTP server listening on %s", server.Addr)
		err := server.ListenAndServe() // 监听 HTTP 请求
		if err != nil && err != http.ErrServerClosed {
			log.Fatalf("http server error: %v", err)
		}
	}()

	// 优雅退出
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	sig := <-sigCh
	fmt.Printf("received signal %s, shutting down...\n", sig.String())

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()
	_ = server.Shutdown(shutdownCtx) // 优雅关闭 HTTP 服务器
	cancel()

	fmt.Println("\n open goeth9 end \n")
}

// subscribeTransferEvents 订阅指定合约的 Transfer 事件，并将事件数据存储到 EventStore 中
func subscribeTransferEvents(ctx context.Context, client *ethclient.Client, parsedABI abi.ABI, contract common.Address, eventStore *EventStore) {
	query := ethereum.FilterQuery{
		Addresses: []common.Address{contract},
	}
	logsCh := make(chan types.Log) // 订阅日志的通道
	sub, err := client.SubscribeFilterLogs(ctx, query, logsCh)
	if err != nil {
		log.Fatalf("failed to subscribe logs: %v", err)
	}

	log.Printf("listening Transfer events of %s", contract.Hex())

	for {
		select {
		case vLog := <-logsCh:
			if len(vLog.Topics) == 0 {
				continue
			}

			// 解码事件
			var event struct {
				From  common.Address
				To    common.Address
				Value *big.Int
			}

			// indexed 地址从 Topics[1], Topics[2]
			if len(vLog.Topics) >= 3 {
				event.From = common.BytesToAddress(vLog.Topics[1].Bytes())
				event.To = common.BytesToAddress(vLog.Topics[2].Bytes())
			}
			// 非 indexed 参数从 Data 解码
			err := parsedABI.UnpackIntoInterface(&event, "Transfer", vLog.Data)
			if err != nil {
				log.Printf("failed to unpack log data: %v", err)
				continue
			}

			eventStore.Add(TransferEvent{
				BlockNumber: vLog.BlockNumber,
				TxHash:      vLog.TxHash.Hex(),
				From:        event.From.Hex(),
				To:          event.To.Hex(),
				Value:       event.Value.String(),
				Timestamp:   time.Now(), // 简化：使用当前时间；可扩展为查询区块时间
			})
		case err := <-sub.Err():
			log.Printf("subscription error: %v", err)
			return
		case <-ctx.Done():
			log.Println("context cancelled, stop subscription")
			return
		}
	}
}

const erc20ABIJSON = `[
  {
    "anonymous": false,
    "inputs": [
      {"indexed": true, "name": "from", "type": "address"},
      {"indexed": true, "name": "to", "type": "address"},
      {"indexed": false, "name": "value", "type": "uint256"}
    ],
    "name": "Transfer",
    "type": "event"
  }
]`

type TransferEvent struct {
	BlockNumber uint64    `json:"block_number"`
	TxHash      string    `json:"tx_hash"`
	From        string    `json:"from"`
	To          string    `json:"to"`
	Value       string    `json:"value"` // 原始 uint256 字符串
	Timestamp   time.Time `json:"timestamp"`
}

type EventStore struct {
	mu     sync.RWMutex
	events []TransferEvent
	limit  int
}

// NewEventStore 创建一个新的事件存储，限制最多保存 limit 条事件
func NewEventStore(limit int) *EventStore {
	return &EventStore{
		events: make([]TransferEvent, 0, limit),
		limit:  limit,
	}
}

// Add 添加一个新的事件，如果超过限制则丢弃最旧的一条
func (s *EventStore) Add(e TransferEvent) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if len(s.events) >= s.limit {
		// 简单环形缓冲：丢弃最旧一条
		s.events = s.events[1:]
	}
	s.events = append(s.events, e)
}

// List 返回当前所有事件的副本
func (s *EventStore) List() []TransferEvent {
	s.mu.RLock()
	defer s.mu.RUnlock()

	out := make([]TransferEvent, len(s.events))
	copy(out, s.events)
	return out
}

// http://localhost:8080/events
// [
// 	{
// 		"block_number": 2,
// 		"tx_hash": "0xd880b7504254bd7cb9b580a84cdc75e19e90d429b9ead5b7313bca6f1144add9",
// 		"from": "0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266",
// 		"to": "0x70997970C51812dc3A010C7d01b50e0d17dc79C8",
// 		"value": "2800000000000000000",
// 		"timestamp": "2026-04-26T06:49:00.3686614+08:00"
// 	}
// ]
