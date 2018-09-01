package config

import (
	"encoding/json"
	"time"
)

type Config struct {
	Id           string          `json:"id"`
	Name         string          `json:"name"`
	LastModified time.Time       `json:"lastModified"`
	Properties   json.RawMessage `json:"properties"`
}
