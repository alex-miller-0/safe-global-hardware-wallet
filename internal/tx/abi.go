package tx

import (
	"encoding/hex"
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/math"
)

// Encode ABI-encodes a set of parameters into a byte array
// It does not support nested transactions
func (d *DecodedData) Encode() ([]byte, error) {
	// Make sure there are no nested transactions
	for _, p := range d.Params {
		if p.ValueDecoded != nil {
			return nil, fmt.Errorf("nested transactions are not supported")
		}
	}
	// go-ethereum's ABI encoding is just too complicated and nearly impossible
	// to use with arbitrary types, so we'll punt and call out to a TS bin that
	// serializes method calls into hex strings
	return nil, nil
}

func getGethTypeValue(t string, value any) (any, error) {
	// Remove array suffixes - they will be implicitly handled by the `value`
	// type handler
	t = stripArray(t)
	switch value.(type) {
	case string:
		v := value.(string)
		// Match the type and parse the value string
		switch {
		// address = common.Address
		case t == "address":
			return common.HexToAddress(v), nil
		// string = string
		case t == "string":
			return v, nil
		// int/uint = big.Int
		case strings.Contains(t, "int"):
			n, ok := math.ParseBig256(v)
			if !ok {
				return nil, fmt.Errorf("error converting %s to big.Int", value)
			}
			return handleNumber(t, n), nil
		// bytes = []byte
		case strings.Contains(t, "bytes"):
			if v[:2] == "0x" {
				v = v[2:]
			}
			data, err := hex.DecodeString(v)
			if err != nil {
				return nil, fmt.Errorf("error decoding %s: %w", v, err)
			}
			return data, nil
		default:
			return nil, fmt.Errorf("unsupported type %s (value=%s)", t, value)
		}
	default:
		return nil, fmt.Errorf("invalid value type %T", value)
	}
}

func stripArray(t string) string {
	if strings.Contains(t, "[") {
		return t[:strings.Index(t, "[")]
	}
	return t
}

func handleNumber(t string, n *big.Int) any {
	// For signed integers, parse as two's complement
	i := math.S256(n)
	if strings.Contains(t, "uint") {
		// Treate uints differently
		i = math.U256(n)
	}
	switch t {
	case "uint8":
		return uint8(i.Uint64())
	case "uint16":
		return uint16(i.Uint64())
	case "uint32":
		return uint32(i.Uint64())
	case "uint64":
		return uint64(i.Uint64())
	case "int8":
		return int8(i.Int64())
	case "int16":
		return int16(i.Int64())
	case "int32":
		return int32(i.Int64())
	case "int64":
		return int64(i.Int64())
	default:
		return i
	}

}
