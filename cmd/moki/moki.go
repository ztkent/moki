// Command moki is an AI assistant for the command line.
package main

import (
	"context"
	"fmt"
	"os"

	"github.com/ztkent/moki/internal/app"
	"github.com/ztkent/moki/internal/config"
)

func main() {
	cfg, err := config.Parse(os.Args[1:])
	if err != nil {
		fmt.Fprintln(os.Stderr, "moki:", err)
		os.Exit(2)
	}
	if err := cfg.Validate(); err != nil {
		fmt.Fprintln(os.Stderr, "moki:", err)
		os.Exit(2)
	}

	if err := app.Run(context.Background(), cfg); err != nil {
		fmt.Fprintln(os.Stderr, "moki:", err)
		os.Exit(1)
	}
}
