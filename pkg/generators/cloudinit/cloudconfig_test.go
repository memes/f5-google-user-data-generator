package cloudinit_test

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/memes/f5-google-declaration-generator/pkg/generators"
	"github.com/memes/f5-google-declaration-generator/pkg/generators/app"
	"github.com/memes/f5-google-declaration-generator/pkg/generators/as3"
	"github.com/memes/f5-google-declaration-generator/pkg/generators/cfe"
	"github.com/memes/f5-google-declaration-generator/pkg/generators/cloudinit"
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

func TestSchemaURL(t *testing.T) {
	t.Parallel()
	if testing.Short() {
		t.Skip("Skipping schema URL test because of short flag")
	}
	ctx, cancel := context.WithTimeout(context.TODO(), 30*time.Second)
	defer cancel()
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, cloudinit.DefaultSchemaURL, http.NoBody)
	if err != nil {
		t.Fatalf("NewRequestWithContext returned an unexpected error: %v", err)
	}
	defer http.DefaultClient.CloseIdleConnections()
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		t.Fatalf("GET %s returned an unexpected error: %v", cloudinit.DefaultSchemaURL, err)
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
	return "Test cloudinit"
}

func (h testHeader) Version() string {
	return "0.0.0"
}

func (h testHeader) Timestamp() time.Time {
	return time.Now()
}

func (h testHeader) SchemaURL() string {
	return cloudinit.DefaultSchemaURL
}

func testDefaultContext(t *testing.T) *cloudinit.Context {
	t.Helper()
	defaultContext := cloudinit.NewDefaultContext()
	defaultContext.Header = testHeader{
		name: t.Name(),
	}
	return defaultContext
}

func testWithRuntimeInitContext(t *testing.T, ctx *cloudinit.Context) *cloudinit.Context {
	t.Helper()
	ctx.RuntimeInit = runtimeinit.NewDefaultContext()
	return ctx
}

func testWithAS3Context(t *testing.T, ctx *cloudinit.Context) *cloudinit.Context {
	t.Helper()
	if ctx.RuntimeInit == nil {
		_ = testWithRuntimeInitContext(t, ctx)
	}
	as3Context := as3.NewDefaultContext()
	as3Context.Label = t.Name()
	ctx.RuntimeInit.ApplicationServices3 = as3Context
	return ctx
}

func testWithCFEContext(t *testing.T, ctx *cloudinit.Context) *cloudinit.Context {
	t.Helper()
	if ctx.RuntimeInit == nil {
		_ = testWithRuntimeInitContext(t, ctx)
	}
	ctx.RuntimeInit.CloudFailover = cfe.NewDefaultContext()
	return ctx
}

func testWithDOContext(t *testing.T, ctx *cloudinit.Context) *cloudinit.Context {
	t.Helper()
	if ctx.RuntimeInit == nil {
		_ = testWithRuntimeInitContext(t, ctx)
	}
	doContext := do.NewDefaultContext()
	doContext.AdminPassword = ctx.RuntimeInit.AdminPassword
	if doContext.Licensing != nil && doContext.Licensing.LicensePool != nil {
		doContext.Licensing.LicensePool.BigIPPassword = ctx.RuntimeInit.AdminPassword
	}
	doContext.InstanceName = ctx.RuntimeInit.InstanceName
	ctx.RuntimeInit.DeclarativeOnboarding = doContext
	return ctx
}

func testWithTSContext(t *testing.T, ctx *cloudinit.Context) *cloudinit.Context {
	t.Helper()
	if ctx.RuntimeInit == nil {
		_ = testWithRuntimeInitContext(t, ctx)
	}
	tsContext := ts.NewDefaultContext()
	tsContext.ProjectID = ctx.RuntimeInit.ProjectID
	tsContext.ServiceAccount = ctx.RuntimeInit.ServiceAccount
	ctx.RuntimeInit.TelemetryStreaming = tsContext
	return ctx
}

func testWithAppContext(t *testing.T, ctx *cloudinit.Context) *cloudinit.Context {
	t.Helper()
	if ctx.RuntimeInit == nil {
		_ = testWithRuntimeInitContext(t, ctx)
	}
	appCtx := app.NewDefaultContext()
	appCtx.Label = t.Name()
	ctx.RuntimeInit.Application = appCtx
	return ctx
}

func testContextPreparer(t *testing.T, newContextFn func(*testing.T) *cloudinit.Context) generators.ContextPreparer {
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
		newContext func(t *testing.T) *cloudinit.Context
		interfaces int
	}{
		{
			name:       "default",
			newContext: testDefaultContext,
		},
		{
			name: "with-full-runtime-init",
			newContext: func(t *testing.T) *cloudinit.Context {
				t.Helper()
				return testWithAppContext(t, testWithTSContext(t, testWithDOContext(t, testWithCFEContext(t, testWithAS3Context(t, testDefaultContext(t))))))
			},
		},
	}
	t.Parallel()
	t.Cleanup(http.DefaultClient.CloseIdleConnections)
	schema, err := jsonschema.Compile(cloudinit.DefaultSchemaURL)
	if err != nil {
		t.Errorf("failed to compile JSON schema: %v", err)
	}
	for _, test := range tests {
		tst := test
		t.Run(tst.name, func(t *testing.T) {
			t.Parallel()
			var buf bytes.Buffer
			generator, err := generators.NewGenerator(
				generators.WithTemplatePath(cloudinit.TemplatePath),
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
			var declaration any
			if err = yaml.Unmarshal(buf.Bytes(), &declaration); err != nil {
				t.Errorf("failed to unmarshal YAML: %v", err)
				t.Logf("YAML declaration:\n%v", buf.String())
			}
			if err := schema.Validate(declaration); err != nil {
				t.Errorf("failed to validate YAML against schema: %#v", err)
				t.Logf("YAML declaration:\n%v", buf.String())
				t.Logf("Mapped declaration:\n%v", declaration)
			}
		})
	}
}
