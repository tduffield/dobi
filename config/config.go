package config

import (
	"io/ioutil"
	"path/filepath"
	"sort"

	log "github.com/Sirupsen/logrus"
)

// Resource is an interface for each configurable type
type Resource interface {
	Dependencies() []string
	Validate(config *Config) error
}

// Config is a data object for a full config file
type Config struct {
	Resources  map[string]Resource
	WorkingDir string
}

// NewConfig returns a new Config object
func NewConfig() *Config {
	return &Config{
		Resources: make(map[string]Resource),
	}
}

// Sorted returns the list of resource names in alphabetical sort order
func (c *Config) Sorted() []string {
	size := len(c.Resources)
	names := make([]string, size, size)
	for name := range c.Resources {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}

func (c *Config) missingResources(names []string) []string {
	missing := []string{}
	for _, name := range names {
		if _, ok := c.Resources[name]; !ok {
			missing = append(missing, name)
		}
	}
	return missing
}

// Load a configuration from a filename
func Load(filename string) (*Config, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	config, err := LoadFromBytes(data)
	if err != nil {
		return nil, err
	}
	log.WithFields(log.Fields{"filename": filename}).Debug("Configuration loaded")

	absPath, err := filepath.Abs(filename)
	if err != nil {
		return nil, err
	}
	config.WorkingDir = filepath.Dir(absPath)

	if err = validate(config); err != nil {
		return nil, err
	}
	return config, nil
}

func validate(config *Config) error {
	for _, resource := range config.Resources {
		if err := resource.Validate(config); err != nil {
			return err
		}
	}
	return nil
}
