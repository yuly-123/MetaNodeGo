package goeth1

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"os"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"gorm.io/gorm"
)

// 连接以太坊节点，打印链 ID 和最新区块高度。
func Run(db *gorm.DB) {
	fmt.Println("\n open goeth1 success \n")

	// 构建 go eth client 连接
	rpcURL := os.Getenv("ETH_RPC_URL")
	if rpcURL == "" {
		log.Fatal("ETH_RPC_URL is not set")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second) // 超时设置为10秒，避免长时间等待连接失败的节点
	defer cancel()

	client, err := ethclient.DialContext(ctx, rpcURL)
	if err != nil {
		log.Fatalf("failed to connect to Ethereum node: %v", err)
	}
	defer client.Close()

	// 开启链上查询
	chainID, err := client.ChainID(ctx)
	if err != nil {
		log.Fatalf("failed to get chain id: %v", err)
	}

	fmt.Println("=== Ethereum Node Info ===")
	fmt.Printf("RPC URL       : %s\n", rpcURL)
	fmt.Printf("Chain ID      : %s\n", chainID.String())
	fmt.Println()
	fmt.Println("⚠️  注意: 'Latest' 区块是节点当前认为的最新区块，可能尚未被所有节点确认")
	fmt.Println("   不同RPC节点可能返回不同的 'latest' 区块，导致与浏览器不匹配")
	fmt.Println("   建议对比 'Safe' 或 'Finalized' 区块（已确认的区块）")
	fmt.Println()

	// 查询最新区块头
	latestHeader, err := client.HeaderByNumber(ctx, nil)
	if err != nil {
		log.Fatalf("failed to get latest block header: %v", err)
	} else {
		fmt.Println("=== Latest Block ===")
		fmt.Printf("Latest Block Number  : %d\n", latestHeader.Number.Uint64())
		fmt.Printf("Latest Block Hash    : %s\n", latestHeader.Hash().Hex())
		fmt.Printf("Latest Block Time    : %s\n", time.Unix(int64(latestHeader.Time), 0).Format(time.RFC3339))
		fmt.Println("==========================")
	}

	// 查询任意指定高度的区块头
	previousBlockNumber := new(big.Int).Sub(latestHeader.Number, big.NewInt(1)) // 获取前一个区块头
	previousHeader, err := client.HeaderByNumber(ctx, previousBlockNumber)
	if err != nil {
		log.Fatalf("failed to get previous block header: %v", err)
	} else {
		fmt.Println("=== Previous Block ===")
		fmt.Printf("Previous Block Number  : %d\n", previousHeader.Number.Uint64())
		fmt.Printf("Previous Block Hash    : %s\n", previousHeader.Hash().Hex())
		fmt.Printf("Previous Block Time    : %s\n", time.Unix(int64(previousHeader.Time), 0).Format(time.RFC3339))
		fmt.Println("==========================")
	}

	// 查询 safe 区块（浏览器通常显示这个）
	safeHeader, safeHash, err := getBlockByTag(ctx, client, "safe")
	if err != nil {
		log.Printf("failed to get safe block header: %v (this may not be supported by all nodes)", err)
	} else {
		fmt.Println("=== Safe Block (推荐对比) ===")
		fmt.Printf("Safe Block Number  : %d\n", safeHeader.Number.Uint64())
		fmt.Printf("Safe Block Hash    : %s (RPC提供的hash, 与浏览器一致)\n", safeHash.Hex())
		fmt.Printf("Safe Block Time    : %s\n", time.Unix(int64(safeHeader.Time), 0).Format(time.RFC3339))
		fmt.Printf("Calculated    : %s (计算出的hash, 可能不匹配)\n", safeHeader.Hash().Hex())
		fmt.Printf("Confirmations : %d\n", latestHeader.Number.Uint64()-safeHeader.Number.Uint64()) // 计算确认数（最新区块高度 - safe区块高度）
		fmt.Println("=============================")
	}

	// 查询 finalized 区块
	finalizedHeader, finalizedHash, err := getBlockByTag(ctx, client, "finalized")
	if err != nil {
		log.Printf("failed to get finalized block header: %v (this may not be supported by all nodes)", err)
	} else {
		fmt.Println("=== Finalized Block ===")
		fmt.Printf("Finalized Block Number  : %d\n", finalizedHeader.Number.Uint64())
		fmt.Printf("Finalized Block Hash    : %s (RPC提供的hash，与浏览器一致)\n", finalizedHash.Hex())
		fmt.Printf("Finalized Block Time    : %s\n", time.Unix(int64(finalizedHeader.Time), 0).Format(time.RFC3339))
		fmt.Printf("Calculated    : %s (计算出的hash, 可能不匹配)\n", finalizedHeader.Hash().Hex())
		fmt.Printf("Confirmations : %d\n", latestHeader.Number.Uint64()-finalizedHeader.Number.Uint64())
		fmt.Println("========================")
	}

	fmt.Println("\n open goeth1 end \n")
}

