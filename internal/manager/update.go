package manager

import (
	"context"
	"flag"
	"fmt"

	"github.com/alex-miller-0/openpgp-secp256k1-wallet/pkg/ux"
	"github.com/alex-miller-0/safe-global-smartcard/internal/db"
	"github.com/alex-miller-0/safe-global-smartcard/internal/networks"
	"github.com/alex-miller-0/safe-global-smartcard/internal/request"
	"github.com/alex-miller-0/safe-global-smartcard/internal/util"
	"github.com/google/subcommands"
)

const (
	UpdateDesc = "Update any address tag in the db. If this is a safe, it will also " +
		"be updated in the db."
)

type Update struct {
	Safe    bool
	Address string
	Tag     string
	Network string
}

func (*Update) Name() string { return "update" }

func (*Update) Synopsis() string {
	return UpdateDesc
}

func (*Update) Usage() string {
	return "update [--safe] <address> <tag>\n"
}

func (u *Update) SetFlags(flagSet *flag.FlagSet) {
	flagSet.BoolVar(
		&u.Safe,
		"safe",
		false,
		"Whether or not this address is a Safe address",
	)
	flagSet.StringVar(
		&u.Network,
		"network",
		networks.Mainnet,
		"[Only used with --safe] The network on which this Safe exists",
	)
}

func (u *Update) Execute(
	_ context.Context,
	flagSet *flag.FlagSet,
	_ ...any,
) subcommands.ExitStatus {
	u.Address = flagSet.Arg(0)
	u.Tag = flagSet.Arg(1)
	if u.Address == "" || u.Tag == "" {
		fmt.Println(u.Usage())
		return subcommands.ExitFailure
	} else if !util.IsEthereumAddress(u.Address) {
		ux.Errorf("Not an Ethereum address: %s", u.Address)
		return subcommands.ExitFailure
	}
	if u.Safe {
		err := u.updateSafe()
		if err != nil {
			ux.Errorf("failed to add safe: %s", err.Error())
			return subcommands.ExitFailure
		}
	}
	err := u.updateTag()
	if err != nil {
		ux.Errorf("failed to add tag: %s", err.Error())
		return subcommands.ExitFailure
	}
	ux.Passf("Successfully added %s to db.\n", u.Address)
	return subcommands.ExitSuccess
}

func (u *Update) updateTag() error {
	err := db.UpdateTag(db.AddressTag{
		Address: u.Address,
		Tag:     u.Tag,
	})
	if err != nil {
		return err
	}
	err = db.Commit()
	if err != nil {
		return fmt.Errorf("could not commit db: %s", err.Error())
	}
	return nil
}

func (u *Update) updateSafe() error {
	if !networks.IsSupportedNetwork(u.Network) {
		return fmt.Errorf("network %s not supported", u.Network)
	}
	record := db.Safe{
		ID:      db.AddressTag{Address: u.Address, Tag: u.Tag},
		Network: u.Network,
	}

	status, err := request.GetStatus(&record)
	if err != nil {
		return err
	} else if status.Address == "" {
		return fmt.Errorf("no Safe found at address")
	}
	record.Threshold = status.Threshold
	record.Owners = status.Owners
	err = db.UpdateSafe(record)
	if err != nil {
		return err
	}
	err = db.Commit()
	if err != nil {
		return fmt.Errorf("could not commit db: %s", err.Error())
	}
	return nil
}
