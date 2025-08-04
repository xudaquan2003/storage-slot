package withdraw

import (
	"bytes"
	"context"
	"errors"

	"github.com/ethereum-optimism/optimism/op-node/withdrawals"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
)

func TxSlot(ctx context.Context, l2c *rpc.Client, l2TxHash common.Hash) (string, error) {
	l2 := ethclient.NewClient(l2c)
	receipt, err := l2.TransactionReceipt(ctx, l2TxHash)
	if err != nil {
		return "", err
	}

	ev, err := withdrawals.ParseMessagePassed(receipt)
	if err != nil {
		return "", err
	}

	withdrawalHash, err := withdrawals.WithdrawalHash(ev)
	if !bytes.Equal(withdrawalHash[:], ev.WithdrawalHash[:]) {
		return "", errors.New("computed withdrawal hash incorrectly")
	}
	if err != nil {
		return "", err
	}
	println("withdrawal txHash:", withdrawalHash.Hex())
	slot := withdrawals.StorageSlotOfWithdrawalHash(withdrawalHash)

	return slot.Hex(), nil
}
