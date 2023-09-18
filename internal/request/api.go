package request

import (
	"encoding/json"
	"fmt"

	"github.com/alex-miller-0/safe-global-smartcard/internal/db"
	"github.com/alex-miller-0/safe-global-smartcard/internal/tx"
)

func GetStatus(safe *db.Safe) (*SafeStatus, error) {
	resp, err := getStatus(safe.ID.Address, safe.Network)
	if err != nil {
		return nil, fmt.Errorf("error getting Safe status: %w", err)
	}
	status := &SafeStatus{}
	err = json.Unmarshal(resp, status)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling Safe status: %w", err)
	}
	if status.Address == "" {
		return nil, fmt.Errorf("no safe at %s", safe.ID.Address)
	}
	return status, nil
}

func GetPendingTransactions(safe *db.Safe) ([]tx.SafeTransaction, error) {
	status, err := GetStatus(safe)
	if err != nil {
		return nil, err
	}
	nonce := status.Nonce
	resp, err := getUnexecutedTxs(safe.ID.Address, safe.Network, nonce)
	if err != nil {
		return nil, fmt.Errorf("error getting pending transactions: %w", err)
	}
	txs := SafeTransactions{}
	err = json.Unmarshal(resp, &txs)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling pending transactions: %w", err)
	}
	return txs.Results, nil
}
