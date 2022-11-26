package app

import "strings"

const infoResponse = `
** Don't Break the Chat ** is an experimental chat-based game using the Bytebot ecosystem. It's a work in progress.

Follow the project on Github at https://github.com/bytebot-chat/dont-break-the-chat
`
const helpResponse = `
** Don't Break the Chat ** is an experimental chat-based game using the Bytebot ecosystem. It's a work in progress.

## How to play
jk there's no way to play yet

## Commands
- !info - Get information about the game
- !help - Get help with the game
- !work - Work for money
- !balance - Check your balance
- !jobs - List available jobs and their requirements

File an issue on Github at https://github.com/bytebot-chat/dont-break-the-chat/issues
`

func handleCommand(a *App, m *Message) error {
	splitCmd := strings.Split(m.Content, " ")
	cmd := splitCmd[0]

	switch cmd {
	case "!info":
		return handleInfo(a, m)
	case "!help":
		return handleHelp(a, m)
	case "!work":
		return handleWork(a, m)
	case "!balance":
		return handleBalance(a, m)
	case "!jobs":
		return handleJobs(a, m)
	default:
		return handleUnknownCommand(a, m)
	}
}

// handleInfo handles the !info command.
func handleInfo(a *App, m *Message) error {
	resp := m.RespondToChannelOrThread("dbtg", infoResponse, true, false)
	return a.handleOutgoingMessage(resp)
}

// handleHelp handles the !help command.
func handleHelp(a *App, m *Message) error {
	return nil
}

// handleUnknownCommand handles an unknown command.
func handleUnknownCommand(a *App, m *Message) error {
	return nil
}