// getBlockByTag 查询指定标签的区块头（safe, finalized, latest 等）
// 返回 Header、RPC 提供的 Hash 和错误
// 注意：需要使用底层 RPC 调用，因为 ethclient 的高级 API 不直接支持这些标签
func getBlockByTag(ctx context.Context, client *ethclient.Client, tag string) (*types.Header, common.Hash, error) {
	// 获取底层 RPC 客户端
	rpcClient := client.Client()

	// 获取区块头数据（使用 false 只获取 header，不包含交易）
	var raw json.RawMessage
	err := rpcClient.CallContext(ctx, &raw, "eth_getBlockByNumber", tag, false)
	if err != nil {
		return nil, common.Hash{}, fmt.Errorf("RPC call failed: %w", err)
	}

	if len(raw) == 0 || string(raw) == "null" {
		return nil, common.Hash{}, fmt.Errorf("%s block not found", tag)
	}

	// 解析完整的区块头字段
	var blockData struct {
		Number      string         `json:"number"`
		Hash        common.Hash    `json:"hash"`
		ParentHash  common.Hash    `json:"parentHash"`
		UncleHash   common.Hash    `json:"sha3Uncles"`
		Coinbase    common.Address `json:"miner"`
		Root        common.Hash    `json:"stateRoot"`
		TxHash      common.Hash    `json:"transactionsRoot"`
		ReceiptHash common.Hash    `json:"receiptsRoot"`
		Bloom       hexutil.Bytes  `json:"logsBloom"`
		Difficulty  *hexutil.Big   `json:"difficulty"`
		GasLimit    hexutil.Uint64 `json:"gasLimit"`
		GasUsed     hexutil.Uint64 `json:"gasUsed"`
		Time        hexutil.Uint64 `json:"timestamp"`
		Extra       hexutil.Bytes  `json:"extraData"`
		MixDigest   common.Hash    `json:"mixHash"`
		Nonce       hexutil.Bytes  `json:"nonce"`
		BaseFee     *hexutil.Big   `json:"baseFeePerGas"`
	}
	if err := json.Unmarshal(raw, &blockData); err != nil {
		return nil, common.Hash{}, fmt.Errorf("failed to unmarshal block header: %w", err)
	}

	// 解析区块号
	num, ok := new(big.Int).SetString(blockData.Number[2:], 16)
	if !ok {
		return nil, common.Hash{}, fmt.Errorf("invalid block number: %s", blockData.Number)
	}

	// 构造完整的 Header
	header := &types.Header{
		ParentHash:  blockData.ParentHash,
		UncleHash:   blockData.UncleHash,
		Coinbase:    blockData.Coinbase,
		Root:        blockData.Root,
		TxHash:      blockData.TxHash,
		ReceiptHash: blockData.ReceiptHash,
		Bloom:       types.BytesToBloom(blockData.Bloom),
		Difficulty:  big.NewInt(0),
		Number:      num,
		GasLimit:    uint64(blockData.GasLimit),
		GasUsed:     uint64(blockData.GasUsed),
		Time:        uint64(blockData.Time),
		Extra:       blockData.Extra,
		MixDigest:   blockData.MixDigest,
		BaseFee:     nil,
	}

	// 设置 Difficulty
	if blockData.Difficulty != nil {
		header.Difficulty = blockData.Difficulty.ToInt()
	}

	// 设置 BaseFee（EIP-1559）
	if blockData.BaseFee != nil {
		header.BaseFee = blockData.BaseFee.ToInt()
	}

	// 设置 Nonce
	if len(blockData.Nonce) >= 8 {
		var nonceBytes [8]byte
		copy(nonceBytes[:], blockData.Nonce[:8])
		header.Nonce = types.BlockNonce(nonceBytes)
	}

	// 返回 Header 和 RPC 提供的 hash
	// 注意：手动构造的 Header 计算出的 hash 可能不准确，因为：
	// 1. RPC 返回的某些字段可能格式不完全匹配 go-ethereum 的内部格式
	// 2. Header 的内部缓存字段可能未正确初始化
	// 因此，我们应该直接使用 RPC 返回的 hash，它与浏览器显示的 hash 一致
	return header, blockData.Hash, nil
}
