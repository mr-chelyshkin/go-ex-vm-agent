package config

import (
	"strings"
)

type Ext string

const (
	FormatYAML Ext = "yaml"
	FormatJSON Ext = "json"
	FormatTOML Ext = "toml"
)

var ValidConfigFormats = []Ext{
	FormatYAML,
	FormatJSON,
	FormatTOML,
}

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
	default:
		return "", parseError("unsupported file format", ext)
	}
}
