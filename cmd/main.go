package main

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/harlequix/godoist"
	"github.com/alecthomas/kong"
)

type loginCmd struct {
	Token string `help:"API Token." koanf:"token"  short:"t" optional:""`
}

type CLI struct {
	Debug       bool     `help:"Enable debug mode." short:"d" optional:"" default:"false"`
	ConfigFiles []string `help:"Path to config file." short:"c" optional:"" type:"existingfile" default:""`
	Login       loginCmd `cmd:"" help:"Login to Todoist."`
}

func main() {
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelInfo})))
	args := &CLI{}
	cli_ctx := kong.Parse(args)
	fmt.Printf("CLI Context: %+v\n", cli_ctx)
	config, err := godoist.BuildConfig(args.ConfigFiles, "TODOIST_", args.Login)
	if err != nil {
		slog.Error("Failed to build config", "error", err)
		os.Exit(1)
	}
	fmt.Printf("args: %+v\n", args)
	fmt.Printf("Config: %+v\n", config)
}
