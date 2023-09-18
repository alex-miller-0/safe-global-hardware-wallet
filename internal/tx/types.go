package tx

type (
	SafeTransaction struct {
		To                   string `json:"to"`
		Nonce                int    `json:"nonce"`
		GasPrice             string `json:"gasPrice"`
		MaxFeePerGas         string `json:"maxFeePerGas"`
		MaxPriorityFeePerGas string `json:"maxPriorityFeePerGas"`
		Value                string `json:"value"`
		Data                 string `json:"data"`

		Safe           string      `json:"safe"`
		Operation      int         `json:"operation"`
		DataDecoded    DecodedData `json:"dataDecoded"`
		Signatures     string      `json:"signatures"`
		BaseGas        int         `json:"baseGas"`
		SafeTxGas      int         `json:"safeTxGas"`
		GasToken       string      `json:"gasToken"`
		RefundReceiver string      `json:"refundReceiver"`

		SafeTxHash string `json:"safeTxHash"`

		ChainId uint64 `json:"-"`
	}

	DecodedData struct {
		Method string         `json:"method"`
		Params []DecodedParam `json:"parameters"`
		Tabs   int            `json:"-"`
	}

	DecodedParam struct {
		Name  string `json:"name"`
		Type  string `json:"type"`
		Value any    `json:"value"`
		// The presence of ValueDecoded indicates that there is
		// a nested transaction
		ValueDecoded []DecodedValue `json:"valueDecoded"`
		Tabs         int            `json:"-"`
	}

	DecodedValue struct {
		Data        string      `json:"data"`
		DataDecoded DecodedData `json:"dataDecoded"`
		Operation   int         `json:"operation"`
		To          string      `json:"to"`
		Type        string      `json:"type"`
		Value       string      `json:"value"`
	}
)
