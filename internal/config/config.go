package config

import (
	"bufio"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"
)

// Predefined Config
var std = New()

// Config object holds configuration properties. All properties are saved as strings as they are strings in env and file form.
// Casting is done by access functions
type Config struct {
	mu         sync.Mutex
	properties map[string]string
}

// New creates a new Config object
func New() *Config {
	return &Config{
		properties: make(map[string]string),
	}
}

// ReadEnv reads properties from environment variables into the Config object.
func (c *Config) ReadEnv() {
	c.mu.Lock()
	defer c.mu.Unlock()

	for _, e := range os.Environ() {
		pair := strings.Split(e, "=")
		c.properties[pair[0]] = pair[1]
	}
}

// ReadPropertyFile reads properties from a .properties file into the Config object.
func (c *Config) ReadPropertyFile(filepath string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	file, err := os.Open(filepath)
	if err != nil {
		log.Printf("Error opening property file %q. Couldn't read properties.", filepath)
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		if line == "" || line[0] == '#' || line[0] == ' ' {
			continue
		}

		pair := strings.Split(line, "=")
		c.properties[pair[0]] = pair[1]
	}

	if err := scanner.Err(); err != nil {
		log.Printf("Error encountered while reading property file: %q. Error: %v", filepath, err)
	}
}

// GetString gets a property and casts it to a string.
func (c *Config) GetString(prop string) string {
	return c.properties[prop]
}

// GetInt gets a property and casts it to an int.
func (c *Config) GetInt(prop string) int {
	i, err := strconv.Atoi(c.properties[prop])

	if err != nil {
		log.Printf("Warning: Error parsing property %q. Trying to parse as an int. Check configuration.", prop)
	}

	return i
}

// GetFloat gets a property and casts it to a float
func (c *Config) GetFloat(prop string) float64 {
	f, err := strconv.ParseFloat(c.properties[prop], 64)
	if err != nil {
		log.Printf("Warning: Error parsing property %q. Trying to parse as a float. Check configuration.", prop)
	}

	return f
}

// GetBool gets a property and casts it to a bool
func (c *Config) GetBool(prop string, defaul bool) bool {
	s, exists := c.properties[prop]
	if !exists {
		return defaul
	}

	b, err := strconv.ParseBool(s)
	if err != nil {
		log.Printf("Warning: Error parsing property %q. Trying to parse as a float. Check configuration.", prop)
	}
	return b
}

// FromEnv calls FromEnv on the standard Config object
func FromEnv() {
	std.ReadEnv()
}

// FromPorpertyFile calls FromPorpertyFile on the standard Config object
func FromPorpertyFile(filepath string) {
	std.ReadPropertyFile(filepath)
}

// GetString calls GetString on the standard Config object
func GetString(prop string) string {
	return std.GetString(prop)
}

// GetInt calls GetInt on the standard Config object
func GetInt(prop string) int {
	return std.GetInt(prop)
}

// GetFloat calls GetFloat on the standard Config object
func GetFloat(prop string) float64 {
	return std.GetFloat(prop)
}

// GetBool calls GetBool on the standard Config object
func GetBool(prop string, defaul bool) bool {
	return std.GetBool(prop, defaul)
}
