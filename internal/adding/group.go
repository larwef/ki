package adding

// Group represents a group object to be added
type Group struct {
	ID      string   `json:"id"`
	Configs []string `json:"configs"`
}
