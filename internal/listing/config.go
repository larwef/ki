package listing

import (
	"encoding/json"
	"time"
)

// Config represents a config resource to be listed
type Config struct {
	ID           string          `json:"id"`
	Name         string          `json:"name"`
	LastModified time.Time       `json:"lastModified"`
	Version      int             `json:"version"`
	Group        string          `json:"group"`
	Properties   json.RawMessage `json:"properties"`
}
