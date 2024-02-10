package as3_test

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/memes/f5-google-declaration-generator/pkg/generators"
	"github.com/memes/f5-google-declaration-generator/pkg/generators/as3"
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
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, as3.DefaultPackageURL, http.NoBody)
	if err != nil {
		t.Fatalf("NewRequestWithContext returned an unexpected error: %v", err)
	}
	defer http.DefaultClient.CloseIdleConnections()
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		t.Fatalf("GET %s returned an unexpected error: %v", as3.DefaultPackageURL, err)
	}
	defer response.Body.Close()
	if response.StatusCode != 200 {
		t.Fatalf("Expected status code 200, got, %d", response.StatusCode)
	}
	hash := sha256.New()
	if _, err := io.Copy(hash, response.Body); err != nil {
		t.Fatalf("io.Copy returned an unexpected error: %v", err)
	}
	expected, err := hex.DecodeString(as3.DefaultPackageSHA)
	if err != nil {
		t.Fatalf("hex.DecodeString failed to decode %s: %v", as3.DefaultPackageSHA, err)
	}
	actual := hash.Sum(nil)
	if !bytes.Equal(expected, actual) {
		t.Errorf("Expected %q, got %x", as3.DefaultPackageSHA, actual)
	}
}

func TestSchemaURL(t *testing.T) {
	t.Parallel()
	if testing.Short() {
		t.Skip("Skipping schema URL test because of short flag")
	}
	ctx, cancel := context.WithTimeout(context.TODO(), 30*time.Second)
	defer cancel()
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, as3.DefaultSchemaURL, http.NoBody)
	if err != nil {
		t.Fatalf("NewRequestWithContext returned an unexpected error: %v", err)
	}
	defer http.DefaultClient.CloseIdleConnections()
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		t.Fatalf("GET %s returned an unexpected error: %v", as3.DefaultSchemaURL, err)
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

func TestContextPrepare(t *testing.T) {
	t.Parallel()
	preparer, err := as3.NewDefaultContext().NewContextPreparer()
	if err != nil {
		t.Fatalf("unexpected failure to create default context: %v", err)
	}
	_, err = preparer(nil)
	if err != nil {
		t.Errorf("unexpected error preparing context: %v", err)
	}
}

func TestIDPattern(t *testing.T) {
	tests := []struct {
		name     string
		id       string
		expected bool
	}{
		{
			name:     "default",
			id:       uuid.NewSHA1(uuid.NameSpaceURL, []byte(as3.DefaultIDGeneratorURL)).String(),
			expected: true,
		},
		{
			name:     "empty",
			id:       "",
			expected: true,
		},
		{
			name:     "whitespace",
			id:       " \t",
			expected: false,
		},
		{
			name:     "hyphens",
			id:       "foo-bar-baz",
			expected: true,
		},
		{
			name:     "space separated",
			id:       "foo bar baz",
			expected: false,
		},
	}
	t.Parallel()
	for _, test := range tests {
		tst := test
		t.Run(tst.name, func(t *testing.T) {
			t.Parallel()
			result := as3.IDPattern.MatchString(tst.id)
			if result != tst.expected {
				t.Errorf("Expected %v, got %v", tst.expected, result)
			}
		})
	}
}

func TestLabelPattern(t *testing.T) {
	tests := []struct {
		name     string
		id       string
		expected bool
	}{
		{
			name:     "default",
			id:       as3.DefaultLabel,
			expected: true,
		},
		{
			name:     "empty",
			id:       "",
			expected: true,
		},
		{
			name:     "whitespace",
			id:       " \t",
			expected: false,
		},
		{
			name:     "hyphens",
			id:       "foo-bar-baz",
			expected: true,
		},
		{
			name:     "space separated",
			id:       "foo bar baz",
			expected: true,
		},
	}
	t.Parallel()
	for _, test := range tests {
		tst := test
		t.Run(tst.name, func(t *testing.T) {
			t.Parallel()
			result := as3.LabelPattern.MatchString(tst.id)
			if result != tst.expected {
				t.Errorf("Expected %v, got %v", tst.expected, result)
			}
		})
	}
}

func TestRemarkPattern(t *testing.T) {
	tests := []struct {
		name     string
		id       string
		expected bool
	}{
		{
			name:     "default",
			id:       as3.DefaultRemark,
			expected: true,
		},
		{
			name:     "empty",
			id:       "",
			expected: true,
		},
		{
			name:     "whitespace",
			id:       " \t",
			expected: false,
		},
		{
			name:     "spaces only",
			id:       "   ",
			expected: true,
		},
		{
			name:     "hyphens",
			id:       "foo-bar-baz",
			expected: true,
		},
		{
			name:     "space separated",
			id:       "foo bar baz",
			expected: true,
		},
	}
	t.Parallel()
	for _, test := range tests {
		tst := test
		t.Run(tst.name, func(t *testing.T) {
			t.Parallel()
			result := as3.RemarkPattern.MatchString(tst.id)
			if result != tst.expected {
				t.Errorf("Expected %v, got %v", tst.expected, result)
			}
		})
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
	return as3.DefaultSchemaURL
}

func testDefaultContext(t *testing.T) *as3.Context {
	t.Helper()
	defaultContext := as3.NewDefaultContext()
	defaultContext.Label = t.Name()
	defaultContext.Header = testHeader{
		name: t.Name(),
	}
	return defaultContext
}

func testContextPreparer(t *testing.T, newContextFn func(*testing.T) *as3.Context) generators.ContextPreparer {
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
		newContext func(t *testing.T) *as3.Context
		interfaces int
	}{
		{
			name:       "default",
			newContext: testDefaultContext,
		},
		{
			name: "with-label",
			newContext: func(t *testing.T) *as3.Context {
				t.Helper()
				ctx := testDefaultContext(t)
				ctx.Label = "test-label"
				return ctx
			},
		},
		{
			name: "with-remark",
			newContext: func(t *testing.T) *as3.Context {
				t.Helper()
				ctx := testDefaultContext(t)
				ctx.Remark = "FooBarBaz"
				return ctx
			},
		},
	}
	t.Parallel()
	t.Cleanup(http.DefaultClient.CloseIdleConnections)
	schema, err := jsonschema.Compile(as3.DefaultSchemaURL)
	if err != nil {
		t.Errorf("failed to compile JSON schema: %v", err)
	}
	for _, test := range tests {
		tst := test
		t.Run(tst.name, func(t *testing.T) {
			t.Parallel()
			var buf bytes.Buffer
			generator, err := generators.NewGenerator(
				generators.WithTemplatePath(as3.TemplatePath),
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
