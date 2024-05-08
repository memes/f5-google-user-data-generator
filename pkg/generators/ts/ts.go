package ts

import (
	"github.com/memes/f5-google-declaration-generator/pkg/generators"
)

const (
	TemplatePath = "templates/ts.yaml"
	// The default Telemetry Streaming version to use in declarations.
	DefaultVersion = "1.32.0"
	// The default Telemetry Streaming RPM package version to use when
	// installing/referring to TS by version. This will be sliced as needed
	// to extract semver.
	DefaultPackageVersion = "1.32.0-2"
	// The default SHA256 checksum that can be used to validate a downloaded
	// TS package for the DefaultPackageVersion.
	DefaultPackageSHA = "a6bf242728a5ba1b8b8f26b59897765567db7e0f0267ba9973f822be3ab387b6"
	// The default download URL for the TS package.
	DefaultPackageURL = "https://github.com/F5Networks/f5-telemetry-streaming/releases/download/v" + DefaultVersion + "/f5-telemetry-" + DefaultPackageVersion + ".noarch.rpm"
	// The default URL to the JSON schema that represents a valid TS declaration.
	DefaultSchemaURL = "https://raw.githubusercontent.com/F5Networks/f5-telemetry-streaming/v" + DefaultVersion + "/src/schema/" + DefaultVersion + "/base_schema.json"
	// The default logging level to use for generated TS declarations.
	DefaultLogLevel = "info"
	// The default GCP Project ID.
	DefaultProjectID = "my-project-id"
	// The default GCP service account email.
	DefaultServiceAccount = "serviceAccount:bigip@my-project-id.iam.gserviceaccount.com"
)

type Context struct {
	Header         generators.Header
	Version        string
	PackageSHA     string
	PackageURL     string
	SchemaURL      string
	LogLevel       string
	ProjectID      string
	ServiceAccount string
}

// Convenience method to build and return a reference to a TS Context with appropriate default values.
func NewDefaultContext() *Context {
	return &Context{
		Header:         nil,
		Version:        DefaultVersion,
		PackageSHA:     DefaultPackageSHA,
		PackageURL:     DefaultPackageURL,
		SchemaURL:      DefaultSchemaURL,
		LogLevel:       DefaultLogLevel,
		ProjectID:      DefaultProjectID,
		ServiceAccount: DefaultServiceAccount,
	}
}

func (c *Context) NewContextPreparer() (generators.ContextPreparer, error) {
	return func(_ []generators.Interface) (any, error) {
		return c, nil
	}, nil
}
