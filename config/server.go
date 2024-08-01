package config

import (
	"os"
	"regexp"
	"strings"
)

var (
	DbSchemaNameWithDot string
	DbSchemaName        string
	DbTableName         string
)

func init() {
	DbSchemaName = ""
	DbSchemaNameWithDot = ""
	envSchemaName := os.Getenv("DB_SCHEMA_NAME")
	if envSchemaName != "" {
		SetDbSchemaName(envSchemaName)
	}

	DbTableName = "app_log"
	envTableName := os.Getenv("DB_TABLE_NAME")
	if envTableName != "" {
		SetDbTableName(envTableName)
	}
}

func SetDbTableName(tableName string) {
	if strings.TrimSpace(tableName) == "" {
		// If the value was empty, then set to default
		DbTableName = "app_log"
		return
	}

	// Match alphanumerics and underscores
	pattern := `^[a-zA-Z0-9_]+$`

	re, err := regexp.Compile(pattern)

	if err != nil {
		return
	}

	// Match the input string against the regular expression
	if re.MatchString(tableName) {
		// Input is alphanumeric with dashes and underscores
		DbTableName = tableName
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
