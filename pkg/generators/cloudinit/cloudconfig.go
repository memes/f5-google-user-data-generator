package cloudinit

import (
	"fmt"

	"github.com/memes/f5-google-declaration-generator/pkg/generators"
	"github.com/memes/f5-google-declaration-generator/pkg/generators/runtimeinit"
)

const (
	TemplatePath = "templates/cloud-config.yaml"
	// The default URL to the JSON schema that represents a valid cloud-config declaration.
	DefaultSchemaURL = "https://raw.githubusercontent.com/canonical/cloud-init/main/cloudinit/config/schemas/versions.schema.cloud-config.json"
)

type Context struct {
	Header      generators.Header
	SchemaURL   string
	ProxyURL    string
	RuntimeInit *runtimeinit.Context
}

// Returns a fully populated cloud-config Context.
func NewDefaultContext() *Context {
	return &Context{
		Header:      nil,
		SchemaURL:   DefaultSchemaURL,
		ProxyURL:    "",
		RuntimeInit: nil,
	}
}

func (c *Context) NewContextPreparer() (generators.ContextPreparer, error) {
	if c.RuntimeInit == nil {
		return func(_ []generators.Interface) (any, error) {
			return c, nil
		}, nil
	}
	nestedPreparer, err := c.RuntimeInit.NewContextPreparer()
	if err != nil {
		return nil, fmt.Errorf("failed to create runtime-init context preparer: %w", err)
	}
	return func(interfaces []generators.Interface) (any, error) {
		if _, err := nestedPreparer(interfaces); err != nil {
			return nil, err
		}
		return c, nil
	}, nil
}
