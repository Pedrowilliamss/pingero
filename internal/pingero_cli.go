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
	CmdStop  CliCommand = 3
	CmdExit  CliCommand = 0
)

type CLI struct {
	MessagePublisher io.Writer
	MessageConsumer  io.Reader
}

var urls = []string{"https://google.com", "https://www.youtube.com/", "https://www.linkedin.com/"}

func (c *CLI) Start() {
	c.displayInitializeMessage()

	startCommand := commands.StartPingeroCommand{MessagePublisher: c.MessagePublisher}
	urlsAggregation := startCommand.Execute(urls)

	monitorCommand := commands.CreateMonitorUrlsCommand(c.MessagePublisher, urlsAggregation)
	viewLogs := commands.CreateViewLogsCommand(c.MessagePublisher, urlsAggregation)
	c.displayMenu()

	for {
		c.handleCommand(monitorCommand, viewLogs)
	}
}

func (c *CLI) handleCommand(monitorUrlsCommand *commands.MonitorUrlsCommand, viewLogs *commands.ViewLogsCommand) {
	switch c.readCommand() {
	case CmdStart:
		viewLogs.Execute()
		monitorUrlsCommand.Execute()
	case CmdView:
		viewLogs.Execute()
	case CmdStop:
		c.displayStopMonitoringMessage()
		monitorUrlsCommand.Stop()
		viewLogs.Stop()
		c.displayMenu()
	case CmdExit:
		monitorUrlsCommand.Stop()
		viewLogs.Stop()
		os.Exit(-1)
	default:
		c.MessagePublisher.Write([]byte("Command not recognized"))
	}
}

func (c *CLI) displayStopMonitoringMessage() {
	c.MessagePublisher.Write([]byte("Stop monitoring..."))
}

func (c *CLI) displayInitializeMessage() {
	c.MessagePublisher.Write([]byte("Welcome to Pingero 0.0.1\n"))
}

func (c *CLI) displayMenu() {
	fmt.Fprint(c.MessagePublisher, `
1 - Start Monitoring
2 - View Logs
3 - Stop Monitoring
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
