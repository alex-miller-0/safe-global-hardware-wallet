package request

import "github.com/alex-miller-0/safe-global-smartcard/internal/tx"

type (
	SafeStatus struct {
		Address   string   `json:"address"`
		Nonce     int      `json:"nonce"`
		Threshold int      `json:"threshold"`
		Owners    []string `json:"owners"`
	}

	SafeTransactions struct {
		Count   int                  `json:"count"`
		Results []tx.SafeTransaction `json:"results"`
	}
)
