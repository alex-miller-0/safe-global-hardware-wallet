package tx

import (
	"encoding/hex"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common/math"
)

const (
	SafeGlobalMultiSignContract = "0x40A2aCCbd92BCA938b02010E17A5b8929b49130D"
	SafeGlobalMultiSignMethod   = "multiSend"
	SafeGlobalMultiSignMethodId = "0x8d80ff0a"

	Uint8   = "uint8"
	Address = "address"
	Uint256 = "uint256"
	Bytes   = "bytes"
)

// ValidateSafeGlobalCall returns nil if the transaction is a valid call to a
// Safe Global contract. These are treated differently by the API, uzing the
// `valueDecoded` field.
func ValidateSafeGlobalCall(s *SafeTransaction) (bool, error) {
	isMultisendContract := s.To == SafeGlobalMultiSignContract
	isMultisendMethod := s.DataDecoded.Method == SafeGlobalMultiSignMethod
	// Currently only the multisign contract is supported
	isSafeGlobalContract := isMultisendContract
	isSafeGlobalMethod := isMultisendMethod

	// Not a call to a Safe Global contract. No error.
	if !isSafeGlobalContract {
		return false, nil
	}

	// Call to a Safe Global contract, but not a supported method. Error.
	if isSafeGlobalContract && !isSafeGlobalMethod {
		return true, fmt.Errorf(
			"invalid method %s for contract %s",
			s.DataDecoded.Method,
			s.To,
		)
	}
	// Checks on individual Safe Global contracts
	if isMultisendContract && isMultisendMethod {
		for _, p := range s.DataDecoded.Params {
			// If there is calldata (`value`), but it has not been decoded, fail out
			if len(fmt.Sprintf("%s", p.Value)) > 2 && len(p.ValueDecoded) == 0 {
				return true, fmt.Errorf(
					"calldata is not decoded for multisend contract %s",
					s.To,
				)
			}
			// Check for recursive calls to Safe Global contracts.
			// For now, these are unsupported
			if len(p.ValueDecoded) > 0 {
				for _, v := range p.ValueDecoded {
					for _, p2 := range v.DataDecoded.Params {
						if p2.ValueDecoded != nil {
							return true, fmt.Errorf(
								"unsupported recursive call to multisend contract %s",
								s.To,
							)
						}
					}
				}
			}
			// Validate the data
			err := validateMultisend(s.Data, p.ValueDecoded)
			if err != nil {
				return true, fmt.Errorf("error validating multisend: %w", err)
			}
		}
		return true, nil
	}
	// We should not get here
	return false, fmt.Errorf("unknown error")
}

func validateMultisend(outerData string, valueDecoded []DecodedValue) error {
	// Start building up `transactions` byte array. This is packed according to
	// the source code of the Safe multisend contract and NOT ABI-packed (although
	// the nested transactions are themselves ABI packed)
	transactionsBytes := []byte{}
	// Serialize each inner transaction into a byte array according to the Safe
	// multisend contract's source code
	for i, v := range valueDecoded {
		// `operation` is the first byte
		enc := []byte{byte(v.Operation)}
		// Add `to`
		to, err := hex.DecodeString(v.To[2:])
		if err != nil {
			return fmt.Errorf("error decoding `to` for tx #%d: %w", i, err)
		}
		enc = append(enc, to...)
		// Add `value`
		n, ok := math.ParseBig256(v.Value)
		if !ok {
			return fmt.Errorf("error converting %s to big.Int", v.Value)
		}
		valueUint256 := math.U256(n)
		enc = append(enc, math.U256Bytes(valueUint256)...)
		// Add `data`
		// Encode the transaction's `data` field and ensure it matches the `value`
		data, err := v.DataDecoded.Encode()
		if err != nil {
			return fmt.Errorf("error encoding data for tx #%d: %v", i, err)
		}
		if "0x"+hex.EncodeToString(data) != v.Data {
			return fmt.Errorf(
				"⚠️  Multisend transaction %d did not validate",
				i,
			)
		}
		// Prefix with size of data
		dataSz := big.NewInt(int64(len(data)))
		enc = append(enc, math.U256Bytes(math.U256(dataSz))...)
		enc = append(enc, data...)
		// Add these bytes to the `transactions` byte array
		transactionsBytes = append(transactionsBytes, enc...)
	}
	// Build the go-ethereum ABI packer for `multisend`
	bytesType, err := abi.NewType(Bytes, "", nil)
	if err != nil {
		return fmt.Errorf("error creating abi type: %w", err)
	}
	multisendAbi := abi.ABI{Methods: map[string]abi.Method{
		"multiSend": {
			Name: "multiSend",
			Inputs: abi.Arguments{
				{Type: bytesType, Name: "transactions"},
			},
		},
	}}
	// ABI-encode the full transaction
	outerEnc, err := multisendAbi.Pack(
		SafeGlobalMultiSignMethod,
		transactionsBytes,
	)
	if err != nil {
		return fmt.Errorf("error packing multisend transactions: %w", err)
	}
	// Finally, validate that the ABI-encoded `transactions` field matches the
	// `data` field of the outer transaction
	if SafeGlobalMultiSignMethodId+hex.EncodeToString(outerEnc) != outerData {
		return fmt.Errorf("⚠️  Full multisend transaction did not validate")
	}

	return nil
}
