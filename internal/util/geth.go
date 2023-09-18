package util

func ToByte32(b []byte) [32]byte {
	var a [32]byte
	copy(a[:], b)
	return a
}
