package app

type Inventory struct {
	// The unique ID of the inventory. This is the same as the snowflake ID of the user in Discord.
	ID string `json:"id"`
	// The user's currency balance
	Balance int `json:"balance"`
	// The User's demerits
	Demerits int `json:"demerits"`
}
