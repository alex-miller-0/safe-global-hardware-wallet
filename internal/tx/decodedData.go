package tx

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/alex-miller-0/safe-global-smartcard/internal/util"
)

const (
	EncoderBin = "encoder"
)

// Encode ABI-encodes a set of parameters into a byte array
// It does not support nested transactions
func (d *DecodedData) Encode() ([]byte, error) {
	// Make sure there are no nested transactions
	for _, p := range d.Params {
		if p.ValueDecoded != nil {
			return nil, fmt.Errorf("nested transactions are not supported")
		}
	}
	// go-ethereum's ABI encoding is just too complicated and nearly impossible
	// to use with arbitrary types, so we'll punt and call out to a JS bin that
	// serializes method calls into hex strings
	arg, err := json.Marshal(d)
	if err != nil {
		return nil, fmt.Errorf("error marshaling params: %w", err)
	}
	// Call out to the JS bin
	dir, err := os.Executable()
	if err != nil {
		return nil, fmt.Errorf("error getting executable path: %w", err)
	}
	path := filepath.Dir(dir) + "/" + EncoderBin
	stdout, err := exec.Command(path, string(arg)).Output()
	if err != nil {
		return nil, fmt.Errorf("error calling encoder: %w", err)
	}
	if len(stdout) == 0 {
		return nil, fmt.Errorf("invalid encoder output: %s", stdout)
	}
	return hex.DecodeString(strings.TrimSpace(string(stdout)))
}

func (d *DecodedData) String() string {
	str := ""
	if d.Method == "" {
		return ""
	}
	if d.Tabs == 0 {
		d.Tabs = 1
	}
	str += fmt.Sprintf("%s{%s}\n", util.PrintTabs(d.Tabs), d.Method)
	for _, p := range d.Params {
		p.Tabs = d.Tabs + 1
		str += p.String()
	}
	return str
}
