package cfe

import (
	"github.com/memes/f5-google-declaration-generator/pkg/generators"
)

const (
	TemplatePath = "templates/cfe.yaml"
	// The default CFE version to use in declarations.
	DefaultVersion = "1.13.0"
	// The default CFE RPM package version to use when installing/referring
	// to CFE by version. This will be sliced as needed to extract semver.
	DefaultPackageVersion = "1.13.0-0"
	// The default SHA256 checksum that can be used to validate a downloaded
	// CFE package for the DefaultPackageVersion.
	DefaultPackageSHA = "93be496d250838697d8a9aca8bd0e6fe7480549ecd43280279f0a63fc741ab50"
	// The default download URL for the CFE package.
	DefaultPackageURL = "https://github.com/F5Networks/f5-cloud-failover-extension/releases/download/v" + DefaultVersion + "/f5-cloud-failover-" + DefaultPackageVersion + ".noarch.rpm"
	// The default URL to the JSON schema that represents a valid CFE declaration.
	DefaultSchemaURL = "https://raw.githubusercontent.com/F5Networks/f5-cloud-failover-extension/v" + DefaultVersion + "/src/nodejs/schema/base_schema.json"
	// The default logging level for CFE declaration.
	DefaultLogLevel = "info"
)

type Context struct {
	Header                  generators.Header
	Interfaces              []generators.Interface
	Version                 string
	PackageSHA              string
	PackageURL              string
	SchemaURL               string
	LivenessHealthCheckPort string
	LogLevel                string
	ScopingTags             map[string]string
	FailoverRoutes          []FailoverRoute
}

type StaticRoute struct {
	Name             string
	NextHopAddresses []string
}

type TaggedRoute struct {
	ScopingTags   map[string]string
	AddressRanges []string
}

type FailoverRoute struct {
	StaticRoute *StaticRoute
	TaggedRoute *TaggedRoute
}

// Convenience method to build and return a reference to a CFE Context with appropriate default values.
func NewDefaultContext() *Context {
	return &Context{
		Header:     nil,
		Interfaces: nil,
		Version:    DefaultVersion,
		PackageSHA: DefaultPackageSHA,
		PackageURL: DefaultPackageURL,
		SchemaURL:  DefaultSchemaURL,
		LogLevel:   DefaultLogLevel,
		ScopingTags: map[string]string{
			"f5_cloud_failover_label": "my_deployment",
		},
		FailoverRoutes: []FailoverRoute{},
	}
}

func (c *Context) NewContextPreparer() (generators.ContextPreparer, error) {
	return func(interfaces []generators.Interface) (any, error) {
		c.Interfaces = interfaces
		return c, nil
	}, nil
}
