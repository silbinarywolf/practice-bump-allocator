package inlinetest

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func TestInline(t *testing.T) {
	// note(jae): 2021-04-21
	// for seeing the cost of inline functions you can use:
	// go build -gcflags=-m=2
	cmd := exec.Command(filepath.Join(runtime.GOROOT(), "bin", "go"), "build", "-gcflags", "-m", "inlinetest.go")
	cmd.Env = os.Environ()
	outputAsBytes, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("error: %+v", err)
		return
	}
	if len(outputAsBytes) == 0 {
		t.Fatalf("error: no output")
		return
	}
	output := string(outputAsBytes)
	if !strings.Contains(output, "inlining call to bump.New") {
		t.Errorf("New is not being inlined")
	}
	if !strings.Contains(output, "inlining call to bump.(*Allocator).New") {
		t.Errorf("Allocator.New is not being inlined")
	}
}
