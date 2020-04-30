package main

import (
	"context"
	"log"
	"os"

	"github.com/int128/awsswitch/pkg/cmd"
)

func main() {
	log.SetFlags(0)
	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Printf("error: %s", err)
		os.Exit(1)
	}
}
