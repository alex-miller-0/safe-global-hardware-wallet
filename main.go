package main

import (
	"context"
	"flag"
	"os"

	"github.com/alex-miller-0/safe-global-smartcard/internal/manager"
	"github.com/google/subcommands"
)

func main() {
	subcommands.Register(subcommands.HelpCommand(), "")
	subcommands.Register(&manager.Status{}, "")
	subcommands.Register(&manager.Add{}, "")
	subcommands.Register(&manager.List{}, "")
	subcommands.Register(&manager.Update{}, "")
	subcommands.Register(&manager.Sign{}, "")
	flag.Parse()
	ctx := context.Background()
	os.Exit(int(subcommands.Execute(ctx)))
}
