package local

// Group represents a group object to be stored
type Group struct {
	ID      string   `json:"id"`
	Configs []string `json:"configs"`
}
