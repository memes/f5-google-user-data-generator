package generators

import (
	"bytes"
	"embed"
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"text/template"
	"time"
	"unicode"

	"github.com/Masterminds/sprig/v3"
	"gopkg.in/yaml.v3"
)

//go:embed templates/*
var templates embed.FS

var (
	ErrInvalidInterfaceCount        = errors.New("interface count must be an integer between 1 and 8 inclusive")
	ErrInvalidIPAddress             = errors.New("not a valid IPv4 or IPv6 address")
	VipIdentifierReplacementPattern = regexp.MustCompile("[^[:xdigit:]]+")
)

type Header interface {
	Name() string
	Description() string
	Version() string
	Timestamp() time.Time
	SchemaURL() string
}

// Defines a function that can prepare a root context object for template
// rendering from a list of Interface implementations.
type ContextPreparer func(interfaces []Interface) (any, error)

type Generator struct {
	writer           io.Writer
	interfaceBuilder InterfaceBuilder
	contextPreparer  ContextPreparer
	templatePath     string
}

type Option func(*Generator) error

func WithWriter(writer io.Writer) Option {
	return func(g *Generator) error {
		if writer != nil {
			g.writer = writer
		}
		return nil
	}
}

func WithInterfaceBuilder(builder InterfaceBuilder) Option {
	return func(g *Generator) error {
		if builder != nil {
			g.interfaceBuilder = builder
		}
		return nil
	}
}

func WithContextPreparer(preparer ContextPreparer) Option {
	return func(g *Generator) error {
		if preparer != nil {
			g.contextPreparer = preparer
		}
		return nil
	}
}

func WithTemplatePath(templatePath string) Option {
	return func(g *Generator) error {
		g.templatePath = templatePath
		return nil
	}
}

func NoopContextPreparer(_ []Interface) (any, error) {
	return struct{}{}, nil
}

func NewGenerator(options ...Option) (*Generator, error) {
	generator := &Generator{
		writer:           os.Stdout,
		interfaceBuilder: StaticInterfaceBuilder,
		contextPreparer:  NoopContextPreparer,
		templatePath:     "",
	}
	for _, option := range options {
		if err := option(generator); err != nil {
			return nil, err
		}
	}
	return generator, nil
}

// The Render function will generate a YAML file from the template and write it
// to the configured writer.
func (g Generator) Execute(count int) error {
	if count < 1 || count > 8 {
		return ErrInvalidInterfaceCount
	}
	interfaces := make([]Interface, 0, count-1)
	for index := 0; index < count; index++ {
		if index == 1 {
			continue
		}
		entry, err := g.interfaceBuilder(index)
		if err != nil {
			return err
		}
		interfaces = append(interfaces, entry)
	}
	context, err := g.contextPreparer(interfaces)
	if err != nil {
		return err
	}
	declarationTemplate := template.New(filepath.Base(g.templatePath))
	funcMap := sprig.TxtFuncMap()
	funcMap["include"] = func(template string, data any) (string, error) {
		buf := bytes.NewBuffer(nil)
		if err := declarationTemplate.ExecuteTemplate(buf, template, data); err != nil {
			return "", fmt.Errorf("failed to execute included template: %w", err)
		}
		return buf.String(), nil
	}
	funcMap["shaveMustache"] = ShaveMustache
	funcMap["chomp"] = Chomp
	funcMap["toYAML"] = ToYAML
	funcMap["vipIdentifier"] = VipIdentifier
	declarationTemplate, err = declarationTemplate.
		Funcs(funcMap).
		ParseFS(templates, g.templatePath, "templates/*")
	if err != nil {
		return fmt.Errorf("failed to parse template %s: %w", g.templatePath, err)
	}
	if err := declarationTemplate.Execute(g.writer, context); err != nil {
		return fmt.Errorf("failed to render template: %w", err)
	}
	return nil
}

func ShaveMustache(text string) string {
	return strings.TrimFunc(text, func(r rune) bool {
		return unicode.IsSpace(r) || r == '{' || r == '}'
	})
}

func Chomp(text string) string {
	return strings.TrimRight(text, "\r\n")
}

func ToYAML(obj any) (string, error) {
	if obj == nil {
		return "", nil
	}
	data, err := yaml.Marshal(obj)
	if err != nil {
		return "", fmt.Errorf("failed to marshal obj to YAML: %w", err)
	}
	return string(data), nil
}

func VipIdentifier(text string) (string, error) {
	addr, netAddr, err := net.ParseCIDR(text)
	if err != nil {
		if addr = net.ParseIP(text); addr != nil {
			err = nil
		}
	}
	if err != nil {
		return "", fmt.Errorf("failed to parse %q as IP address or CIDR: %w", text, err)
	}

	var buf strings.Builder
	buf.WriteString("vip_")
	switch {
	case netAddr != nil:
		buf.Write(VipIdentifierReplacementPattern.ReplaceAll([]byte(netAddr.String()), []byte("_")))
	case addr.To4() != nil:
		buf.Write(VipIdentifierReplacementPattern.ReplaceAll([]byte(addr.String()), []byte("_")))
		buf.WriteString("_32")
	default:
		buf.Write(VipIdentifierReplacementPattern.ReplaceAll([]byte(addr.String()), []byte("_")))
		buf.WriteString("_128")
	}
	return buf.String(), nil
}
