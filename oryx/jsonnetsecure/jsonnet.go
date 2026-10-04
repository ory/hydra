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
	"github.com/pkg/errors"
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

// ErrNoProcessPool is returned by MakeSecureVM when called without a process
// pool. It is a distinct error because the alternative — quietly returning an
// in-process VM — would strip the isolation callers of this package rely on.
var ErrNoProcessPool = errors.New("jsonnetsecure: a process pool is required; use MakeInProcessVM to evaluate in this process without isolation")

// MakeSecureVM returns a VM that evaluates snippets in a worker process taken
// from p, so that a snippet which exhausts memory, spins on the CPU, or crashes
// takes down only that worker.
//
// p is a required argument rather than an option because a VM without a pool
// offers no isolation at all. Passing a nil pool returns ErrNoProcessPool.
func MakeSecureVM(p Pool, opts ...Option) (VM, error) {
	// A nil *pool inside a non-nil Pool interface is not reachable from outside
	// this package (Pool has an unexported method), but check the concrete
	// value anyway so a future in-package mistake cannot slip through.
	concrete, _ := p.(*pool)
	if p == nil || concrete == nil {
		return nil, errors.WithStack(ErrNoProcessPool)
	}

	options := newVMOptions()
	for _, o := range opts {
		o(options)
	}
	options.pool = concrete

	return NewProcessPoolVM(options), nil
}

// MakeInProcessVM returns a Jsonnet VM that evaluates in the calling process
// with imports disabled. It provides no isolation: a malicious or buggy snippet
// can exhaust this process's memory and CPU, so it is only safe for trusted
// input. The two legitimate uses are the jsonnet subcommand, which is itself
// the isolation boundary, and offline CLI linting. Everything that evaluates
// tenant-supplied Jsonnet must use MakeSecureVM.
func MakeInProcessVM() *jsonnet.VM {
	vm := jsonnet.MakeVM()
	vm.Importer(new(ErrorImporter))
	return vm
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
