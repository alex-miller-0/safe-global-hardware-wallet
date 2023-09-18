package tx

import (
	"encoding/hex"
	"fmt"
	"math/big"

	"github.com/alex-miller-0/safe-global-smartcard/internal/util"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/ethereum/go-ethereum/crypto"
)

// NOTE: This works for v1.3.0 of the Safe contracts, but it is not guaranteed
// to work for other versions.
// https://github.com/safe-global/safe-contracts/blob/v1.3.0/contracts/GnosisSafe.sol
const (
	DomainSeparatorType = "EIP712Domain(uint256 chainId,address verifyingContract)"
	SafeTxType          = "SafeTx(address to,uint256 value,bytes data," +
		"uint8 operation,uint256 safeTxGas,uint256 baseGas,uint256 gasPrice," +
		"address gasToken,address refundReceiver,uint256 nonce)"
)

var (
	EncHeader = []byte{0x19, 0x01}
)

func safeTxHash(tx *SafeTransaction) ([]byte, error) {
	args := abi.Arguments{
		{Type: abiType("bytes32"), Name: "SAFE_TX_TYPEHASH"},
		{Type: abiType("address"), Name: "to"},
		{Type: abiType("uint256"), Name: "value"},
		{Type: abiType("bytes32"), Name: "keccak256(data)"},
		{Type: abiType("uint8"), Name: "operation"},
		{Type: abiType("uint256"), Name: "safeTxGas"},
		{Type: abiType("uint256"), Name: "baseGas"},
		{Type: abiType("uint256"), Name: "gasPrice"},
		{Type: abiType("address"), Name: "gasToken"},
		{Type: abiType("address"), Name: "refundReceiver"},
		{Type: abiType("uint256"), Name: "_nonce"},
	}
	// Convert all fields into types that go-ethereum can handle
	safeTxTypehashVal := util.ToByte32(crypto.Keccak256([]byte(SafeTxType)))
	dataBytes, err := hex.DecodeString(tx.Data[2:])
	if err != nil {
		return nil, fmt.Errorf("error hex-decoding `data`: %w", err)
	}
	dataHashVal := util.ToByte32(crypto.Keccak256([]byte(dataBytes)))
	value, ok := math.ParseBig256(tx.Value)
	if !ok {
		return nil, fmt.Errorf("error converting %s to big.Int", tx.Value)
	}
	valueVal := math.U256(value)
	operationVal := uint8(tx.Operation)
	// NOTE: Normally I would parse these as big.Ints, but int64 is easier and
	// gas cannot realistically reach 2^31, as it is bounded to ~10M on mainnet
	safeTxGasVal := big.NewInt(int64(tx.SafeTxGas))
	baseGasVal := big.NewInt(int64(tx.BaseGas))
	gasPrice, ok := math.ParseBig256(tx.GasPrice)
	if !ok {
		return nil, fmt.Errorf("error converting %s to big.Int", tx.GasPrice)
	}
	gasPriceVal := math.U256(gasPrice)
	// Same story here - nonce is not goint to approach 2^31
	nonceVal := big.NewInt(int64(tx.Nonce))
	// ABI encode this data
	enc, err := args.Pack(
		safeTxTypehashVal,
		common.HexToAddress(tx.To),
		valueVal,
		dataHashVal,
		operationVal,
		safeTxGasVal,
		baseGasVal,
		gasPriceVal,
		common.HexToAddress(tx.GasToken),
		common.HexToAddress(tx.RefundReceiver),
		nonceVal,
	)
	if err != nil {
		return nil, err
	}
	// Hash and return
	return crypto.Keccak256(enc), nil
}

func abiType(t string) abi.Type {
	a, _ := abi.NewType(t, "", nil)
	return a
}
