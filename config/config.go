package config

import (
	"encoding/json"
	"time"
)

type Config struct {
	Id           string          `json:"id"`
	Name         string          `json:"name"`
	LastModified time.Time       `json:"lastModified"`
	Version      int             `json:"version"`
	Group        string          `json:"group"`
	Properties   json.RawMessage `json:"properties"`
}
