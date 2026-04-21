package goeth4

import (
	"context"
	"flag"
	"fmt"
	"log"
	"math"
	"math/big"
	"os"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"gorm.io/gorm"
)

// 查询账户 ETH 余额（Wei 与 ETH）。
// 可以查最新余额，也可以查历史区块的余额（需要区块高度）。
// go run main.go --address 0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266 --block 0
func Run(db *gorm.DB) {
	fmt.Println("\n open goeth4 success \n")

	addrHex := flag.String("address", "", "account address (required)")
	blockNumber := flag.Int64("block", -1, "block number to query (-1 means latest)")
	flag.Parse()

	if *addrHex == "" {
		log.Fatal("missing --address flag")
	}

	rpcURL := os.Getenv("ETH_RPC_URL")
	if rpcURL == "" {
		log.Fatal("ETH_RPC_URL is not set")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	client, err := ethclient.DialContext(ctx, rpcURL)
	if err != nil {
		log.Fatalf("failed to connect to Ethereum node: %v", err)
	}
	defer client.Close()

	address := common.HexToAddress(*addrHex)

	var blockNum *big.Int
	if *blockNumber >= 0 {
		blockNum = big.NewInt(*blockNumber)
	}

	// 传nil就是查最新的，传区块高度就是查历史区块的余额
	balanceWei, err := client.BalanceAt(ctx, address, blockNum)
	if err != nil {
		log.Fatalf("failed to get balance: %v", err)
	}

	fmt.Println("=== Account Balance ===")
	fmt.Printf("Address     : %s\n", address.Hex())
	if blockNum == nil {
		fmt.Printf("Block       : latest\n")
	} else {
		fmt.Printf("Block       : %d\n", blockNum.Uint64())
	}
	fmt.Printf("Balance Wei : %s\n", balanceWei.String())

	ethValue := new(big.Float).Quo(
		new(big.Float).SetInt(balanceWei),
		big.NewFloat(math.Pow10(18)),
	)
	fmt.Printf("Balance ETH : %s\n", ethValue.Text('f', 6))

	fmt.Println("\n open goeth4 end \n")
}
