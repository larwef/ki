package config

import (
	"bufio"
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
	mu         sync.RWMutex
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

	for _, envVar := range os.Environ() {
		pair := strings.Split(envVar, "=")
		c.properties[pair[0]] = pair[1]
	}
}

// ReadEnv calls ReadEnv on the standard Config object
func ReadEnv() {
	std.ReadEnv()
}

// ReadPropertyFile reads properties from a .properties file into the Config object.
func (c *Config) ReadPropertyFile(filepath string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	file, err := os.Open(filepath)
	if err != nil {
		return err
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

	return scanner.Err()
}

// ReadPorpertyFile calls ReadPorpertyFile on the standard Config object
func ReadPorpertyFile(filepath string) {
	std.ReadPropertyFile(filepath)
}

// GetString gets a property and casts it to a string.
func (c *Config) GetString(prop string) string {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return c.properties[prop]
}

// GetString calls GetString on the standard Config object
func GetString(prop string) string {
	return std.GetString(prop)
}

// GetInt gets a property and casts it to an int.
func (c *Config) GetInt(prop string) (int, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return strconv.Atoi(c.properties[prop])
}

// GetInt calls GetInt on the standard Config object
func GetInt(prop string) (int, error) {
	return std.GetInt(prop)
}

// GetFloat gets a property and casts it to a float
func (c *Config) GetFloat(prop string) (float64, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return strconv.ParseFloat(c.properties[prop], 64)
}

// GetFloat calls GetFloat on the standard Config object
func GetFloat(prop string) (float64, error) {
	return std.GetFloat(prop)
}

// GetBool gets a property and casts it to a bool
func (c *Config) GetBool(prop string, defaul bool) (bool, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	s, exists := c.properties[prop]
	if !exists {
		return defaul, nil
	}

	return strconv.ParseBool(s)
}

// GetBool calls GetBool on the standard Config object
func GetBool(prop string, defaul bool) (bool, error) {
	return std.GetBool(prop, defaul)
}
