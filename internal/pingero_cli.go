package internal

import (
	"fmt"
	"io"
	"os"

	"github.com/pedrowilliamss/pingero-cli/internal/commands"
)

type CliCommand int

const (
	CmdStart CliCommand = 1
	CmdView  CliCommand = 2
	CmdStop  CliCommand = 0
)

type Command interface {
	Execute()
}

type CLI struct {
	MessagePublisher io.Writer
	MessageConsumer  io.Reader
}

var urls = []string{"https://google.com", "https://www.youtube.com/", "https://www.linkedin.com/"}

func (c *CLI) Start() {
	c.displayInitializeMessage()

	startCommand := commands.StartPingeroCommand{MessagePublisher: c.MessagePublisher}
	urlsFileName := startCommand.Start(urls)

	monitorCommand := commands.CreateMonitorUrlsCommand(c.MessagePublisher, urlsFileName)
	c.displayMenu()

	for {
		c.handleCommand(monitorCommand)
	}
}

func (c *CLI) handleCommand(monitorUrlsCommand *commands.MonitorUrlsCommand) {
	switch c.readCommand() {
	case CmdStart:
		monitorUrlsCommand.Execute()
	case CmdStop:
		monitorUrlsCommand.Stop()
		os.Exit(-1)
	default:
		c.MessagePublisher.Write([]byte("Command not recognized"))
	}
}

func (c *CLI) displayInitializeMessage() {
	c.MessagePublisher.Write([]byte("Welcome to Pingero 0.0.1"))
}

func (c *CLI) displayMenu() {
	fmt.Fprint(c.MessagePublisher, `
1 - Start Monitoring
2 - View Logs
0 - Exit the Program
`)
}

func (c *CLI) readCommand() CliCommand {
	var command int
	_, err := fmt.Fscanf(c.MessageConsumer, "%d\n", &command)
	if err != nil {
		c.MessagePublisher.Write([]byte("Failed to read command\n"))
		return -1
	}
	return CliCommand(command)
}
