package config

import (
	"strings"
)

// Ext represents a configuration file format as a string.
type Ext string

const (
	// FormatYAML represents the YAML configuration file format.
	FormatYAML Ext = "yaml"

	// FormatJSON represents the JSON configuration file format.
	FormatJSON Ext = "json"

	// FormatTOML represents the TOML configuration file format.
	FormatTOML Ext = "toml"
)

// ValidConfigFormats is a list of supported configuration file formats.
var ValidConfigFormats = []Ext{
	FormatYAML,
	FormatJSON,
	FormatTOML,
}

// IsValid checks if the Ext value matches any of the supported configuration formats in ValidConfigFormats.
func (e Ext) IsValid() bool {
	for _, valid := range ValidConfigFormats {
		if e == valid {
			return true
		}
	}
	return false
}

func getConfigFormatByExt(ext string) (Ext, error) {
	switch strings.ToLower(ext) {
	case ".yaml", ".yml":
		return FormatYAML, nil
	case ".json":
		return FormatJSON, nil
	case ".toml":
		return FormatTOML, nil
	case "":
		return "", parseError("empty file extension")
	default:
		return "", parseError("unsupported file format: %s", ext)
	}
}
