package app

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
	return nil
}

// handleJobList handles the !jobs list command. It lists all available jobs.
func handleJobsList(a *App, m *Message, splitCmd []string) error {
	return nil
}

// handleJobStart handles the !job start command. It starts a job for the user.
func handleJobsStart(a *App, m *Message, splitCmd []string) error {
	return nil
}

// generateJobs generates a slice of jobs that the user can take on
func generateJobs() []Job {
	// Initialize empty slice of jobs
	jobs := []Job{}

	return jobs
}

// simpleJob is a simple job that pays a fixed amount of money and takes a fixed amount of time
func simpleJob(name, description string, payout, time int) Job {
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
