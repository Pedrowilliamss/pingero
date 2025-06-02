package main

import (
	"os"

	"github.com/pedrowilliamss/pingero-cli/internal"
)

func main() {
	cli := internal.CLI{
		MessagePublisher: os.Stdout,
		MessageConsumer:  os.Stdin,
	}

	cli.Start()
}
