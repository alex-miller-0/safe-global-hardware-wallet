package request

type (
	SafeStatus struct {
		Address   string   `json:"address"`
		Nonce     int      `json:"nonce"`
		Threshold int      `json:"threshold"`
		Owners    []string `json:"owners"`
	}

	SafeTransactions struct {
		Count   int               `json:"count"`
		Results []SafeTransaction `json:"results"`
	}

	SafeTransaction struct {
		To                   string `json:"to"`
		Nonce                int    `json:"nonce"`
		GasPrice             string `json:"gasPrice"`
		MaxFeePerGas         string `json:"maxFeePerGas"`
		MaxPriorityFeePerGas string `json:"maxPriorityFeePerGas"`
		Value                string `json:"value"`
		Data                 string `json:"data"`

		Operation   int         `json:"operation"`
		DataDecoded DecodedData `json:"dataDecoded"`
		Signatures  string      `json:"signatures"`
	}

	DecodedData struct {
		Method string         `json:"method"`
		Params []DecodedParam `json:"parameters"`
		Tabs   int            `json:"-"`
	}

	DecodedParam struct {
		Name         string         `json:"name"`
		Type         string         `json:"type"`
		Value        any            `json:"value"`
		ValueDecoded []DecodedValue `json:"valueDecoded"`
		Tabs         int            `json:"-"`
	}

	DecodedValue struct {
		DataDecoded DecodedData `json:"dataDecoded"`
		Operation   int         `json:"operation"`
		To          string      `json:"to"`
		Type        string      `json:"type"`
		Value       string      `json:"value"`
	}
)
