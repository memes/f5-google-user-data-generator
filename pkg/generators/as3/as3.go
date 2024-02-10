package as3

import (
	"errors"
	"regexp"

	"github.com/google/uuid"
	"github.com/memes/f5-google-declaration-generator/pkg/generators"
)

const (
	TemplatePath = "templates/as3.yaml"
	// The default AS3 version to use in declarations.
	DefaultVersion = "3.41.0"
	// The default AS3 RPM package version to use when installing/referring
	// to AS3 by version. This will be sliced as needed to extract semver.
	DefaultPackageVersion = "3.41.0-1"
	// The default SHA256 checksum that can be used to validate a downloaded
	// AS3 package for the DefaultPackageVersion.
	DefaultPackageSHA = "ced0948208f4dc29af7c0ea3a925a28bf8b8690a263588374e3c3d2689999490"
	// The default download URL for the AS3 package.
	DefaultPackageURL = "https://github.com/F5Networks/f5-appsvcs-extension/releases/download/v" + DefaultVersion + "/f5-appsvcs-" + DefaultPackageVersion + ".noarch.rpm"
	// The default URL to the JSON schema that represents a valid AS3 declaration.
	DefaultSchemaURL = "https://raw.githubusercontent.com/F5Networks/f5-appsvcs-extension/v" + DefaultVersion + "/schema/" + DefaultVersion + "/as3-schema.json"
	DefaultLogLevel  = "warning"
	// The TCP port that will be used for liveness health check.
	DefaultLivenessHealthCheckPort = "26000"
	// The default URL to use as an ID for generated AS3 declarations.
	DefaultIDGeneratorURL = "https://raw.githubusercontent.com/memes/f5-google-declaration-generator/main/pkg/generators/templates/as3.yaml"
	// The default label to use for generated AS3 declarations.
	DefaultLabel = "f5-google-declaration-generator"
	// The default remark to use for generated AS3 declarations.
	DefaultRemark = "Sample generated AS3"
)

var (
	IDPattern        = regexp.MustCompile("^[^\x00-\x20\x22'<>\x5c^`|\x7f]{0,255}$")
	LabelPattern     = regexp.MustCompile("^[^\x00-\x1f\x22#&*<>?\x5b\x5c\x5d`\x7f]{0,64}$")
	RemarkPattern    = regexp.MustCompile("^[^\x00-\x1f\x22\x5c\x7f]{0,64}$")
	ErrInvalidID     = errors.New("id doesn't pass AS3 validation")
	ErrInvalidLabel  = errors.New("label doesn't pass AS3 validation")
	ErrInvalidRemark = errors.New("remark doesn't pass AS3 validation")
)

func ValidateID(id string) error {
	if IDPattern.MatchString(id) {
		return nil
	}
	return ErrInvalidID
}

func ValidateLabel(label string) error {
	if LabelPattern.MatchString(label) {
		return nil
	}
	return ErrInvalidLabel
}

func ValidateRemark(remark string) error {
	if RemarkPattern.MatchString(remark) {
		return nil
	}
	return ErrInvalidRemark
}

type Context struct {
	Header                  generators.Header
	Interfaces              []generators.Interface
	Version                 string
	PackageSHA              string
	PackageURL              string
	SchemaURL               string
	LogLevel                string
	LivenessHealthCheck     bool
	LivenessHealthCheckPort string
	ID                      string
	Label                   string
	Remark                  string
}

// Convenience method to build and return a reference to an AS3 Context with appropriate default values.
func NewDefaultContext() *Context {
	return &Context{
		Header:                  nil,
		Interfaces:              nil,
		Version:                 DefaultVersion,
		PackageSHA:              DefaultPackageSHA,
		PackageURL:              DefaultPackageURL,
		SchemaURL:               DefaultSchemaURL,
		LogLevel:                DefaultLogLevel,
		LivenessHealthCheck:     true,
		LivenessHealthCheckPort: DefaultLivenessHealthCheckPort,
		ID:                      uuid.NewSHA1(uuid.NameSpaceURL, []byte(DefaultIDGeneratorURL)).String(),
		Label:                   DefaultLabel,
		Remark:                  DefaultRemark,
	}
}

func (c *Context) NewContextPreparer() (generators.ContextPreparer, error) {
	return func(interfaces []generators.Interface) (any, error) {
		c.Interfaces = interfaces
		return c, nil
	}, nil
}
