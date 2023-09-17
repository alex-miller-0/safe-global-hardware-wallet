package manager

import (
	"context"
	"flag"
	"fmt"

	"github.com/alex-miller-0/openpgp-secp256k1-wallet/pkg/ux"
	"github.com/alex-miller-0/safe-global-smartcard/internal/db"
	"github.com/alex-miller-0/safe-global-smartcard/internal/request"
	"github.com/google/subcommands"
)

const (
	StatusDesc = "Get status for a given Safe. " +
		"You may provide a Safe address or tag."
)

type Status struct {
	Network string
}

func (*Status) Name() string { return "status" }

func (*Status) Synopsis() string {
	return StatusDesc
}

func (*Status) Usage() string {
	return "status [--network <network>] <safe>\n"
}

func (s *Status) SetFlags(flagSet *flag.FlagSet) {
	flagSet.StringVar(
		&s.Network,
		"network",
		"ethereum",
		"The network where this Safe exists",
	)
}

func (s *Status) Execute(
	_ context.Context,
	flagSet *flag.FlagSet,
	_ ...any,
) subcommands.ExitStatus {
	if s.Network != "ethereum" {
		ux.Errorln("Only Ethereum is supported at this time.")
		return subcommands.ExitFailure
	}
	safeArg := flagSet.Arg(0)
	if safeArg == "" {
		fmt.Println(s.Usage())
		return subcommands.ExitFailure
	}
	safe := db.SearchSafe(safeArg)
	if safe == nil {
		ux.Errorf("no safe found: %s\n", safeArg)
		return subcommands.ExitFailure
	} else if safe.Network != s.Network {
		ux.Errorf("safe %s is not on network %s\n", safeArg, s.Network)
		return subcommands.ExitFailure
	}
	status, err := request.GetStatus(safe)
	if err != nil {
		ux.Errorln(err.Error())
		return subcommands.ExitFailure
	}
	fmt.Printf(
		"Safe Status:\nAddress: %s\nNonce: %d\nThreshold: %d\nOwners:\n",
		status.Address,
		status.Nonce,
		status.Threshold,
	)
	for _, owner := range status.Owners {
		fmt.Printf("  - %s\n", owner)
	}
	return subcommands.ExitSuccess
}
