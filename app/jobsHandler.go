package app

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/rs/zerolog/log"
)

/*
The Job System

Yes, I am aware there's also "!work" right now. I guess I'll merge the two experimental systems into one eventually.
The job system is meant to be an evolution of the idea of "!work" that allows users to take on jobs of various difficulty and risk.
Upon the first request in the calendar day for that user, the app will generate a list of jobs that the user can take on.
Jobs will have an ID, a name, and must bring their own functions for computing time and rewards.
The user will be able to select a job and the app will start a timer for that job. When the timer is up, the user will receive their reward.
A user will have an "active job" field in their profile that will be set to the job they are currently working on.
If they have no active job, they will be able to start a new job. If they have an active job, they must wait or quit the job.
Quitting a job should carry a penalty, but I haven't decided what that penalty should be yet.
*/

const jobsHelpResponse = `
** Jobs **
The jobs system allows you to earn money by taking on randomized jobs.
Jobs are scaled to your level, so the higher your level, the more money you can earn.

** Commands **
- !jobs 			- Get a list of available jobs
- !jobs help 		- Get help with the jobs system (you're looking at it)
- !jobs list 		- Get a list of available jobs
- !jobs refresh 	- Refresh the list of available jobs
- !jobs take <job> 	- Take a job
`

const jobsUnknownCommandResponse = `I don't know what you mean by that. Try !jobs help.`

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
	case "refresh":
		return handleJobsRefresh(a, m, splitCmd)
	case "start":
		return handleJobsStart(a, m, splitCmd)
	case "take":
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
	// This also initializes the user's profile if it doesn't exist
	a.logger.Info().
		Str("user", m.Author.Username).
		Msg("getting user profile")
	profile, err := a.getProfile(m.Author.ID)
	if err != nil {
		return err
	}

	// Check for existing available jobs
	jobs, err := a.getAvailableJobs(profile)
	// If there's a redis nil error, it means the user has no jobs so we skip error handling
	// So we handle everything else
	if err != nil && err.Error() != "redis: nil" {
		return err
	}

	// If there are no jobs, tell the user in the channel
	if len(jobs) == 0 {
		a.logger.Info().
			Str("user", m.Author.Username).
			Msg("no jobs available")
		msg := m.RespondToChannelOrThread("dbtg", "There are no jobs available right now. Looking for new work...", true, false)
		a.handleOutgoingMessage(msg)
		time.Sleep(5 * time.Second) // Give the impression that the bot is working on something

		// Generate a new list of jobs
		jobs, err = a.generateJobs(profile, 10)

		// If there's an error, return it
		if err != nil {
			m.RespondToChannelOrThread("dbtg", "Nobody's hiring, kid. Come back later.", true, false)
			return err
		}

		// Set the jobs in redis and proceed with the rest of the function
		err = a.setAvailableJobs(profile, jobs)
		if err != nil {
			return err
		}
	}

	// Respond to the user with the list of jobs
	a.logger.Info().
		Str("user", m.Author.Username).
		Msg("sending job list")
	jobString := []string{}
	jobString = append(jobString, "There's some folks looking for help. Here's what they need:\n")
	jobString = append(jobString, "```")
	for _, job := range jobs {
		jobString = append(jobString, job.InfoString())
	}
	jobString = append(jobString, "```")
	jobString = append(jobString, "To take a job, type `!jobs take <job ID>`")
	msg := m.RespondToChannelOrThread("dbtg", strings.Join(jobString, "\n"), true, false)
	return a.handleOutgoingMessage(msg)
}

// handleJobRefresh handles the !jobs refresh command. It generates a new list of jobs.
func handleJobsRefresh(a *App, m *Message, splitCmd []string) error {
	// Get the user's profile
	// This also initializes the user's profile if it doesn't exist
	a.logger.Info().
		Str("user", m.Author.Username).
		Msg("getting user profile")
	profile, err := a.getProfile(m.Author.ID)
	if err != nil {
		return err
	}

	// Generate a new list of jobs
	jobs, err := a.generateJobs(profile, 10)
	if err != nil {
		log.Error().
			Err(err).
			Str("user", m.Author.Username).
			Msg("error generating jobs")
		m.RespondToChannelOrThread("dbtg", "Nobody's hiring, kid. Come back later.", true, false)
		return err
	}

	// Set the jobs in redis and proceed with the rest of the function
	err = a.setAvailableJobs(profile, jobs)
	if err != nil {
		log.Error().
			Err(err).
			Str("user", m.Author.Username).
			Msg("error setting jobs in redis")
		m.RespondToChannelOrThread("dbtg", "I found some jobs for you but I lost the paperwork on the way over. Better luck next time.", true, false)
		return err
	}

	// Respond to the user with the list of jobs
	// TODO: This is a copy/paste of handleJobsList. Refactor this into a function.
	a.logger.Info().
		Str("user", m.Author.Username).
		Msg("sending job list")

	// Format a multi-line string with the job info
	jobString := []string{}
	jobString = append(jobString, "There's some folks looking for help. Here's what they need:")
	jobString = append(jobString, "```") // Use a code block to make it look nice

	// Loop through the jobs and append the info to the string
	for _, job := range jobs {
		jobString = append(jobString, job.InfoString())
	}

	jobString = append(jobString, "```")                                                  // Close the code block
	jobString = append(jobString, "To take a job, type `!jobs take <job ID>`")            // Tell the user how to take a job
	msg := m.RespondToChannelOrThread("dbtg", strings.Join(jobString, "\n"), true, false) // Generate a reply message
	return a.handleOutgoingMessage(msg)                                                   // Send the message
}

