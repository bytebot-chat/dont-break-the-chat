package app

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"strconv"
	"time"

	uuid "github.com/satori/go.uuid"
)

const JOBS_REDIS_KEY = "jobs"

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

// InfoString returns a string containing the information about the job
func (j *Job) InfoString() string {
	// Name, Description, Duration, Payout
	return fmt.Sprintf("**%s**\t%s\tDuration: %d\tPayout: %d", j.Name, j.Description, j.Duration(), j.Payout)
}

// Duration returns the duration of the job in hours
func (j *Job) Duration() int {
	return int(j.ExpiresAt - j.CreatedAt)
}

// getAvailableJobs gets the aviailable jobs for the given user profile
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
		j := Job{
			ID:          uuid.NewV4(),
			Name:        "Job " + strconv.Itoa(i),                                           // TODO: Generate a random name
			Description: "This is job " + strconv.Itoa(i),                                   // TODO: Generate a random description
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
