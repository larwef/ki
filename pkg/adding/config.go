package adding

import (
	"encoding/json"
	"time"
)

// Config represents a config resource to be added
type Config struct {
	ID           string
	Name         string
	LastModified time.Time
	Version      int
	Group        string
	Properties   json.RawMessage
}
