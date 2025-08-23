package jsonnetsecure

import (
	"bytes"
	"os"
	"os/exec"
	"path"
	"runtime"
	"testing"
)

type (
	VM interface {
		EvaluateAnonymousSnippet(filename string, snippet string) (json string, formattedErr error)
		ExtCode(key string, val string)
		ExtVar(key string, val string)
		TLACode(key string, val string)
		TLAVar(key string, val string)
	}
)

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
