package template

import (
	"fmt"
	"strings"

	"github.com/machanirobotics/buffman/internal"
	"github.com/machanirobotics/buffman/internal/install"
)

// Template holds version information and comment generation
type Template struct {
	languagePrefix string
	FlatC          string
	Buffman        string
}

// NewTemplate creates a new Template instance
func NewTemplate(languagePrefix string) (*Template, error) {
	installer := install.NewManager()
	installer.RegisterInstaller(install.FlatbuffersInstaller, install.NewInstaller(install.FlatbuffersInstaller))
	flatc := installer.GetInstaller(install.FlatbuffersInstaller)
	version, err := flatc.GetVersion()
	if err != nil {
		return nil, err
	}
	return &Template{
		languagePrefix: languagePrefix,
		FlatC:          version,
		Buffman:        internal.Version,
	}, nil
}

// BuildCustomComment allows building comments with custom content
func (t *Template) BuildCustomComment(lines ...string) string {
	var builder strings.Builder

	for _, line := range lines {
		builder.WriteString(t.languagePrefix)
		builder.WriteString(" ")
		builder.WriteString(line)
		builder.WriteString("\n")
	}

	return builder.String()
}

// BuildCustomCommentf creates a single formatted comment line
func (t *Template) BuildCustomCommentf(format string, args ...any) string {
	formattedLine := fmt.Sprintf(format, args...)
	return t.BuildCustomComment(formattedLine)
}
