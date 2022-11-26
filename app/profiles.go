package app

import (
	"encoding/json"
	"math/rand"
	"strconv"
)

// Prefix for consistent key names in the database.
const REDIS_PROFILE_PREFIX = "profile:"

// Profile is a struct that represents a user's profile.
// It's the main data structure for the game and tracks the state for a user.
// State is maintained in redis under the top-level key "profile:<user_id>".
// Profiles have the ID field of the user to facilitate lookups and writing the profile back to Redis.
type Profile struct {
	ID        string    `json:"id"`          // The unique ID of the profile. This is the same as the snowflake ID of the user in Discord.
	Inventory Inventory `json:"inventory"`   // The user's inventory
	Balance   int       `json:"balance"`     // The user's available spending balance
	ActiveJob Job       `json:"current_job"` // Active job
}

// getProfile gets the profile for the given user ID.
// If the profile does not exist, it will be created.
func (a *App) getProfile(userID string) (*Profile, error) {
	// Check for the profile in the database
	p, err := a.redis.Get(a.context, "profile:"+userID).Result()
	if err != nil {
		// If the profile does not exist, create a new profile
		profile := &Profile{
			ID:        userID,
			Inventory: Inventory{},
			Balance:   0,
		}

		// Marshal the profile
		profileBytes, err := json.Marshal(profile)
		if err != nil {
			return nil, err
		}

		// Save the profile to the database
		err = a.redis.Set(a.context, "profile:"+userID, profileBytes, 0).Err()
		if err != nil {
			return nil, err
		}

		return profile, nil
	}

	// If the profile exists, unmarshal it and return it
	var profile Profile
	err = json.Unmarshal([]byte(p), &profile)
	if err != nil {
		return nil, err
	}

	return &profile, nil
}

// saveProfile saves the profile to the database.
func (p *Profile) save(a *App) error {
	// Marshal the profile into a json string
	profileBytes, err := json.Marshal(p)
	if err != nil {
		return err
	}

	// Save the profile to the database
	// TODO: This looks prone to race conditions. Make it not that.
	return a.redis.Set(a.context, "profile:"+p.ID, profileBytes, 0).Err()
}

// work is a method on the profile that handles the work command from chat.
// It should generate some amount of currency and add it to the user's balance.
// This method will eventually take the user's profile and return a scaled or leveled amount of currency.
// For now it just returns a random amount of currency between 1 and 100.
// This function is probably going to be the biggest source of bugs for a long time.
// I have a feeling that we will see a lot of race conditions and concurrency issues here. I can't wait.
func (p *Profile) work(a *App) (int, error) {
	// Generate a random amount of currency between 1 and 100
	// TODO: This should eventually be scaled or leveled based on the user's profile/level/job/other attributes
	// Make sure always earn at least 1 buck
	earned := rand.Intn(100) + 1 // nolint:gosec // This is not a security issue

	// Add the amount earned to the user's balance
	p.Balance += earned

	// Marshal the profile
	profileBytes, err := json.Marshal(p)
	if err != nil {
		return 0, err
	}

	// Save the profile to the database
	err = a.redis.Set(a.context, REDIS_PROFILE_PREFIX+":"+p.ID, profileBytes, 0).Err()
	if err != nil {
		return 0, err
	}

	return earned, nil
}

// getBalance returns the user's balance.
func (p *Profile) getBalance() int {
	return p.Balance
}

// getBalanceString returns the user's balance as a string.
func (p *Profile) getBalanceString() string {
	return strconv.Itoa(p.Balance)
}
