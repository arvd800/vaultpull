package envfile

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type schemaFile struct {
	Rules map[string]struct {
		Required bool   `yaml:"required"`
		Pattern  string `yaml:"pattern"`
	} `yaml:"rules"`
}

// LoadSchema reads a YAML schema file and returns a Schema.
// Example schema file:
//
//	rules:
//	  DATABASE_URL:
//	    required: true
//	  PORT:
//	    required: false
//	    pattern: '^[0-9]+$'
func LoadSchema(path string) (Schema, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("schema: read file: %w", err)
	}

	var sf schemaFile
	if err := yaml.Unmarshal(data, &sf); err != nil {
		return nil, fmt.Errorf("schema: parse yaml: %w", err)
	}

	schema := make(Schema, len(sf.Rules))
	for k, v := range sf.Rules {
		schema[k] = SchemaRule{
			Required: v.Required,
			Pattern:  v.Pattern,
		}
	}
	return schema, nil
}
