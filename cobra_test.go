package cobra

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"
	"testing"
)

// MockCheckErr mimics CheckErr without calling os.Exit
func MockCheckErr(msg interface{}) {
    if msg != nil {
        fmt.Fprintln(os.Stderr, "Error:", msg)
    }
}

func TestGPTGt(t *testing.T) {
	if !Gt(2, 1) {
		t.Errorf("Expected Gt(2, 1) to be true")
	}
	if Gt(1, 2) {
		t.Errorf("Expected Gt(1, 2) to be false")
	}
	if !Gt("3", "2") {
		t.Errorf("Expected Gt('3', '2') to be true")
	}
	if Gt("2", "3") {
		t.Errorf("Expected Gt('2', '3') to be false")
	}
	if !Gt([]int{1, 2, 3}, []int{1, 2}) {
		t.Errorf("Expected Gt([1, 2, 3], [1, 2]) to be true")
	}
	if Gt([]int{1, 2}, []int{1, 2, 3}) {
		t.Errorf("Expected Gt([1, 2], [1, 2, 3]) to be false")
	}
}

func TestGPTEq(t *testing.T) {
	if !Eq(1, 1) {
		t.Errorf("Expected Eq(1, 1) to be true")
	}
	if Eq(1, 2) {
		t.Errorf("Expected Eq(1, 2) to be false")
	}
	if !Eq("test", "test") {
		t.Errorf("Expected Eq('test', 'test') to be true")
	}
	if Eq("test", "Test") {
		t.Errorf("Expected Eq('test', 'Test') to be false")
	}
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Expected Eq to panic on unsupported types")
		}
	}()
	Eq([]int{1, 2}, []int{1, 2})
}

func TestGPTTrimRightSpace(t *testing.T) {
	if result := trimRightSpace("test   "); result != "test" {
		t.Errorf("Expected trimRightSpace('test   ') to be 'test', got '%s'", result)
	}
	if result := trimRightSpace("test\t\t\t"); result != "test" {
		t.Errorf("Expected trimRightSpace('test\\t\\t\\t') to be 'test', got '%s'", result)
	}
	if result := trimRightSpace("test\n\n"); result != "test" {
		t.Errorf("Expected trimRightSpace('test\\n\\n') to be 'test', got '%s'", result)
	}
}

func TestGPTAppendIfNotPresent(t *testing.T) {
	if result := appendIfNotPresent("test", ""); result != "test" {
		t.Errorf("Expected appendIfNotPresent('test', '') to be 'test', got '%s'", result)
	}
	if result := appendIfNotPresent("test", "append"); result != "test append" {
		t.Errorf("Expected appendIfNotPresent('test', 'append') to be 'test append', got '%s'", result)
	}
	if result := appendIfNotPresent("test append", "append"); result != "test append" {
		t.Errorf("Expected appendIfNotPresent('test append', 'append') to be 'test append', got '%s'", result)
	}
}

func TestGPTrpad(t *testing.T) {
	if result := rpad("test", 8); result != "test    " {
		t.Errorf("Expected rpad('test', 8) to be 'test    ', got '%s'", result)
	}
	if result := rpad("test", 4); result != "test" {
		t.Errorf("Expected rpad('test', 4) to be 'test', got '%s'", result)
	}
}

func TestMockCheckErr(t *testing.T) {
    // Save the original os.Stderr
    origStderr := os.Stderr
    defer func() { os.Stderr = origStderr }()

    // Create a pipe to capture os.Stderr output
    r, w, _ := os.Pipe()
    os.Stderr = w

    // Call MockCheckErr
    MockCheckErr("error")

    // Close the writer and read the captured output
    w.Close()
    var buf bytes.Buffer
    io.Copy(&buf, r)

    // Verify the output
    expectedOutput := "Error: error\n"
    if buf.String() != expectedOutput {
        t.Errorf("Expected MockCheckErr to write '%s', got '%s'", expectedOutput, buf.String())
    }
}

func TestGPTWriteStringAndCheck(t *testing.T) {
	var b strings.Builder
	WriteStringAndCheck(&b, "test")
	if result := b.String(); result != "test" {
		t.Errorf("Expected WriteStringAndCheck to write 'test', got '%s'", result)
	}
}

func TestGPTLd(t *testing.T) {
	if result := ld("test", "test", false); result != 0 {
		t.Errorf("Expected ld('test', 'test', false) to be 0, got '%d'", result)
	}
	if result := ld("test", "Test", false); result != 1 {
		t.Errorf("Expected ld('test', 'Test', false) to be 1, got '%d'", result)
	}
	if result := ld("test", "Test", true); result != 0 {
		t.Errorf("Expected ld('test', 'Test', true) to be 0, got '%d'", result)
	}
}

func TestGPTStringinSlice(t *testing.T) {
	if !stringInSlice("test", []string{"test", "example"}) {
		t.Errorf("Expected stringInSlice('test', {'test', 'example'}) to be true")
	}
	if stringInSlice("notfound", []string{"test", "example"}) {
		t.Errorf("Expected stringInSlice('notfound', {'test', 'example'}) to be false")
	}
}

func TestGPTTmpl(t *testing.T) {
	var b strings.Builder
	err := tmpl(&b, "{{.}}", "test")
	if err != nil {
		t.Errorf("Expected tmpl to execute without error, got '%s'", err)
	}
	if result := b.String(); result != "test" {
		t.Errorf("Expected tmpl to write 'test', got '%s'", result)
	}
}
