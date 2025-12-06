package jsonnetsecure

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path"
	"runtime"
	"testing"

	"github.com/google/go-jsonnet"
)

type (
	VM interface {
		EvaluateAnonymousSnippet(filename string, snippet string) (json string, formattedErr error)
		ExtCode(key string, val string)
		ExtVar(key string, val string)
		TLACode(key string, val string)
		TLAVar(key string, val string)
	}

	kv struct {
		Key, Value string
	}
	processParameters struct {
		Filename, Snippet                    string
		TLACodes, TLAVars, ExtCodes, ExtVars []kv
	}

	vmOptions struct {
		jsonnetBinaryPath string
		args              []string
		ctx               context.Context
		pool              *pool
	}

	Option func(o *vmOptions)
)

func (pp *processParameters) EncodeTo(w io.Writer) error {
	return json.NewEncoder(w).Encode(pp)
}

func (pp *processParameters) Decode(d []byte) error {
	return json.Unmarshal(d, pp)
}

func newVMOptions() *vmOptions {
	jsonnetBinaryPath, _ := os.Executable()
	return &vmOptions{
		jsonnetBinaryPath: jsonnetBinaryPath,
		ctx:               context.Background(),
	}
}

func WithContext(ctx context.Context) Option {
	return func(o *vmOptions) {
		o.ctx = ctx
	}
}

func WithProcessPool(p Pool) Option {
	return func(o *vmOptions) {
		pool, _ := p.(*pool)
		o.pool = pool
	}
}

func WithJsonnetBinary(jsonnetBinaryPath string) Option {
	return func(o *vmOptions) {
		o.jsonnetBinaryPath = jsonnetBinaryPath
	}
}

func WithProcessArgs(args ...string) Option {
	return func(o *vmOptions) {
		o.args = args
	}
}

func MakeSecureVM(opts ...Option) VM {
	options := newVMOptions()
	for _, o := range opts {
		o(options)
	}

	if options.pool != nil {
		return NewProcessPoolVM(options)
	} else {
		vm := jsonnet.MakeVM()
		vm.Importer(new(ErrorImporter))
		return vm
	}
}

// ErrorImporter errors when calling "import".
type ErrorImporter struct{}

// Import fetches data from a map entry.
// All paths are treated as absolute keys.
func (importer *ErrorImporter) Import(importedFrom, importedPath string) (contents jsonnet.Contents, foundAt string, err error) {
	return jsonnet.Contents{}, "", fmt.Errorf("import not available %v", importedPath)
}

func JsonnetTestBinary(t testing.TB) string {
	t.Helper()

	// We can force the usage of a given jsonnet executable.
	// Useful to test different versions, or run the tests under wine.
	if s := os.Getenv("ORY_JSONNET_PATH"); s != "" {
		return s
	}

	var stderr bytes.Buffer
	// Using `t.TempDir()` results in permissions errors on Windows, sometimes.
	outPath := path.Join(os.TempDir(), "jsonnet")
	if runtime.GOOS == "windows" {
		outPath = outPath + ".exe"
	}
	cmd := exec.Command("go", "build", "-o", outPath, "github.com/ory/x/jsonnetsecure/cmd")
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil || stderr.Len() != 0 {
		t.Fatalf("building the Go binary returned error: %v\n%s", err, stderr.String())
	}

	return outPath
}
