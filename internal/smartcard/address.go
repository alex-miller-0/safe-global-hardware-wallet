package smartcard

import (
	"fmt"

	"github.com/alex-miller-0/openpgp-secp256k1-wallet/pkg/api"
	"github.com/alex-miller-0/openpgp-secp256k1-wallet/pkg/ux"
	"github.com/alex-miller-0/safe-global-smartcard/internal/util"
)

func Signer(pin string) (string, error) {
	if pin == "" {
		ux.PromptForSecret("Enter PIN: ", &pin)
	}
	pubBytes, err := api.GetPub(pin)
	if err != nil {
		return "", fmt.Errorf("error getting smartcard pubkey: %s", err.Error())
	}
	address, err := util.GetEthereumAddress(pubBytes)
	if err != nil {
		return "", fmt.Errorf("error converting pubkey to address: %s", err.Error())
	}
	return address, nil
}
