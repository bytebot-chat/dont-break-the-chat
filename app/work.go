package app

/*
The Work System

The work system is a simple system that allows users to earn currency by working. It's designed as a way
to experiment with stateful applications in bytebot using only the core features of the framework.

Normally, sane people here would use a real database to store state. But we're not sane people.

At this point in time, the work system is pretty simple. Punch the clock, get paid. No cooldown. It's a start.

Right now what I'd like to see is just that users can work, get paid, and check their balance. Eventually we
can add a cooldown, and maybe even a way to check how long until you can work again.

*/

// handleWork handles the !work command. It represents the entrypoint for the work system.
// to parse the commands, each function should strip the 0th element from the splitCmd slice
// and pass the rest of the slice to the next function until the command is fully parsed.
func handleWork(a *App, m *Message) error {
	return nil
}
