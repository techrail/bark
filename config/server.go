package config

import (
	"os"
	"regexp"
	"strings"
)

var DbSchemaName, DbSchemaNameWithDot string

func init() {
	DbSchemaName = ""
	DbSchemaNameWithDot = ""
	envValue := os.Getenv("DB_SCHEMA_NAME")
	if envValue != "" {
		SetDbSchemaName(envValue)
	}
}

func SetDbSchemaName(schemaName string) {
	if strings.TrimSpace(schemaName) == "" {
		DbSchemaName = ""
		DbSchemaNameWithDot = ""
		return
	}

	// Match alphanumerics and dashes and underscores
	pattern := `^[a-zA-Z0-9_-]+$`

	re, err := regexp.Compile(pattern)

	if err != nil {
		return
	}

	// Match the input string against the regular expression
	if re.MatchString(schemaName) {
		// Input is alphanumeric with dashes and underscores
		DbSchemaName = schemaName
		DbSchemaNameWithDot = schemaName + "."
	}
}