// handleJobStart handles the !job start command. It starts a job for the user.
// Jobs are stored in a slice in the user's profile
// Profile -> Jobs -> Job by index
// An active job must be removed from the AvailableJobs slice and added to the ActiveJobs field
// And then a goroutine must be started to act as a timer for the job.
// It should sleep for the duration of the job and then grant the user the reward and remove the job from the ActiveJobs slice when it returns.
func handleJobsStart(a *App, m *Message, splitCmd []string) error {

	// Make sure splitCmd is not empty and contains an integer
	if len(splitCmd) == 0 {
		a.logger.Info().
			Str("user", m.Author.Username).
			Msg("no job ID provided")
		a.handleOutgoingMessage(m.RespondToChannelOrThread("dbtg", "You need to provide a job ID. Type `!jobs list` to see a list of available jobs.", true, false))
		return errors.New("no job ID provided")
	}

	// Make sure splitCmd[0] is a valid, positive integer
	jobID, err := strconv.Atoi(splitCmd[0])
	if err != nil || jobID < 0 {
		a.logger.Error().
			Err(err).
			Str("user", m.Author.Username).
			Msg("error converting job ID to valid integer for indexing")
		a.handleOutgoingMessage(m.RespondToChannelOrThread("dbtg", "That's not a valid job ID. Type `!jobs list` to see a list of available jobs.", true, false))
		return err
	}

	// Get the user's profile
	profile, err := a.getProfile(m.Author.ID)
	if err != nil {
		a.logger.Error().
			Err(err).
			Str("user", m.Author.Username).
			Msg("error getting user profile")
		return err
	}
	a.logger.Info().
		Str("user", m.Author.Username).
		Msg("User profile retrieved")

	// Check for existing available jobs
	jobs, err := a.getAvailableJobs(profile)
	if err != nil {
		a.logger.Error().
			Err(err).
			Str("user", m.Author.Username).
			Msg("error getting available jobs")
		return err
	}
	a.logger.Info().
		Str("user", m.Author.Username).
		Msgf("%d available jobs retrieved", len(jobs))

	// Make sure the job ID is valid
	if jobID < 0 || jobID >= len(jobs) {
		a.logger.Error().
			Str("user", m.Author.Username).
			Int("jobID", jobID).
			Int("numJobs", len(jobs)).
			Msg("Invalid job ID: job ID out of range")
		return errors.New("invalid job ID")
	}

	// Assign the job to the user's ActiveJob field
	// Make a copy of the job so we don't modify the original
	activeJob := jobs[jobID]
	profile.ActiveJob = activeJob

	// Save profile
	err = profile.save(a)
	if err != nil {
		a.logger.Error().
			Err(err).
			Str("user", m.Author.Username).
			Msg("error saving profile")
		return err
	}

	a.logger.Info().
		Str("user", m.Author.Username).
		Str("job", profile.ActiveJob.Name).
		Str("id", profile.ActiveJob.ID.String()).
		Msg("Job assigned to user")

	// Remove the job from the AvailableJobs slice
	// and tell the app to save the jobs to the database
	jobs = append(jobs[:jobID], jobs[jobID+1:]...)
	err = a.setAvailableJobs(profile, jobs)
	if err != nil {
		a.logger.Error().
			Err(err).
			Str("user", m.Author.Username).
			Msg("error saving profile")
		return err
	}

	// Start a goroutine to act as a timer for the job
	go profile.ActiveJob.work(a, profile.ID)
	a.logger.Info().
		Str("user", m.Author.Username).
		Str("job", profile.ActiveJob.Name).
		Str("id", profile.ActiveJob.ID.String()).
		Msg("Job started")

	// Construct a message to send to the user with the job info and a timer
	jobAcceptedMessage := fmt.Sprintf("'%s', eh? I'll let the boss know you're on that one. Get lost.", profile.ActiveJob.Name)
	jobAcceptedMessage += fmt.Sprintf(" You've got %d seconds to get it done.", profile.ActiveJob.Duration())
	jobAcceptedMessage += fmt.Sprintf(" If you don't get it done in time, I'll be taking your %d credits.", profile.ActiveJob.Payout)

	// Send the message
	return a.handleOutgoingMessage(m.RespondToChannelOrThread("dbtg", jobAcceptedMessage, true, false))
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
