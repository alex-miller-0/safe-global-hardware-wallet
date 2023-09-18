package manager

import (
	"context"
	"flag"
	"fmt"

	"github.com/alex-miller-0/openpgp-secp256k1-wallet/pkg/ux"
	"github.com/alex-miller-0/safe-global-smartcard/internal/db"
	"github.com/alex-miller-0/safe-global-smartcard/internal/smartcard"
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
	signer, err := smartcard.Signer(s.Pin)
	if err != nil {
		ux.Errorf(err.Error())
		return subcommands.ExitFailure
	}
	safes := db.GetOwnedSafes(signer)
	str := fmt.Sprintf("\n-----\nConnected Smartcard: %s\n-----\n", signer)
	if len(safes) == 0 {
		str += "No owned safes found.\n"
	} else {
		str += "Owned safes:\n"
	}
	for _, safe := range safes {
		str += safe.String()
	}
	ux.Infoln(str)
	return subcommands.ExitSuccess
}
