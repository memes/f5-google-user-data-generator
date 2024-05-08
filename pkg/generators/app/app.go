package app

import (
	"regexp"

	"github.com/google/uuid"
	"github.com/memes/f5-google-declaration-generator/pkg/generators"
)

const (
	TemplatePath = "templates/app_as3.yaml"
	// The default AS3 version to use in declarations.
	DefaultVersion = "3.41.0"
	// The default URL to the JSON schema that represents a valid AS3 declaration.
	DefaultSchemaURL = "https://raw.githubusercontent.com/F5Networks/f5-appsvcs-extension/v" + DefaultVersion + "/schema/" + DefaultVersion + "/as3-schema.json"
	DefaultLogLevel  = "warning"
	// The default URL to use as an ID for generated AS3 declarations.
	DefaultIDGeneratorURL = "https://raw.githubusercontent.com/memes/f5-google-declaration-generator/main/pkg/generators/templates/app_as3.yaml"
	// The default label to use for generated AS3 declarations.
	DefaultLabel = "f5-google-declaration-generator"
	// The default remark to use for generated AS3 declarations.
	DefaultRemark = "Sample generated AS3"
	// The default tenant to use for generated AS3 declarations.
	DefaultTenant = "app"
	// The TCP port that will be used for readiness health check for HA.
	DefaultReadinessHealthCheckPort = "26000"
	// The default application name applied in the tenant.
	DefaultName = "app1"
)

var (
	IDPattern     = regexp.MustCompile("^[^\x00-\x20\x22'<>\x5c^`|\x7f]{0,255}$")
	LabelPattern  = regexp.MustCompile("^[^\x00-\x1f\x22#&*<>?\x5b\x5c\x5d`\x7f]{0,64}$")
	RemarkPattern = regexp.MustCompile("^[^\x00-\x1f\x22\x5c\x7f]{0,64}$")
)

type Context struct {
	Header                   generators.Header
	Interfaces               []generators.Interface
	Version                  string
	SchemaURL                string
	LogLevel                 string
	ID                       string
	Label                    string
	Remark                   string
	Tenant                   string
	VIPs                     []string
	ReadinessHealthCheckPort string
	Name                     string
}

// Convenience method to build and return a reference to an AS3 Context with appropriate default values.
func NewDefaultContext() *Context {
	return &Context{
		Header:                   nil,
		Interfaces:               nil,
		Version:                  DefaultVersion,
		SchemaURL:                DefaultSchemaURL,
		LogLevel:                 DefaultLogLevel,
		ID:                       uuid.NewSHA1(uuid.NameSpaceURL, []byte(DefaultIDGeneratorURL)).String(),
		Label:                    DefaultLabel,
		Remark:                   DefaultRemark,
		Tenant:                   DefaultTenant,
		ReadinessHealthCheckPort: DefaultReadinessHealthCheckPort,
		Name:                     DefaultName,
	}
}

func (c *Context) NewContextPreparer() (generators.ContextPreparer, error) {
	return func(interfaces []generators.Interface) (any, error) {
		c.Interfaces = interfaces
		return c, nil
	}, nil
}
