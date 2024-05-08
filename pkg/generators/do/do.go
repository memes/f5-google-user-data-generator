package do

import (
	"errors"

	"github.com/memes/f5-google-declaration-generator/pkg/generators"
)

const (
	TemplatePath = "templates/do.yaml"
	// The default DO version to use in declarations.
	DefaultVersion = "1.34.0"
	// The default DO RPM package version to use when installing/referring
	// to DO by version. This will be sliced as needed to extract semver.
	DefaultPackageVersion = "1.34.0-5"
	// The default SHA256 checksum that can be used to validate a downloaded
	// DO package for the DefaultPackageVersion.
	DefaultPackageSHA = "5e58bc15a4c436494599dfc509c87f02400339e6c0ce8275df259d5f1585146b"
	// The default download URL for DO package.
	DefaultPackageURL = "https://github.com/F5Networks/f5-declarative-onboarding/releases/download/v" + DefaultVersion + "/f5-declarative-onboarding-" + DefaultPackageVersion + ".noarch.rpm"
	// The default URL to the JSON schema that represents a valid DO declaration.
	DefaultSchemaURL = "https://raw.githubusercontent.com/F5Networks/f5-declarative-onboarding/v" + DefaultVersion + "/src/schema/" + DefaultVersion + "/base.schema.json"
	// The default label to use for generated DO declarations.
	DefaultLabel = "f5-google-declaration-generator"
	// The default hostname to assign to the instance.
	DefaultInstanceName     = "bigip"
	DefaultDomainName       = "example.com"
	GoogleMetadataService   = "169.254.169.254"
	DefaultTimezone         = "UTC"
	DefaultAdminPassword    = "Foob@r1234!"
	DefaultProvisionLTMName = "ltm"
	DefaultProvisionValue   = "nominal"
	DefaultRegKey           = "AAAAA-BBBBB-CCCCC-DDDDD-EEEEEEE"
)

type (
	RegKey struct {
		LicenseType string `yaml:"licenseType"`
		RegKey      string `yaml:"regKey"`
	}
	LicensePool struct {
		LicenseType   string `yaml:"licenseType"`
		Host          string `yaml:"bigIqHost"`
		Username      string `yaml:"bigIqUsername"`
		Password      string `yaml:"bigIqPassword"`
		PoolName      string `yaml:"licensePool"`
		BigIPUsername string `yaml:"bigIpUsername"`
		BigIPPassword string `yaml:"bigIpPassword"`
	}
	Licensing struct {
		Class       string       `yaml:"class"`
		RegKey      *RegKey      `yaml:"inline,omitempty"`
		LicensePool *LicensePool `yaml:"inline,omitempty"`
	}
	Context struct {
		Header        generators.Header
		Interfaces    []generators.Interface
		Version       string
		PackageSHA    string
		PackageURL    string
		SchemaURL     string
		Label         string
		InstanceName  string
		DomainName    string
		DNSServers    []string
		Timezone      string
		NTPServers    []string
		AdminPassword string
		SSHKeys       []string
		Provision     map[string]string
		Licensing     *Licensing
	}
)

var ErrInvalidLicensingParameter = errors.New("invalid licensing parameter")

// Convenience method to build and return a reference to a DO Context with appropriate default values.
func NewDefaultContext() *Context {
	return &Context{
		Header:       nil,
		Interfaces:   nil,
		Version:      DefaultVersion,
		PackageSHA:   DefaultPackageSHA,
		PackageURL:   DefaultPackageURL,
		SchemaURL:    DefaultSchemaURL,
		Label:        DefaultLabel,
		InstanceName: DefaultInstanceName,
		DomainName:   DefaultDomainName,
		DNSServers: []string{
			GoogleMetadataService,
		},
		Timezone: DefaultTimezone,
		NTPServers: []string{
			GoogleMetadataService,
		},
		AdminPassword: DefaultAdminPassword,
		SSHKeys:       []string{},
		Provision: map[string]string{
			DefaultProvisionLTMName: DefaultProvisionValue,
		},
		Licensing: nil,
	}
}

func (c *Context) NewContextPreparer() (generators.ContextPreparer, error) {
	return func(interfaces []generators.Interface) (any, error) {
		c.Interfaces = interfaces
		return c, nil
	}, nil
}

func (c *Context) WithRegKeyLicensing(regKey string) error {
	if regKey == "" {
		return ErrInvalidLicensingParameter
	}

	c.Licensing = &Licensing{
		Class: "License",
		RegKey: &RegKey{
			LicenseType: "regKey",
			RegKey:      regKey,
		},
		LicensePool: nil,
	}
	return nil
}

func (c *Context) WithLicensePool(poolName, bigIqHost, bigIqUsername, bigIqPassword string) error {
	c.Licensing = &Licensing{
		Class:  "License",
		RegKey: nil,
		LicensePool: &LicensePool{
			LicenseType:   "licensePool",
			Host:          bigIqHost,
			Username:      bigIqUsername,
			Password:      bigIqPassword,
			PoolName:      poolName,
			BigIPUsername: "admin",
			BigIPPassword: c.AdminPassword,
		},
	}
	return nil
}
