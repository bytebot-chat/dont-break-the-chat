package app

import (
	"errors"
	"strings"
)

/*
The Work System

The work system is a simple system that allows users to earn currency by working. It's designed as a way
to experiment with stateful applications in bytebot using only the core features of the framework.

Normally, sane people here would use a real database to store state. But we're not sane people.

At this point in time, the work system is pretty simple. Punch the clock, get paid. No cooldown. It's a start.

Right now what I'd like to see is just that users can work, get paid, and check their balance. Eventually we
can add a cooldown, and maybe even a way to check how long until you can work again.

*/

const workHelpResponse = `
** Working **
Achieve class consciousness by punching the clock and earning your daily wage.

## Commands
- !work - Punch the clock and earn your daily wage
- !work help - Get help with the work system (you're looking at it)
`

const workUnknownCommandResponse = `I don't know what you mean by that. Try !work help.`

// handleWork handles the !work command. It represents the entrypoint for the work system.
// to parse the commands, each function should strip the 0th element from the splitCmd slice
// and pass the rest of the slice to the next function until the command is fully parsed.
func handleWork(a *App, m *Message) error {
	// Split the incoming message into a slice of strings
	splitCmd := strings.Split(m.Content, " ")

	// Pop the first element off the slice to get the command
	cmd, splitCmd := splitCmd[0], splitCmd[1:]

	// Make sure the command is !work
	if cmd != "!work" {
		return errors.New("invalid command for handleWork. expected !work, got " + cmd)
	}

	// If there are no more elements in the slice, we're at the end of the command chain
	// and we can handle the command.
	if len(splitCmd) == 0 {
		return handleWorkCommand(a, m, []string{})
	}

	// Otherwise, we need to keep parsing the command.
	// Pop the next element off the slice to get the subcommand
	subCmd, splitCmd := splitCmd[0], splitCmd[1:]

	// Switch on the subcommand
	switch subCmd {
	case "help":
		return handleWorkHelp(a, m, splitCmd)
	default:
		return handleWorkUnknownCommand(a, m, splitCmd)
	}
}

// handleWorkCommand handles the bare !work command.
// It should generate some amount of currency and add it to the user's balance.
// The amount of currency should eventually come from a function that takes the user's
// profile and returns a scaled or leveled amount of currency.
func handleWorkCommand(a *App, m *Message, args []string) error {
	resp := m.RespondToChannelOrThread("dbtg", "You punched the clock and earned 0 currency.", true, false)
	return a.handleOutgoingMessage(resp)
}

// handleWorkHelp handles the !work help command.
func handleWorkHelp(a *App, m *Message, args []string) error {
	resp := m.RespondToChannelOrThread("dbtg", workHelpResponse, true, false)
	return a.handleOutgoingMessage(resp)
}

// handleWorkUnknownCommand handles an unknown command for the work system.
func handleWorkUnknownCommand(a *App, m *Message, args []string) error {
	resp := m.RespondToChannelOrThread("dbtg", workUnknownCommandResponse, true, false)
	return a.handleOutgoingMessage(resp)
}
