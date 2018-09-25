package memory

import (
	"encoding/json"
	"time"
)

// Config represents a config resource to be stored
type Config struct {
	ID           string
	Name         string
	LastModified time.Time
	Version      int
	Group        string
	Properties   json.RawMessage
}
