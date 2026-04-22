package goeth7

import (
	"context"
	"fmt"
	"log"
	"math"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"gorm.io/gorm"
)

// 展示订阅断线后的简单重连策略（示意实现）。
func Run(db *gorm.DB) {
	fmt.Println("\n open goeth7 success \n")

	rpcURL := os.Getenv("ETH_WS_URL")
	if rpcURL == "" {
		rpcURL = os.Getenv("ETH_RPC_URL")
	}
	if rpcURL == "" {
		log.Fatal("ETH_WS_URL or ETH_RPC_URL must be set")
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// // 捕获 Ctrl+C 退出，监听系统信号以优雅关闭
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-sigCh
		fmt.Printf("received signal %s, shutting down...\n", sig.String())
		cancel()
	}()

	runWithReconnect(ctx, rpcURL)

	fmt.Println("\n open goeth7 end \n")
}

func runWithReconnect(ctx context.Context, rpcURL string) {
	var attempt int

	for {
		select {
		case <-ctx.Done():
			fmt.Println("context cancelled, stop reconnect loop")
			return
		default:
		}

		attempt++
		log.Printf("connect attempt #%d to %s", attempt, rpcURL)

		client, err := ethclient.DialContext(ctx, rpcURL) // 重新连接以太坊节点
		if err != nil {
			log.Printf("failed to connect: %v", err)
			sleepWithBackoff(ctx, attempt)
			continue
		}

		headers := make(chan *types.Header)
		sub, err := client.SubscribeNewHead(ctx, headers)
		if err != nil {
			log.Printf("failed to subscribe new heads: %v", err)
			client.Close()
			sleepWithBackoff(ctx, attempt) // 订阅失败，等待后重试
			continue
		}

		log.Println("subscription established")

		// 订阅循环：如果 sub.Err() 返回错误，则跳出重新连接
		for {
			select {
			case h := <-headers:
				if h == nil {
					continue
				}
				fmt.Printf("New Block: %d, Hash: %s\n", h.Number.Uint64(), h.Hash().Hex())
			case err := <-sub.Err():
				log.Printf("subscription error: %v", err)
				client.Close()
				sleepWithBackoff(ctx, attempt)
				goto RECONNECT // 订阅出错，跳出循环重连
			case <-ctx.Done():
				log.Println("context cancelled, closing client")
				client.Close()
				return
			}
		}

	RECONNECT:
		// 下一轮 for 循环将尝试重连
	}
}

func sleepWithBackoff(ctx context.Context, attempt int) {
	// 简单指数退避，最大 1 分钟，time.Second 表示 1 秒，值 = 1e9 纳秒
	sec := int(math.Min(60, math.Pow(2, float64(attempt))))
	d := time.Duration(sec) * time.Second // 计算下一次重试的等待时间，单位为纳秒
	log.Printf("will retry in %s", d)

	t := time.NewTimer(d)
	defer t.Stop()

	select {
	case <-t.C:
	case <-ctx.Done():
	}
}
