package manager

import (
	"context"
	"flag"
	"fmt"

	"github.com/alex-miller-0/openpgp-secp256k1-wallet/pkg/api"
	"github.com/alex-miller-0/openpgp-secp256k1-wallet/pkg/ux"
	"github.com/alex-miller-0/safe-global-smartcard/internal/db"
	"github.com/alex-miller-0/safe-global-smartcard/internal/util"
	"github.com/google/subcommands"
)

const (
	StatusDesc = "Print the Ethereum address on the current smartcard as well " +
		"as all saved safes for which it is an owner."
)

type Status struct {
	Pin string
}

func (*Status) Name() string { return "status" }

func (*Status) Synopsis() string {
	return StatusDesc
}

func (*Status) Usage() string {
	return "status\n"
}

func (s *Status) SetFlags(flagSet *flag.FlagSet) {
	flagSet.StringVar(
		&s.Pin,
		"pin",
		"",
		"The PIN for the smartcard",
	)
}

func (s *Status) Execute(
	_ context.Context,
	flagSet *flag.FlagSet,
	_ ...any,
) subcommands.ExitStatus {
	// Check on connected smartcard
	if s.Pin == "" {
		ux.PromptForSecret("Enter PIN: ", &s.Pin)
	}
	pubBytes, err := api.GetPub(s.Pin)
	if err != nil {
		ux.Errorf("error getting pubkey from smartcard: %s", err.Error())
		return subcommands.ExitFailure
	}
	address, err := util.GetEthereumAddress(pubBytes)
	if err != nil {
		ux.Errorf("error converting pubkey to address: %s", err.Error())
		return subcommands.ExitFailure
	}
	var ownedSafes []db.Safe
	safes := db.GetSafes()
	for _, safe := range safes {
		for _, owner := range safe.Owners {
			if owner == address {
				ownedSafes = append(ownedSafes, safe)
			}
		}
	}
	str := fmt.Sprintf("\n-----\nConnected Smartcard: %s\n-----\n", address)
	if len(ownedSafes) == 0 {
		str += "No owned safes found.\n"
	} else {
		str += "Owned safes:\n"
	}
	for _, safe := range ownedSafes {
		str += safe.String()
	}
	ux.Infoln(str)
	return subcommands.ExitSuccess
}
