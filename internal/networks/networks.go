package networks

import "fmt"

const (
	Mainnet = "ethereum"
)

func BaseRequestURL(network string) (string, error) {
	switch network {
	case Mainnet:
		return "https://safe-transaction-mainnet.safe.global/api/", nil
	default:
		return "", fmt.Errorf("network %s not supported", network)
	}
}

func IsSupportedNetwork(network string) bool {
	switch network {
	case Mainnet:
		return true
	default:
		return false
	}
}

func ChainIDFromNetwork(network string) (uint64, error) {
	switch network {
	case Mainnet:
		return 1, nil
	default:
		return 0, fmt.Errorf("network %s not supported", network)
	}
}
