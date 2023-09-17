package manager

import (
	"context"
	"flag"
	"fmt"

	"github.com/alex-miller-0/openpgp-secp256k1-wallet/pkg/ux"
	"github.com/alex-miller-0/safe-global-smartcard/internal/db"
	"github.com/google/subcommands"
)

const (
	ListDesc = "List all known address tags and safes in the db."
)

type List struct{}

func (*List) Name() string { return "list" }

func (*List) Synopsis() string {
	return ListDesc
}

func (*List) SetFlags(flagSet *flag.FlagSet) {}

func (*List) Usage() string {
	return "list\n"
}

func (*List) Execute(
	_ context.Context,
	flagSet *flag.FlagSet,
	_ ...any,
) subcommands.ExitStatus {
	tags := db.GetTags()
	str := "\n-----\nTags:\n"
	for _, t := range tags {
		str += fmt.Sprintf("  - %s: %s\n", t.Tag, t.Address)
	}
	str += "\n-----\nSafes:\n"
	safes := db.GetSafes()
	for _, s := range safes {
		str += fmt.Sprintf(
			"  - %s: %s\n    threshold: %d\n    owners: %v",
			s.ID.Tag,
			s.ID.Address,
			s.Threshold,
			s.Owners,
		)
	}
	ux.Infoln(str)
	return subcommands.ExitSuccess
}
