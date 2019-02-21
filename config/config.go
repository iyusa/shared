package config

import (
	"strings"

	"gopkg.in/ini.v1"
)

// Config general configuration from ini file
type Config struct {
	cfg *ini.File
}

// LoadIni from disk
func (c *Config) LoadIni(iniFile string) (e error) {
	c.cfg, e = ini.Load(iniFile)
	return
}

// GetString from config
func (c *Config) GetString(section string, key string, defValue string) string {
	if c.cfg.Section(section).HasKey(key) {
		return c.cfg.Section(section).Key(key).String()
	}
	return defValue
}

// GetInt from config
func (c *Config) GetInt(section string, key string, defValue int) int {
	if c.cfg.Section(section).HasKey(key) {
		val, e := c.cfg.Section(section).Key(key).Int()
		if e == nil {
			return val
		}
	}
	return defValue
}

// GetBool from config
func (c *Config) GetBool(section string, key string, defValue bool) bool {
	if c.cfg.Section(section).HasKey(key) {
		val := c.cfg.Section(section).Key(key).String()
		val = strings.ToUpper(val)
		return val == "YES" || val == "TRUE" || val == "1"
	}
	return defValue
}
