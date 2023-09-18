package tx

import "fmt"

// Verify ABI-encodes the decoded data and makes sure it matches `data`
func (s *SafeTransaction) Verify() error {
	err := ValidateSafeGlobalCall(s)
	if err != nil {
		return fmt.Errorf("error validating Safe Global call: %w", err)
	}
	return nil
}
