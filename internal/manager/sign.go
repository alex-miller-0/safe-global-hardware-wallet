package manager

import (
	"context"
	"flag"
	"fmt"

	"github.com/alex-miller-0/openpgp-secp256k1-wallet/pkg/ux"
	"github.com/alex-miller-0/safe-global-smartcard/internal/db"
	"github.com/alex-miller-0/safe-global-smartcard/internal/networks"
	"github.com/alex-miller-0/safe-global-smartcard/internal/request"
	"github.com/alex-miller-0/safe-global-smartcard/internal/smartcard"
	"github.com/google/subcommands"
)

const (
	SignDesc = "Interactively sign pending transactions with the connected " +
		"smartcard."
)

type Sign struct {
	Safe    string
	Pin     string
	Network string
}

func (*Sign) Name() string { return "sign" }

func (*Sign) Synopsis() string {
	return SignDesc
}

func (*Sign) Usage() string {
	return "sign [--network <network>, --pin <pin> --safe <safe>]\n"
}

func (s *Sign) SetFlags(flagSet *flag.FlagSet) {
	flagSet.StringVar(
		&s.Network,
		"network",
		networks.Mainnet,
		"The network where this Safe exists",
	)
	flagSet.StringVar(
		&s.Pin,
		"pin",
		"",
		"The PIN for the smartcard",
	)
	flagSet.StringVar(
		&s.Safe,
		"safe",
		"",
		"The address or tag of the Safe to sign transactions for",
	)
}

func (s *Sign) Execute(
	_ context.Context,
	flagSet *flag.FlagSet,
	_ ...any,
) subcommands.ExitStatus {
	if !networks.IsSupportedNetwork(s.Network) {
		ux.Errorf("Unsupported network (%s)", s.Network)
		return subcommands.ExitFailure
	}
	signer, err := smartcard.Signer(s.Pin)
	if err != nil {
		ux.Errorf(err.Error())
		return subcommands.ExitFailure
	}
	ux.Infof("Using smartcard signer %s\n", signer)
	// safes := db.GetOwnedSafes(signer)

	safe := db.SearchSafe(s.Safe)
	txs, err := request.GetPendingTransactions(safe)
	if err != nil {
		ux.Errorln(err.Error())
		return subcommands.ExitFailure
	}
	for i, tx := range txs {
		var msg string
		err := tx.Verify()
		if err != nil {
			msg = fmt.Sprintf("Not verified: %s", err.Error())
		} else {
			msg = "✅  Verified"
		}
		str := fmt.Sprintf("Transaction %d/%d (%s)\n", i+1, len(txs), msg)
		str += tx.String()
		if err != nil {
			ux.Warnln(str)
		} else {
			ux.Passln(str)
		}
	}

	return subcommands.ExitSuccess
}
