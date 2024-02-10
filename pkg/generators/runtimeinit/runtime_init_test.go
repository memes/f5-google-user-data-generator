package runtimeinit_test

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"io"
	"net/http"
	"regexp"
	"testing"
	"time"

	"github.com/memes/f5-google-declaration-generator/pkg/generators"
	"github.com/memes/f5-google-declaration-generator/pkg/generators/app"
	"github.com/memes/f5-google-declaration-generator/pkg/generators/as3"
	"github.com/memes/f5-google-declaration-generator/pkg/generators/cfe"
	"github.com/memes/f5-google-declaration-generator/pkg/generators/do"
	"github.com/memes/f5-google-declaration-generator/pkg/generators/runtimeinit"
	"github.com/memes/f5-google-declaration-generator/pkg/generators/ts"
	"github.com/santhosh-tekuri/jsonschema/v5"
	_ "github.com/santhosh-tekuri/jsonschema/v5/httploader"
	"go.uber.org/goleak"
	"gopkg.in/yaml.v3"
)

func TestMain(m *testing.M) {
	goleak.VerifyTestMain(m)
}

func TestPackageDownload(t *testing.T) {
	t.Parallel()
	if testing.Short() {
		t.Skip("Skipping package download test because of short flag")
	}
	ctx, cancel := context.WithTimeout(context.TODO(), 30*time.Second)
	defer cancel()
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, runtimeinit.DefaultPackageURL, http.NoBody)
	if err != nil {
		t.Fatalf("NewRequestWithContext returned an unexpected error: %v", err)
	}
	defer http.DefaultClient.CloseIdleConnections()
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		t.Fatalf("GET %s returned an unexpected error: %v", runtimeinit.DefaultPackageURL, err)
	}
	defer response.Body.Close()
	if response.StatusCode != 200 {
		t.Fatalf("Expected status code 200, got, %d", response.StatusCode)
	}
	hash := sha256.New()
	if _, err := io.Copy(hash, response.Body); err != nil {
		t.Fatalf("io.Copy returned an unexpected error: %v", err)
	}
	expected, err := hex.DecodeString(runtimeinit.DefaultPackageSHA)
	if err != nil {
		t.Fatalf("hex.DecodeString failed to decode %s: %v", runtimeinit.DefaultPackageSHA, err)
	}
	actual := hash.Sum(nil)
	if !bytes.Equal(expected, actual) {
		t.Errorf("Expected %q, got %x", runtimeinit.DefaultPackageSHA, actual)
	}
}

func TestSchemaURL(t *testing.T) {
	t.Parallel()
	if testing.Short() {
		t.Skip("Skipping schema URL test because of short flag")
	}
	ctx, cancel := context.WithTimeout(context.TODO(), 30*time.Second)
	defer cancel()
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, runtimeinit.DefaultSchemaURL, http.NoBody)
	if err != nil {
		t.Fatalf("NewRequestWithContext returned an unexpected error: %v", err)
	}
	defer http.DefaultClient.CloseIdleConnections()
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		t.Fatalf("GET %s returned an unexpected error: %v", runtimeinit.DefaultSchemaURL, err)
	}
	defer response.Body.Close()
	if response.StatusCode != 200 {
		t.Fatalf("Expected status code 200, got, %d", response.StatusCode)
	}
	buf := new(bytes.Buffer)
	if _, err := io.Copy(buf, response.Body); err != nil {
		t.Fatalf("io.Copy returned an unexpected error: %v", err)
	}
	if !json.Valid(buf.Bytes()) {
		t.Errorf("http body is not valid JSON")
	}
}

type testHeader struct {
	name string
}

func (h testHeader) Name() string {
	return h.name
}

func (h testHeader) Description() string {
	return "Test application"
}

func (h testHeader) Version() string {
	return "0.0.0"
}

func (h testHeader) Timestamp() time.Time {
	return time.Now()
}

func (h testHeader) SchemaURL() string {
	return runtimeinit.DefaultSchemaURL
}

func testDefaultContext(t *testing.T) *runtimeinit.Context {
	t.Helper()
	defaultContext := runtimeinit.NewDefaultContext()
	defaultContext.Header = testHeader{
		name: t.Name(),
	}
	return defaultContext
}

func testWithAS3Context(t *testing.T, ctx *runtimeinit.Context) *runtimeinit.Context {
	t.Helper()
	as3Context := as3.NewDefaultContext()
	as3Context.Label = t.Name()
	ctx.ApplicationServices3 = as3Context
	return ctx
}

func testWithCFEContext(t *testing.T, ctx *runtimeinit.Context) *runtimeinit.Context {
	t.Helper()
	ctx.CloudFailover = cfe.NewDefaultContext()
	return ctx
}

