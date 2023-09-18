package tx

import (
	"encoding/hex"
	"fmt"
	"math/big"

	"github.com/alex-miller-0/safe-global-smartcard/internal/db"
	"github.com/alex-miller-0/safe-global-smartcard/internal/util"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

// Verify ABI-encodes the decoded data and makes sure it matches `data`
func (s *SafeTransaction) Verify() error {
	isSafeGlobalCall, err := ValidateSafeGlobalCall(s)
	if err != nil {
		return fmt.Errorf("error validating Safe Global call: %w", err)
	}
	if !isSafeGlobalCall {
		enc, err := s.DataDecoded.Encode()
		if err != nil {
			return fmt.Errorf("error encoding data: %w", err)
		}
		if s.Data != "0x"+hex.EncodeToString(enc) {
			return fmt.Errorf("⚠️  Transaction data could not be verified")
		}
	}
	// Check the hash
	hash, err := s.Hash()
	if err != nil {
		return fmt.Errorf("error hashing transaction: %w", err)
	}
	if s.SafeTxHash != "0x"+hex.EncodeToString(hash) {
		return fmt.Errorf("⚠️  Transaction hash could not be verified")
	}
	return nil
}

// Hash produces the transaction hash to be signed. Note that `safe` must be
// passed, which
func (s *SafeTransaction) Hash() ([]byte, error) {
	// Get the domain separator
	domainSeparatorVal, err := s.domainSeparator()
	if err != nil {
		return nil, fmt.Errorf("could not get domain separator: %w", err)
	}
	// Get the safeTxHash
	safeTxHashVal, err := safeTxHash(s)
	if err != nil {
		return nil, fmt.Errorf("could not calculate safe TX hash: %w", err)
	}
	// Tightly pack and hash. There is no abi.encodePacked equivalent in
	// go-ethereum (of course), but the following has the same effect as
	// keccak256(abi.encodePacked(EncHeader, domainSeparatorVal, safeTxHashVal))
	return crypto.Keccak256(EncHeader, domainSeparatorVal, safeTxHashVal), nil
}

func (s *SafeTransaction) domainSeparator() ([]byte, error) {
	// We have to copy to a fixed sized array. Great example of why go-ethereum's
	// ABI functionality sucks
	dstHash := util.ToByte32(crypto.Keccak256([]byte(DomainSeparatorType)))
	chainId := big.NewInt(0).SetUint64(s.ChainId)
	verifyingContract := common.HexToAddress(s.Safe)
	args := abi.Arguments{
		{Type: abiType("bytes32"), Name: "sep"},
		{Type: abiType("uint256"), Name: "chainId"},
		{Type: abiType("address"), Name: "verifyingContract"},
	}
	enc, err := args.Pack(dstHash, chainId, verifyingContract)
	if err != nil {
		return nil, err
	}
	return crypto.Keccak256(enc), nil
}

func (s *SafeTransaction) String() string {
	str := fmt.Sprintf(
		"To: %s\nNonce: %d\nValue: %s\n",
		db.SwapAddress(s.To),
		s.Nonce,
		s.Value,
	)
	if s.DataDecoded.Method != "" {
		str += "Data:\n"
		str += s.DataDecoded.String()
	}
	return str
}
