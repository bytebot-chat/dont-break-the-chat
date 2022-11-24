package app

import (
	"errors"
	"strings"
)

/*
The Currency System

This file represents the entrypoint for the currency system. It handles the !balance and future
commands related to managing a users's wallet or balance.

At this point in time, the system is dead simple:
- Integers only
- Users cannot send currency to other users
- Users cannot check the balance of other users
- Users cannot check the inventory of other users
- Users can only earn money. They cannot spend it (yet).

*/

const balanceHelpResponse = `
** Balance **
Check your balance and see how much money you have.

## Commands
- !balance - Check your balance
- !balance help - Get help with the balance system (you're looking at it)
`

const balanceUnknownCommandResponse = `I don't know what you mean by that. Try !balance help.`

// handleBalance handles the !balance command. It represents the entrypoint for the currency system.
// to parse the commands, each successive handler function should strip the 0th element from the splitCmd slice
// and pass the rest of the slice to the next function until the command is fully parsed.
func handleBalance(a *App, m *Message) error {
	// Split the incoming message into a slice of strings
	splitCmd := strings.Split(m.Content, " ")

	// Pop the first element off the slice to get the command
	cmd, splitCmd := splitCmd[0], splitCmd[1:]

	// Make sure the command is !balance
	if cmd != "!balance" {
		return errors.New("invalid command for handleBalance. expected !balance, got " + cmd)
	}

	// If there are no more elements in the slice, we're at the end of the command chain
	// and we can handle the command.
	if len(splitCmd) == 0 {
		return handleBalanceCommand(a, m, []string{})
	}

	// Otherwise, we need to keep parsing the command.
	// Pop the next element off the slice to get the subcommand
	subCmd, splitCmd := splitCmd[0], splitCmd[1:]

	// Switch on the subcommand
	switch subCmd {
	case "help":
		return handleBalanceHelpCommand(a, m, splitCmd)
	default:
		return handleBalanceUnknownCommand(a, m, splitCmd)
	}

	return nil

}

// handleBalanceCommand handles the bare !balance command.
// It should return the user's current balance.
func handleBalanceCommand(a *App, m *Message, args []string) error {
	// Get the user's profile
	profile, err := a.getProfile(m.Author.ID)
	if err != nil {
		return err
	}

	// Respond to the user with their balance
	msg := m.RespondToChannelOrThread("dbtg", "Your balance is "+profile.getBalanceString()+" dollars", true, false)

	return a.handleOutgoingMessage(msg)
}

// handleBalanceHelpCommand handles the !balance help command.
// It should return a help message for the currency system.
func handleBalanceHelpCommand(a *App, m *Message, args []string) error {
	return a.handleOutgoingMessage(m.RespondToChannelOrThread("dbtg", balanceHelpResponse, true, false))
}

// handleBalanceUnknownCommand handles an unknown subcommand for the !balance command.
// It should return an error message.
func handleBalanceUnknownCommand(a *App, m *Message, args []string) error {
	return a.handleOutgoingMessage(m.RespondToChannelOrThread("dbtg", balanceUnknownCommandResponse, true, false))
}
