package app

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

/*
The Job System

*/

const jobsHelpResponse = `
** Jobs **
The jobs system allows you to earn money by taking on randomized jobs.
Jobs are scaled to your level, so the higher your level, the more money you can earn.

** Commands **
- !jobs - Get a list of available jobs
- !jobs help - Get help with the jobs system (you're looking at it)
- !jobs list - Get a list of available jobs
- !jobs take <job> - Take a job
`

const jobsUnknownCommandResponse = `I don't know what you mean by that. Try !jobs help.`

// the Job struct represents a quest that a user can take on
type Job struct {
	// the name of the job
	Name string
	// the description of the job
	Description string
	// the amount of money the user will earn if they complete the job
	// This is an interface to allow jobs to compute their own reward based on the user's stats
	Payout func(p *Profile) int
	// the amount of time in seconds the job will take to complete
	// This is an interface to allow jobs to compute their own time based on the user's stats
	Time func(p *Profile) int
}

// handleJob handles the !job command. It represents the entrypoint for the job system.
// to parse the commands, each successive handler function should strip the 0th element from the splitCmd slice
// and pass the rest of the slice to the next function until the command is fully parsed.
func handleJobs(a *App, m *Message) error {
	// Split the incoming message into a slice of strings
	splitCmd := strings.Split(m.Content, " ")

	// Pop the first element off the slice to get the command
	cmd, splitCmd := splitCmd[0], splitCmd[1:]

	// Make sure the command is !job
	if cmd != "!jobs" {
		return errors.New("invalid command for handleJob. expected !job, got " + cmd)
	}

	// If there are no more elements in the slice, we're at the end of the command chain
	// and we can handle the command.
	if len(splitCmd) == 0 {
		return handleJobsList(a, m, []string{})
	}

	// Otherwise, we need to keep parsing the command.
	// Pop the next element off the slice to get the subcommand
	subCmd, splitCmd := splitCmd[0], splitCmd[1:]

	// Switch on the subcommand
	switch subCmd {
	case "list":
		return handleJobsList(a, m, splitCmd)
	case "start":
		return handleJobsStart(a, m, splitCmd)
	case "help":
		return handleJobsHelp(a, m, splitCmd)
	default:
		return handleJobsUnknownCommand(a, m, splitCmd)
	}
}

// handleJobList handles the !jobs list command. It lists all available jobs.
func handleJobsList(a *App, m *Message, splitCmd []string) error {
	// Get the user's profile

	a.logger.Info().
		Str("user", m.Author.Username).
		Msg("getting user profile")

	profile, err := a.getProfile(m.Author.ID)
	if err != nil {
		return err
	}

	a.logger.Info().
		Str("user", m.Author.Username).
		Msg("generating jobs")

	// Generate a list of jobs
	jobs := generateJobs()

	// Initialize empty string to hold the list of jobs
	jobList := ""

	// Loop through the jobs and add them to the list
	a.logger.Info().
		Str("user", m.Author.Username).
		Int("jobs", len(jobs)).
		Msg("adding jobs to list")
	for _, job := range jobs {
		jobList += fmt.Sprintf("- %s - %s - %d seconds\n", job.Name, job.Description, job.Time(profile))
	}

	// Send the list of jobs to the user
	a.logger.Info().
		Str("user", m.Author.Username).
		Str("jobs", jobList).
		Str("channel", m.ChannelID).
		Msg("sending job list to user")
	msg := m.RespondToChannelOrThread("dbtg", jobList, true, false)

	return a.handleOutgoingMessage(msg)
}

// handleJobStart handles the !job start command. It starts a job for the user.
func handleJobsStart(a *App, m *Message, splitCmd []string) error {
	a.logger.Info().
		Str("user", m.Author.Username).
		Msg("User requested to start a job")
	return nil
}

// handleJobHelp handles the !job help command. It displays help for the job system.
func handleJobsHelp(a *App, m *Message, splitCmd []string) error {
	a.logger.Info().
		Str("user", m.Author.Username).
		Msg("User requested help message")
	// Generate the help message
	help := jobsHelpResponse

	// Send the help message to the user
	msg := m.RespondToChannelOrThread("dbtg", help, true, false)

	return a.handleOutgoingMessage(msg)
}

// handleJobUnknownCommand handles an unknown command. It displays help for the job system.
func handleJobsUnknownCommand(a *App, m *Message, splitCmd []string) error {
	a.logger.Info().
		Str("user", m.Author.Username).
		Str("command", strings.Join(splitCmd, " ")).
		Str("channel", m.ChannelID).
		Str("guild", m.GuildID).
		Msg("User requested unknown command")
	// Generate the help message
	help := jobsUnknownCommandResponse

	// Send the help message to the user
	msg := m.RespondToChannelOrThread("dbtg", help, true, false)

	return a.handleOutgoingMessage(msg)
}

// generateJobs generates a slice of jobs that the user can take on
func generateJobs() []Job {
	// Initialize empty slice of jobs
	jobs := []Job{}

	// Add jobs to the slice
	for i := 0; i < 10; i++ {
		jobs = append(jobs, newSimpleJob("Job "+strconv.Itoa(i), "This is a simple job", 100, 10))
	}
	return jobs
}

// newSimpleJob is a simple job that pays a fixed amount of money and takes a fixed amount of time
func newSimpleJob(name, description string, payout, time int) Job {
	return Job{
		Name:        name,
		Description: description,
		Payout: func(p *Profile) int {
			return payout
		},
		Time: func(p *Profile) int {
			return time
		},
	}
}
