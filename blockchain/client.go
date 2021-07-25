package blockchain

import (
	"context"
	"encoding/hex"
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/spf13/viper"
)

type QueryTransactionClientFn func(ctx context.Context, chainId int, txnHash string) (*TransactionInfo, bool, error)

func NewQueryTransactionClientFn(ethCli *ethclient.Client, bscCli *ethclient.Client) QueryTransactionClientFn {
	return func(ctx context.Context, chainId int, txnHash string) (*TransactionInfo, bool, error) {
		var cli *ethclient.Client
		switch chainId {
		case viper.GetInt("blockchain.ethereum.chainId"):
			cli = ethCli
		case viper.GetInt("blockchain.binance.chainId"):
			cli = bscCli
		default:
			cli = bscCli
		}

		tx, pending, err := cli.TransactionByHash(ctx, common.HexToHash(txnHash))
		if err != nil {
			return nil, pending, err
		}
		if pending {
			return nil, pending, nil
		}

		receipt, err := cli.TransactionReceipt(ctx, tx.Hash())
		if err != nil {
			return nil, pending, err
		}

		msg, err := tx.AsMessage(types.NewEIP155Signer(tx.ChainId()), tx.GasFeeCap())
		if err != nil {
			return nil, pending, err
		}

		block, err := cli.BlockByNumber(ctx, receipt.BlockNumber)
		if err != nil {
			return nil, pending, err
		}

		data := fmt.Sprintf("%x", tx.Data())
		abi, err := abi.JSON(strings.NewReader(bep20Abi))
		if err != nil {
			return nil, pending, err
		}
		decodeSig, err := hex.DecodeString(data[:8])
		if err != nil {
			return nil, pending, err
		}
		method, err := abi.MethodById(decodeSig)
		if err != nil {
			return nil, pending, err
		}
		decodeData, err := hex.DecodeString(data[8:])
		if err != nil {
			return nil, pending, err
		}

		params, err := method.Inputs.Unpack(decodeData)
		if err != nil {
			return nil, pending, err
		}

		// value, _ := weiToEther(tx.Value()).Float64()
		// gasPrice, _ := weiToEther(tx.GasPrice()).Float64()
		// txnFee, _ := weiToEther(new(big.Int).Mul(big.NewInt(int64(receipt.GasUsed)), tx.GasPrice())).Float64()
		// txnFee, _ := new(big.Float).SetInt(CalcGasCost(tx.Gas(), tx.GasPrice())).Float64()

		value, _ := ToDecimal(tx.Value(), 18).Float64()
		gasPrice, _ := ToDecimal(tx.GasPrice(), 18).Float64()
		txnFee, _ := ToDecimal(CalcGasCost(tx.Gas(), tx.GasPrice()), 18).Float64()
		amount, _ := ToDecimal(params[1].(*big.Int), 18).Float64()

		txnInfo := TransactionInfo{
			TxnHash:        tx.Hash().Hex(),
			Status:         int64(receipt.Status),
			Block:          receipt.BlockNumber.Int64(),
			Timestamp:      int64(block.Time()),
			From:           msg.From().Hex(),
			InteractedWith: tx.To().Hex(), // Token address
			TokenTransfer: TokenTransfer{
				From:   msg.From().Hex(),
				To:     params[0].(common.Address).Hex(),
				Amount: amount,
			},
			Value:    value,                  // BNB หน่วย ether
			TxnFee:   txnFee,                 // BNB
			GasPrice: gasPrice,               // BNB
			GasLimit: int64(tx.Gas()),        // amount
			GasUsed:  int64(receipt.GasUsed), // amount
			Nonce:    int64(tx.Nonce()),
		}
		return &txnInfo, pending, nil
	}
}
