package request

import (
	"fmt"
	"io"
	"net/http"

	"github.com/alex-miller-0/safe-global-smartcard/internal/networks"
)

func getStatus(safe, network string) ([]byte, error) {
	base, err := networks.BaseRequestURL(network)
	if err != nil {
		return nil, err
	}
	resp, err := http.Get(fmt.Sprintf("%s/v1/safes/%s", base, safe))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return io.ReadAll(resp.Body)
}

func getUnexecutedTxs(safe, network string, nonceGt int) ([]byte, error) {
	base, err := networks.BaseRequestURL(network)
	if err != nil {
		return nil, err
	}
	resp, err := http.Get(fmt.Sprintf(
		"%s/v1/safes/%s/multisig-transactions/?executed=false&nonce__gte=%d",
		base,
		safe,
		nonceGt,
	))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return io.ReadAll(resp.Body)
}
