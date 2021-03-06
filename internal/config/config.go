package config

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"
)

// ErrorHandling defines how Config behaves on and error.
type ErrorHandling int

// These constants determines error behaviour for Config.
const (
	ReturnError ErrorHandling = iota // Return a descriptive error.
	ExitOnError                      // Call os.Exit(2).
)

// MissingPropertyError is used when calling a getter with required = true and the property is missing.
type MissingPropertyError string

func (m MissingPropertyError) Error() string {
	return fmt.Sprintf("Property %q is missing", m)
}

// Predefined Config
var std = New(ReturnError)

// Config object holds configuration properties. All properties are saved as strings as they are strings in env and file form.
// Casting is done by access functions. All getter functions will block until Init is run.
type Config struct {
	properties    map[string]string
	errorhandling ErrorHandling
	initOnce      sync.Once
	waitForInit   sync.WaitGroup
}

// New creates a new Config object
func New(errorHandling ErrorHandling) *Config {
	conf := &Config{
		properties:    make(map[string]string),
		errorhandling: errorHandling,
	}

	conf.waitForInit.Add(1)

	return conf
}

// SetErrorhandling sets the error behaviour on Config object
func (c *Config) SetErrorhandling(errorHandling ErrorHandling) {
	c.errorhandling = errorHandling
}

// SetErrorhandling calls SetErrorhandling on the standard Config object.
func SetErrorhandling(errorHandling ErrorHandling) {
	std.SetErrorhandling(errorHandling)
}

// Init reads properties from several files and the environment if set to true. Properties with the same name will be overwritten
// by following files and environment variables will be read last. Is only called once. Will return on first error.
func (c *Config) Init(readEnv bool, paths ...string) error {
	var err error
	init := func() {
		defer c.waitForInit.Done()
		for _, path := range paths {
			if path == "" {
				continue
			}
			if e := c.handleError(c.readPropertyFile(path)); e != nil {
				err = c.handleError(e)
				return
			}
		}

		if readEnv {
			c.readEnv()
		}
	}

	c.initOnce.Do(init)

	return err
}

// Init calls Init on the standard config object.
func Init(readEnv bool, paths ...string) error {
	return std.Init(readEnv, paths...)
}

// readEnv reads properties from environment variables into the Config object.
func (c *Config) readEnv() {
	for _, envVar := range os.Environ() {
		pair := strings.Split(envVar, "=")
		c.properties[pair[0]] = pair[1]
	}
}

// readPropertyFile reads properties from a .properties file into the Config object.
func (c *Config) readPropertyFile(filepath string) error {
	file, err := os.Open(filepath)
	if err != nil {
		return c.handleError(err)
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

	return c.handleError(scanner.Err())
}

// GetString gets a property and casts it to a string.
func (c *Config) GetString(prop string, required bool) (string, error) {
	c.waitForInit.Wait()

	val, exists := c.properties[prop]
	if required && !exists {
		return val, c.handleError(MissingPropertyError(prop))
	}

	return val, nil
}

// GetString calls GetString on the standard Config object
func GetString(prop string, required bool) (string, error) {
	return std.GetString(prop, required)
}

// GetInt gets a property and casts it to an int.
func (c *Config) GetInt(prop string, required bool) (int, error) {
	c.waitForInit.Wait()

	var i int
	val, exists := c.properties[prop]
	if required && !exists {
		return i, c.handleError(MissingPropertyError(prop))
	}

	if !exists {
		return i, nil
	}

	i, err := strconv.Atoi(val)
	return i, c.handleError(err)
}

// GetInt calls GetInt on the standard Config object
func GetInt(prop string, required bool) (int, error) {
	return std.GetInt(prop, required)
}

// GetFloat gets a property and casts it to a float
func (c *Config) GetFloat(prop string, required bool) (float64, error) {
	c.waitForInit.Wait()

	var f float64
	val, exists := c.properties[prop]
	if required && !exists {
		return f, c.handleError(MissingPropertyError(prop))
	}

	if !exists {
		return f, nil
	}

	f, err := strconv.ParseFloat(val, 64)
	return f, c.handleError(err)
}

// GetFloat calls GetFloat on the standard Config object
func GetFloat(prop string, required bool) (float64, error) {
	return std.GetFloat(prop, required)
}

// GetBool gets a property and casts it to a bool
func (c *Config) GetBool(prop string, defaul bool, required bool) (bool, error) {
	c.waitForInit.Wait()

	val, exists := c.properties[prop]
	if required && !exists {
		return defaul, c.handleError(MissingPropertyError(prop))
	}

	if !exists {
		return defaul, nil
	}

	b, err := strconv.ParseBool(val)
	return b, c.handleError(err)
}

// GetBool calls GetBool on the standard Config object
func GetBool(prop string, defaul bool, required bool) (bool, error) {
	return std.GetBool(prop, defaul, required)
}

func (c *Config) handleError(err error) error {
	switch c.errorhandling {
	case ReturnError:
		return err
	case ExitOnError:
		os.Exit(2)
	}

	return nil
}
