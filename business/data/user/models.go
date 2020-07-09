package user

// User represents someone with access to the system.
type User struct {
	ID           string `json:"id"`
	SourceID     string `json:"source_id"`
	Source       string `json:"source"`
	ScreenName   string `json:"screen_name"`
	Name         string `json:"name"`
	Location     string `json:"location"`
	FriendsCount int    `json:"friends_count"`
	Friends      []User `json:"friends"`
}

// NewUser contains information needed to create a new User.
type NewUser struct {
	SourceID     string `json:"source_id"`
	Source       string `json:"source"`
	ScreenName   string `json:"screen_name"`
	Name         string `json:"name"`
	Location     string `json:"location"`
	FriendsCount int    `json:"friends_count"`
	Friends      []User `json:"friends"`
}

type addResult struct {
	AddUser struct {
		User []struct {
			ID string `json:"id"`
		} `json:"user"`
	} `json:"addUser"`
}

func (addResult) document() string {
	return `{
		user {
			id
		}
	}`
}
