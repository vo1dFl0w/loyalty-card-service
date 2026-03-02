package main

import (
	"context"
	"log"
	"os"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if err := run(ctx); err != nil {
		log.Println(ctx, "startup", "err", err)
		os.Exit(1)
	}
}

func run(ctx context.Context) error {
	// TODO: load config
	// TODO: load logger
	// TODO: start http server
	// TODO: gracefull shutdown
}