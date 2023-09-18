package tx

import (
	"encoding/hex"
	"fmt"

	"github.com/alex-miller-0/safe-global-smartcard/internal/db"
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
	return nil
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
