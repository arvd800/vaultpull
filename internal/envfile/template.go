package envfile

import (
	"bytes"
	"fmt"
	"os"
	"text/template"
)

// RenderTemplate renders a Go text/template file using secrets as the data source.
// The template receives a map[string]string of secret key/value pairs.
func RenderTemplate(templatePath string, secrets map[string]string) (string, error) {
	tmplBytes, err := os.ReadFile(templatePath)
	if err != nil {
		return "", fmt.Errorf("read template: %w", err)
	}

	tmpl, err := template.New("env").Option("missingkey=error").Parse(string(tmplBytes))
	if err != nil {
		return "", fmt.Errorf("parse template: %w", err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, secrets); err != nil {
		return "", fmt.Errorf("render template: %w", err)
	}

	return buf.String(), nil
}

// RenderTemplateToFile renders a template and writes the result to destPath.
func RenderTemplateToFile(templatePath, destPath string, secrets map[string]string) error {
	output, err := RenderTemplate(templatePath, secrets)
	if err != nil {
		return err
	}
	return os.WriteFile(destPath, []byte(output), 0600)
}
