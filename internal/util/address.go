package util

import (
	"crypto/ecdsa"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/crypto/secp256k1"
)

func IsEthereumAddress(address string) bool {
	return common.IsHexAddress(address)
}

func GetEthereumAddress(pubBytes []byte) (string, error) {
	if len(pubBytes) != 65 {
		return "", fmt.Errorf("invalid public key from smartcard: not 65 bytes")
	}
	var x, y big.Int
	x.SetBytes(pubBytes[1:33])
	y.SetBytes(pubBytes[33:])
	pub := &ecdsa.PublicKey{X: &x, Y: &y, Curve: secp256k1.S256()}
	return crypto.PubkeyToAddress(*pub).String(), nil
}
