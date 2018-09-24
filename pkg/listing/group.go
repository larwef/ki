package listing

// Group represents a group object to be listed
type Group struct {
	ID      string   `json:"id"`
	Configs []string `json:"configs"`
}
