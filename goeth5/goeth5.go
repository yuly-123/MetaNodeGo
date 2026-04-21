package goeth5

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"gorm.io/gorm"
)

// 通过 SubscribeNewHead 订阅新区块头。
// 注意：大多数节点要求使用 WebSocket RPC，例如：ws://127.0.0.1:8546 或 wss://...
func Run(db *gorm.DB) {
	fmt.Println("\n open goeth5 success \n")

	rpcURL := os.Getenv("ETH_WS_URL")
	if rpcURL == "" {
		// 回退到 ETH_RPC_URL，便于在只配置了 HTTP 的环境中看到错误提示
		rpcURL = os.Getenv("ETH_RPC_URL")
	}
	if rpcURL == "" {
		log.Fatal("ETH_WS_URL or ETH_RPC_URL must be set")
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	client, err := ethclient.DialContext(ctx, rpcURL)
	if err != nil {
		log.Fatalf("failed to connect to Ethereum node: %v", err)
	}
	defer client.Close()

	// 订阅新区块头
	headers := make(chan *types.Header)
	sub, err := client.SubscribeNewHead(ctx, headers)
	if err != nil {
		log.Fatalf("failed to subscribe new heads: %v", err)
	}

	fmt.Printf("Subscribed to new blocks via %s\n", rpcURL)

	// 捕获 Ctrl+C 退出
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	for {
		select {
		case h := <-headers:
			if h == nil {
				continue
			}
			printHeaderInfo(h)
		case err := <-sub.Err():
			log.Printf("subscription error: %v", err)
			return
		case sig := <-sigCh:
			fmt.Printf("received signal %s, shutting down...\n", sig.String())
			return
		case <-ctx.Done():
			fmt.Println("context cancelled, exiting...")
			return
		}
	}

	fmt.Println("\n open goeth5 end \n")
}

func printHeaderInfo(h *types.Header) {
	fmt.Printf("Time       : [%s]\n", time.Now().Format(time.RFC3339))
	fmt.Printf("Hash       : [%s]\n", h.Hash().Hex())
	fmt.Printf("ParentHash : [%s]\n", h.ParentHash.Hex())
	fmt.Printf("UncleHash  : [%s]\n", h.UncleHash.Hex())
	fmt.Printf("Coinbase   : [%s]\n", h.Coinbase.Hex())
	fmt.Printf("Root       : [%s]\n", h.Root.Hex())
	fmt.Printf("TxHash     : [%s]\n", h.TxHash.Hex())
	fmt.Printf("ReceiptHash : [%s]\n", h.ReceiptHash.Hex())
	fmt.Printf("Bloom      : [%x]\n", h.Bloom.Bytes())
	fmt.Printf("Difficulty  : [%s]\n", h.Difficulty.String())
	fmt.Printf("Number      : [%d]\n", h.Number.Uint64())
	fmt.Printf("GasLimit    : [%d]\n", h.GasLimit)
	fmt.Printf("GasUsed     : [%d]\n", h.GasUsed)
	fmt.Printf("Time        : [%s]\n", time.Unix(int64(h.Time), 0))
	fmt.Printf("Extra       : [%x]\n", h.Extra)
	fmt.Printf("MixDigest   : [%s]\n", h.MixDigest.Hex())
	fmt.Printf("Nonce       : [%d]\n", h.Nonce)
	fmt.Printf("BaseFee     : [%s]\n", h.BaseFee.String())
	fmt.Printf("WithdrawalsHash : [%s]\n", h.WithdrawalsHash.Hex())
	fmt.Printf("BlobGasUsed : [%d]\n", h.BlobGasUsed)
	fmt.Printf("ExcessBlobGas : [%d]\n", h.ExcessBlobGas)
	fmt.Printf("ParentBeaconBlockRoot : [%s]\n", h.ParentBeaconRoot.Hex())
	fmt.Printf("RequestsHash : [%s]\n", h.RequestsHash.Hex())
	fmt.Printf("SlotNumber : [%d]\n", h.SlotNumber)
	fmt.Println("=======================================================")
}
