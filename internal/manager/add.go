package manager

import (
	"context"
	"flag"
	"fmt"

	"github.com/alex-miller-0/openpgp-secp256k1-wallet/pkg/ux"
	"github.com/alex-miller-0/safe-global-smartcard/internal/db"
	"github.com/alex-miller-0/safe-global-smartcard/internal/request"
	"github.com/alex-miller-0/safe-global-smartcard/internal/util.go"
	"github.com/google/subcommands"
)

const (
	AddDesc = "Add any address tag to the db. If this is a safe, it will also " +
		"be added to the db."
)

type Add struct {
	Safe    bool
	Address string
	Tag     string
	Network string
}

func (*Add) Name() string { return "add" }

func (*Add) Synopsis() string {
	return AddDesc
}

func (*Add) Usage() string {
	return "add [--safe] <address> <tag>\n"
}

func (a *Add) SetFlags(flagSet *flag.FlagSet) {
	flagSet.BoolVar(
		&a.Safe,
		"safe",
		false,
		"Whether or not this address is a Safe address",
	)
	flagSet.StringVar(
		&a.Network,
		"network",
		"ethereum",
		"[Only used with --safe] The network on which this Safe exists",
	)
}

func (a *Add) Execute(
	_ context.Context,
	flagSet *flag.FlagSet,
	_ ...any,
) subcommands.ExitStatus {
	a.Address = flagSet.Arg(0)
	a.Tag = flagSet.Arg(1)
	if a.Address == "" || a.Tag == "" {
		fmt.Println(a.Usage())
		return subcommands.ExitFailure
	} else if !util.IsEthereumAddress(a.Address) {
		ux.Errorf("Not an Ethereum address: %s", a.Address)
		return subcommands.ExitFailure
	}
	if a.Safe {
		err := a.addSafe()
		if err != nil {
			ux.Errorf("failed to add safe: %s", err.Error())
			return subcommands.ExitFailure
		}
	}
	err := a.addTag()
	if err != nil {
		ux.Errorf("failed to add tag: %s", err.Error())
		return subcommands.ExitFailure
	}
	ux.Passf("Successfully added %s to db.\n", a.Address)
	return subcommands.ExitSuccess
}

func (a *Add) addTag() error {
	err := db.AddTag(db.AddressTag{Address: a.Address, Tag: a.Tag})
	if err != nil {
		return err
	}
	err = db.Commit()
	if err != nil {
		return fmt.Errorf("could not commit db: %s", err.Error())
	}
	return nil
}

func (a *Add) addSafe() error {
	if a.Network != "ethereum" {
		return fmt.Errorf("only Ethereum is supported at this time")
	}
	record := db.Safe{
		ID:      db.AddressTag{Address: a.Address, Tag: a.Tag},
		Network: a.Network,
	}
	status, err := request.GetStatus(&record)
	if err != nil {
		return err
	} else if status.Address == "" {
		return fmt.Errorf("no Safe found at address")
	}
	record.Threshold = status.Threshold
	record.Owners = status.Owners
	err = db.AddSafe(record)
	if err != nil {
		return err
	}
	err = db.Commit()
	if err != nil {
		return fmt.Errorf("could not commit db: %s", err.Error())
	}
	return nil
}
