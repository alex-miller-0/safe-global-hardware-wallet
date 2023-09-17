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
	TxsDesc = "Get transactions for a given Safe. " +
		"You may provide a Safe address or tag."
)

type Txs struct {
	Safe    string
	Pending bool
	Network string
}

func (*Txs) Name() string { return "txs" }

func (*Txs) Synopsis() string {
	return TxsDesc
}

func (*Txs) Usage() string {
	return "txs [--network <network>] <safe>\n"
}

func (t *Txs) SetFlags(flagSet *flag.FlagSet) {
	flagSet.StringVar(
		&t.Network,
		"network",
		"ethereum",
		"The network where this Safe exists",
	)
	flagSet.BoolVar(
		&t.Pending,
		"pending",
		true,
		"Return only pending transactions",
	)
}

func (t *Txs) Execute(
	_ context.Context,
	flagSet *flag.FlagSet,
	_ ...any,
) subcommands.ExitStatus {
	if t.Network != "ethereum" {
		ux.Errorln("Only Ethereum is supported at this time.")
		return subcommands.ExitFailure
	}
	safeArg := flagSet.Arg(0)
	if safeArg == "" {
		ux.Errorln("Must provide a Safe address or tag")
		return subcommands.ExitFailure
	}
	safe := db.SearchSafe(safeArg)
	if safe == nil {
		ux.Errorf("no safe found: %s\n", safeArg)
		return subcommands.ExitFailure
	} else if safe.Network != t.Network {
		ux.Errorf("safe %s is not on network %s\n", safeArg, t.Network)
		return subcommands.ExitFailure
	}
	txs, err := request.GetPendingTransactions(safe)
	if err != nil {
		ux.Errorln(err.Error())
		return subcommands.ExitFailure
	}
	str := fmt.Sprintf("\nFound %d pending transactions:\n", len(txs))
	for i, tx := range txs {
		str += fmt.Sprintf("\n--- [Transaction %d/%d] ---\n", i+1, len(txs))
		str += tx.String()
		str += "-------------------------"
	}
	ux.Infoln(str)

	return subcommands.ExitSuccess
}
