package group

// Group represents a group object
type Group struct {
	ID      string   `json:"id"`
	Configs []string `json:"configs"`
}
