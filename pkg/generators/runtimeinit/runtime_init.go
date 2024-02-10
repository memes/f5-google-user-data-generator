package runtimeinit

import (
	"github.com/memes/f5-google-declaration-generator/pkg/generators"
	"github.com/memes/f5-google-declaration-generator/pkg/generators/app"
	"github.com/memes/f5-google-declaration-generator/pkg/generators/as3"
	"github.com/memes/f5-google-declaration-generator/pkg/generators/cfe"
	"github.com/memes/f5-google-declaration-generator/pkg/generators/do"
	"github.com/memes/f5-google-declaration-generator/pkg/generators/ts"
)

const (
	TemplatePath = "templates/runtime-init.yaml"
	// The default runtime-init version to use in declarations.
	DefaultVersion = "1.5.2"
	// The default runtime-init RPM package version to use when
	// installing/referring to runtime-init by version. This will be sliced
	// as needed to extract semver.
	DefaultPackageVersion = "1.5.2-1"
	// The default SHA256 checksum that can be used to validate a downloaded
	// runtime-init package for the DefaultPackageVersion.
	DefaultPackageSHA = "b9eea6a7b2627343553f47d18f4ebbb2604cec38a6e761ce4b79d518ac24b2d4"
	// The default download URL for the runtime-init package.
	DefaultPackageURL = "https://github.com/F5Networks/f5-bigip-runtime-init/releases/download/" + DefaultVersion + "/f5-bigip-runtime-init-" + DefaultPackageVersion + ".gz.run"
	// The default URL to the JSON schema that represents a valid runtime-init declaration.
	DefaultSchemaURL = "https://raw.githubusercontent.com/F5Networks/f5-bigip-runtime-init/" + DefaultVersion + "/src/schema/base_schema.json"
	// The default logging level to use for generated runtime-init declarations.
	DefaultLogLevel       = "info"
	DefaultAdminPassword  = "{{{ ADMIN_PASSWORD }}}"
	DefaultProjectID      = "{{{ PROJECT_ID }}}"
	DefaultServiceAccount = "{{{ SERVICE_ACCOUNT }}}"
	DefaultInstanceName   = "{{{ INSTANCE_NAME }}}"
)

type Context struct {
	Header                generators.Header
	Version               string
	PackageSHA            string
	PackageURL            string
	SchemaURL             string
	LogLevel              string
	AdminPassword         string
	ProjectID             string
	ServiceAccount        string
	InstanceName          string
	Interfaces            []generators.Interface
	ApplicationServices3  *as3.Context
	CloudFailover         *cfe.Context
	DeclarativeOnboarding *do.Context
	TelemetryStreaming    *ts.Context
	Application           *app.Context
}

// Convenience method to build and return a reference to a runtime-init Context with appropriate default values.
func NewDefaultContext() *Context {
	return &Context{
		Header:                nil,
		Version:               DefaultVersion,
		PackageSHA:            DefaultPackageSHA,
		PackageURL:            DefaultPackageURL,
		SchemaURL:             DefaultSchemaURL,
		LogLevel:              DefaultLogLevel,
		AdminPassword:         DefaultAdminPassword,
		ProjectID:             DefaultProjectID,
		ServiceAccount:        DefaultServiceAccount,
		InstanceName:          DefaultInstanceName,
		ApplicationServices3:  nil,
		CloudFailover:         nil,
		DeclarativeOnboarding: nil,
		TelemetryStreaming:    nil,
		Application:           nil,
	}
}

func (c *Context) NewContextPreparer() (generators.ContextPreparer, error) {
	return func(interfaces []generators.Interface) (any, error) {
		c.Interfaces = interfaces
		if c.ApplicationServices3 != nil {
			c.ApplicationServices3.Interfaces = interfaces
		}
		if c.CloudFailover != nil {
			c.CloudFailover.Interfaces = interfaces
		}
		if c.DeclarativeOnboarding != nil {
			c.DeclarativeOnboarding.Interfaces = interfaces
		}
		if c.Application != nil {
			c.Application.Interfaces = interfaces
		}
		return c, nil
	}, nil
}
