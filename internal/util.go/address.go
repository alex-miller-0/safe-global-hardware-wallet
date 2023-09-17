package util

import "github.com/ethereum/go-ethereum/common"

func IsEthereumAddress(address string) bool {
	return common.IsHexAddress(address)
}
