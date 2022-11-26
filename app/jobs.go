package app

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"time"

	uuid "github.com/satori/go.uuid"
)

const JOBS_REDIS_KEY = "jobs"

var QUEST_TITLE_VERBS = [...]string{
	"Find",
	"Collect",
	"Deliver",
	"Steal",
	"Retrieve",
	"Return",
	"Destroy",
	"Kill",
	"Capture",
	"Rescue",
	"Escort",
	"Protect",
	"Defend",
	"Assassinate",
}

var QUEST_TITLE_NOUNS = [...]string{
	"the airlock",
	"the bridge",
	"a stim pack",
	"the captain",
	"the captain's cat",
	"your boss's space suit",
	"a rat meat pie",
	"the last iPod shuffle",
}

// The Job struct is the base struct for all jobs
// Eventually this should become an interface to allow for more complex jobs
// But for now I don't want to figure out how to mess with the JSON unmarshaller to handle interfaces
type Job struct {
	ID          uuid.UUID `json:"id"`          // The ID of the job
	Name        string    `json:"name"`        // The name of the job
	Description string    `json:"description"` // The description of the job
	Payout      int       `json:"payout"`      // Amount of currency the user gets for completing the job
	CreatedAt   int64     `json:"created_at"`  // The time the job was created
	ExpiresAt   int64     `json:"expires_at"`  // The time the job expires
}

// work does the work for the given job
// It sleeps for the duration of the job and then updates the user's balance with the payout before returning
func (j *Job) work(a *App, userID string) error {
	// Get the user's profile
	profile, err := a.getProfile(userID)
	if err != nil {
		return err
	}

	// Calculate the duration and payout of the job before sleeping
	duration := j.Duration()
	payout := j.Payout

	// Sleep for the duration of the job
	a.logger.Debug().
		Str("user", userID).
		Str("job", j.ID.String()).
		Int("duration", duration).
		Msg("Job started")

	time.Sleep(time.Duration(duration) * time.Second)

	// Update the user's balance with the payout
	profile.Balance += payout

	// Save the user's profile
	err = profile.save(a)
	if err != nil {
		return err
	}

	// Log the job completion
	a.logger.Debug().
		Str("user", userID).
		Str("job", j.ID.String()).
		Int("duration", duration).
		Int("payout", payout).
		Msg("Job completed")

	return nil
}

// saveAvailableJobs saves the list of available jobs to the database
func (a *App) setAvailableJobs(p *Profile, jobs []Job) error {
	// Marshal the jobs into a json string
	jobsBytes, err := json.Marshal(jobs)
	if err != nil {
		return err
	}

	// Save the jobs to the database
	err = a.redis.Set(a.context, JOBS_REDIS_KEY+":"+p.ID, jobsBytes, 0).Err()
	if err != nil {
		return err
	}

	return nil
}

// Duration returns the duration of the job in hours
func (j *Job) Duration() int {
	return int(j.ExpiresAt - j.CreatedAt)
}

// timeRemaining returns the time remaining for the job in seconds
func (j *Job) timeRemaining() int {
	return int(j.ExpiresAt - time.Now().Unix())
}

// getAvailableJobs gets the available jobs for the given user profile
// if the user has no available jobs, it returns an empty list
func (a *App) getAvailableJobs(p *Profile) ([]Job, error) {
	// Get the list of available jobs from the database
	j, err := a.redis.Get(a.context, JOBS_REDIS_KEY+":"+p.ID).Result()
	if err != nil {
		return []Job{}, nil
	}

	// Unmarshal the list of jobs
	var jobs []Job
	err = json.Unmarshal([]byte(j), &jobs)
	if err != nil {
		return nil, err
	}

	// Return the list of jobs
	return jobs, nil
}

// generateJobs generates a new set of jobs for the given user profile
func (a *App) generateJobs(p *Profile, count int) ([]Job, error) {
	// Generate a list of jobs with randomized names, descriptions, durations, and payouts
	jobs := []Job{}
	for i := 0; i < count; i++ {
		// Generate a new job
		name, desc := generateJobNameAndDescription()
		j := Job{
			ID:          uuid.NewV4(),
			Name:        name,                                                               // TODO: Generate a random name
			Description: desc,                                                               // TODO: Generate a random description
			Payout:      rand.Intn(950) + 50,                                                // 50 minimum, 1000 maximum
			CreatedAt:   time.Now().Unix(),                                                  // Now
			ExpiresAt:   time.Now().Add(time.Duration(rand.Intn(24)+24) * time.Hour).Unix(), // 24-48 hours from now
		}
		// Add the job to the list of jobs
		jobs = append(jobs, j)
	}

	// Debug log the generated jobs
	a.logger.Debug().
		Str("user", p.ID).
		Interface("jobs", jobs).
		Msg("Generated new jobs")

	// Return the list of jobs
	return jobs, nil
}

// generateJobNameAndDescription generates a random job name and description
func generateJobNameAndDescription() (string, string) {
	// Generate a random name
	name := fmt.Sprintf("%s %s", QUEST_TITLE_VERBS[rand.Intn(len(QUEST_TITLE_VERBS))], QUEST_TITLE_NOUNS[rand.Intn(len(QUEST_TITLE_NOUNS))])

	// Generate a random description
	desc := fmt.Sprintf("I need you to %s.", name)

	// Return the name and description
	return name, desc
}