func testWithDOContext(t *testing.T, ctx *runtimeinit.Context) *runtimeinit.Context {
	t.Helper()
	doContext := do.NewDefaultContext()
	doContext.AdminPassword = ctx.AdminPassword
	if doContext.Licensing != nil && doContext.Licensing.LicensePool != nil {
		doContext.Licensing.LicensePool.BigIPPassword = ctx.AdminPassword
	}
	doContext.InstanceName = ctx.InstanceName
	ctx.DeclarativeOnboarding = doContext
	return ctx
}

func testWithTSContext(t *testing.T, ctx *runtimeinit.Context) *runtimeinit.Context {
	t.Helper()
	tsContext := ts.NewDefaultContext()
	tsContext.ProjectID = ctx.ProjectID
	tsContext.ServiceAccount = ctx.ServiceAccount
	ctx.TelemetryStreaming = tsContext
	return ctx
}

func testWithAppContext(t *testing.T, ctx *runtimeinit.Context) *runtimeinit.Context {
	t.Helper()
	appCtx := app.NewDefaultContext()
	appCtx.Label = t.Name()
	ctx.Application = appCtx
	return ctx
}

func testContextPreparer(t *testing.T, newContextFn func(*testing.T) *runtimeinit.Context) generators.ContextPreparer {
	t.Helper()
	ctx := newContextFn(t)
	preparer, err := ctx.NewContextPreparer()
	if err != nil {
		t.Errorf("failed to create context preparer: %v", err)
		return nil
	}
	return preparer
}

func TestGeneration(t *testing.T) {
	tests := []struct {
		name       string
		newContext func(t *testing.T) *runtimeinit.Context
		interfaces int
	}{
		{
			name:       "default",
			newContext: testDefaultContext,
		},
		{
			name: "with-full-runtime-init",
			newContext: func(t *testing.T) *runtimeinit.Context {
				t.Helper()
				return testWithAppContext(t, testWithTSContext(t, testWithDOContext(t, testWithCFEContext(t, testWithAS3Context(t, (testDefaultContext(t)))))))
			},
		},
		{
			name:       "four",
			newContext: testDefaultContext,
			interfaces: 4,
		},
		{
			name:       "five",
			newContext: testDefaultContext,
			interfaces: 5,
		},
		{
			name:       "six",
			newContext: testDefaultContext,
			interfaces: 6,
		},
		{
			name:       "seven",
			newContext: testDefaultContext,
			interfaces: 7,
		},
		{
			name:       "eight",
			newContext: testDefaultContext,
			interfaces: 8,
		},
	}
	t.Parallel()
	t.Cleanup(http.DefaultClient.CloseIdleConnections)
	schema, err := jsonschema.Compile(runtimeinit.DefaultSchemaURL)
	if err != nil {
		t.Errorf("failed to compile JSON schema: %v", err)
	}
	for _, test := range tests {
		tst := test
		t.Run(tst.name, func(t *testing.T) {
			t.Parallel()
			var buf bytes.Buffer
			generator, err := generators.NewGenerator(
				generators.WithTemplatePath(runtimeinit.TemplatePath),
				generators.WithInterfaceBuilder(runtimeinit.InterfaceBuilder),
				generators.WithContextPreparer(testContextPreparer(t, tst.newContext)),
				generators.WithWriter(&buf),
			)
			if err != nil {
				t.Errorf("failed to create a generator: %v", err)
			}
			interfaceCount := tst.interfaces
			if interfaceCount == 0 {
				interfaceCount = 3
			}
			if err := generator.Execute(interfaceCount); err != nil {
				t.Errorf("generator raised an error: %v", err)
			}
			// MTU is an int, so the handlebar declaration in template is not quoted, which causes a
			// YAML parsing and schema validation failure. Find and substitute the GCP default of 1460 for
			// all handlebar MTU fields.
			mtuPattern := regexp.MustCompile(`mtu: \{+ *[A-Z_]+_MTU *\}+`)
			sanitized := mtuPattern.ReplaceAll(buf.Bytes(), []byte("mtu: 1460"))
			var declaration any
			if err = yaml.Unmarshal(sanitized, &declaration); err != nil {
				t.Errorf("failed to unmarshal YAML: %v", err)
				t.Logf("YAML declaration:\n%v", string(sanitized))
			}
			if err := schema.Validate(declaration); err != nil {
				t.Errorf("failed to validate YAML against schema: %#v", err)
				t.Logf("YAML declaration:\n%v", string(sanitized))
				t.Logf("Mapped declaration:\n%v", declaration)
			}
		})
	}
}
