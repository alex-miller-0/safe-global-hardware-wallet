package manager

import (
	"context"
	"flag"

	"github.com/alex-miller-0/safe-global-smartcard/internal/tx"
	"github.com/google/subcommands"
)

const (
	SignDesc = "Interactively sign transactions with the connected smartcard."
)

type Sign struct{}

func (*Sign) Name() string { return "sign" }

func (*Sign) Synopsis() string {
	return SignDesc
}

func (*Sign) SetFlags(flagSet *flag.FlagSet) {}

func (*Sign) Usage() string {
	return "sign\n"
}

func (*Sign) Execute(
	_ context.Context,
	flagSet *flag.FlagSet,
	_ ...any,
) subcommands.ExitStatus {
	tx := tx.SafeTransaction{}
	tx.Verify()
	return subcommands.ExitSuccess
}
