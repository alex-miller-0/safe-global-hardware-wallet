package tx

import (
	"encoding/hex"
	"fmt"

	"github.com/ethereum/go-ethereum/accounts/abi"
)

const (
	SafeGlobalMultiSignContract = "0x40A2aCCbd92BCA938b02010E17A5b8929b49130D"
	SafeGlobalMultiSignMethod   = "multiSend"

	Uint8   = "uint8"
	Address = "address"
	Uint256 = "uint256"
	Bytes   = "bytes"
)

// ValidateSafeGlobalCall returns nil if the transaction is a valid call to a
// Safe Global contract. These are treated differently by the API, uzing the
// `valueDecoded` field.
func ValidateSafeGlobalCall(s *SafeTransaction) error {
	isMultisendContract := s.To == SafeGlobalMultiSignContract
	isMultisendMethod := s.DataDecoded.Method == SafeGlobalMultiSignMethod
	// Currently only the multisign contract is supported
	isSafeGlobalContract := isMultisendContract
	isSafeGlobalMethod := isMultisendMethod

	// Not a call to a Safe Global contract. No error.
	if !isSafeGlobalContract {
		return nil
	}

	// Call to a Safe Global contract, but not a supported method. Error.
	if isSafeGlobalContract && !isSafeGlobalMethod {
		return fmt.Errorf(
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
				return fmt.Errorf(
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
							return fmt.Errorf(
								"unsupported recursive call to multisend contract %s",
								s.To,
							)
						}
					}
				}
			}
			// Validate the data
			err := validateMultisend([]byte(s.Data), p.ValueDecoded)
			if err != nil {
				return fmt.Errorf("error validating multisend: %w", err)
			}
		}
		return nil
	}

	return nil
}

func validateMultisend(outerData []byte, valueDecoded []DecodedValue) error {
	// Build the go-ethereum ABI packer for `multisend`
	bytesType, err := abi.NewType(Bytes, "", nil)
	if err != nil {
		return fmt.Errorf("error creating abi type: %w", err)
	}
	uint8Type, err := abi.NewType(Uint8, "", nil)
	if err != nil {
		return fmt.Errorf("error creating abi type: %w", err)
	}
	addressType, err := abi.NewType(Address, "", nil)
	if err != nil {
		return fmt.Errorf("error creating abi type: %w", err)
	}
	uint256Type, err := abi.NewType(Uint256, "", nil)
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
	// Each transaction will be encoded according to a fixed set of arguments
	multisendTxArgs := abi.Arguments{
		{Type: uint8Type, Name: "operation"},
		{Type: addressType, Name: "to"},
		{Type: uint256Type, Name: "value"},
		{Type: bytesType, Name: "data"},
	}
	// Start building up `transactions`
	transactionsBytes := []byte{}
	for i, v := range valueDecoded {
		// Encode the transaction's `data` field and ensure it matches the `value`
		data, err := v.DataDecoded.Encode()
		if err != nil {
			return fmt.Errorf("error encoding data for tx #%d: %w", i, err)
		}
		fmt.Printf("got data %x", data)
		if "0x"+hex.EncodeToString(data) != v.Value {
			return fmt.Errorf(
				"[⚠️] encoded inner transaction %d did not match expected: %w",
				i,
				err,
			)
		}
		// Encode the transaction's `operation`, `to`, and `value` fields
		// NOTE: `operation` must be an int to unmarsal API response, but we must
		// convert it to a string so that it will work with getGethTypeValue
		operation, err := getGethTypeValue(Uint8, fmt.Sprintf("%d", v.Operation))
		if err != nil {
			return fmt.Errorf("error converting operation for tx #%d: %w", i, err)
		}
		to, err := getGethTypeValue(Address, v.To)
		if err != nil {
			return fmt.Errorf("error converting to for tx #%d: %w", i, err)
		}
		value, err := getGethTypeValue(Uint256, v.Value)
		if err != nil {
			return fmt.Errorf("error converting value for tx #%d: %w", i, err)
		}
		// Pack the transaction's fields
		enc, err := multisendTxArgs.Pack(operation, to, value, data)
		if err != nil {
			return fmt.Errorf("error packing args for tx #%d: %w", i, err)
		}
		fmt.Printf("enc %x", enc)
		// Add these bytes to the `transactions` byte array
		transactionsBytes = append(transactionsBytes, enc...)
	}
	fmt.Printf("tx bytes %x", transactionsBytes)
	// ABI-encode `transactions` field
	outerEnc, err := multisendAbi.Pack(
		SafeGlobalMultiSignMethod,
		transactionsBytes,
	)
	if err != nil {
		return fmt.Errorf("error packing multisend transactions: %w", err)
	}
	fmt.Printf("outerEnc %x", outerEnc)
	// Finally, validate that the ABI-encoded `transactions` field matches the
	// `data` field of the outer transaction
	if "0x"+hex.EncodeToString(outerEnc) != "0x"+hex.EncodeToString(outerData) {
		return fmt.Errorf(
			"[⚠️] encoded transaction did not match expected: %w",
			err,
		)
	}

	return nil
}
