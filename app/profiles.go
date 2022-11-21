package app

type Profile struct {
	// The unique ID of the profile. This is the same as the snowflake ID of the user in Discord.
	ID string `json:"id"`
	// The user's inventory
	Inventory Inventory `json:"inventory"`
}
