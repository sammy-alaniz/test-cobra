package cobra

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"testing"
	"text/template"
)

func TestGPTAddTemplateFunctions(t *testing.T) {
	AddTemplateFunc("t", func() bool { return true })
	AddTemplateFuncs(template.FuncMap{
		"f": func() bool { return false },
		"h": func() string { return "Hello," },
		"w": func() string { return "world." }})

	c := &Command{}
	c.SetUsageTemplate(`{{if t}}{{h}}{{end}}{{if f}}{{h}}{{end}} {{w}}`)

	const expected = "Hello, world."
	if got := c.UsageString(); got != expected {
		t.Errorf("Expected UsageString: %v\nGot: %v", expected, got)
	}
}

func TestGPTOnInitialize(t *testing.T) {
	var initialized bool
	OnInitialize(func() { initialized = true })
	for _, initFunc := range initializers {
		initFunc()
	}
	if !initialized {
		t.Error("Expected initializer function to be called, but it was not.")
	}
}

func TestGPTOnFinalize(t *testing.T) {
	var finalized bool
	OnFinalize(func() { finalized = true })
	for _, finalizeFunc := range finalizers {
		finalizeFunc()
	}
	if !finalized {
		t.Error("Expected finalizer function to be called, but it was not.")
	}
}

func TestGPTGt(t *testing.T) {
	tests := []struct {
		a        interface{}
		b        interface{}
		expected bool
	}{
		{3, 2, true},
		{2, 3, false},
		{3, 3, false},
		{"4", "3", true},
		{"3", "4", false},
		{"4", "4", false},
		{[]int{1, 2, 3}, []int{1, 2}, true},
		{[]int{1, 2}, []int{1, 2, 3}, false},
	}

	for _, tt := range tests {
		if got := Gt(tt.a, tt.b); got != tt.expected {
			t.Errorf("Gt(%v, %v) = %v; want %v", tt.a, tt.b, got, tt.expected)
		}
	}
}

func TestGPTEq(t *testing.T) {
	tests := []struct {
		a        interface{}
		b        interface{}
		expected bool
	}{
		{3, 3, true},
		{3, 2, false},
		{"hello", "hello", true},
		{"hello", "world", false},
	}

	for _, tt := range tests {
		if got := Eq(tt.a, tt.b); got != tt.expected {
			t.Errorf("Eq(%v, %v) = %v; want %v", tt.a, tt.b, got, tt.expected)
		}
	}

	// Test for panic condition
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Expected panic for unsupported type, but did not get one")
		}
	}()
	Eq([]int{1, 2, 3}, []int{1, 2, 3})
}

func TestGPTAppendIfNotPresent(t *testing.T) {
	tests := []struct {
		s              string
		stringToAppend string
		expected       string
	}{
		{"hello", "world", "hello world"},
		{"hello world", "world", "hello world"},
		{"", "world", " world"},
		{"hello", "", "hello"},
	}

	for _, tt := range tests {
		if got := appendIfNotPresent(tt.s, tt.stringToAppend); got != tt.expected {
			t.Errorf("appendIfNotPresent(%v, %v) = %v; want %v", tt.s, tt.stringToAppend, got, tt.expected)
		}
	}
}

func TestGPTCheckErr(t *testing.T) {
	tests := []struct {
		name string
		msg  interface{}
	}{
		{"Non-nil error", errors.New("error message")},
		{"Nil error", nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.msg != nil {
				// Intercept os.Exit call
				exitCalled := false
				osExit = func(code int) {
					exitCalled = true
					panic(code)
				}

				// Capture standard error output
				oldStderr := os.Stderr
				r, w, _ := os.Pipe()
				os.Stderr = w

				defer func() {
					w.Close()
					os.Stderr = oldStderr
				}()

				defer func() {
					if r := recover(); r != nil {
						if !exitCalled {
							t.Errorf("Expected os.Exit to be called, but it was not")
						}
						if r != 1 {
							t.Errorf("Expected exit code 1, got %v", r)
						}
					} else {
						t.Errorf("Expected panic for non-nil error, but did not get one")
					}
				}()
				fmt.Println("HEREEEEE")
				CheckErr(tt.msg)
				fmt.Println("HEREEEEE")
				w.Close()
				var buf bytes.Buffer
				io.Copy(&buf, r)
				if !strings.Contains(buf.String(), "Error: error message") {
					t.Errorf("Expected 'Error: error message' in output, got %v", buf.String())
				}
			} else {
				CheckErr(tt.msg)
			}
		})
	}
}

func TestGPTWriteStringAndCheck(t *testing.T) {
	tests := []struct {
		name  string
		input string
		err   error
	}{
		{"No error", "hello", nil},
		{"With error", "hello", errors.New("error writing")},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			writer := &mockStringWriter{err: tt.err}
			if tt.err != nil {
				// Intercept os.Exit call
				exitCalled := false
				osExit = func(code int) {
					exitCalled = true
					panic(code)
				}

				// Capture standard error output
				oldStderr := os.Stderr
				r, w, _ := os.Pipe()
				os.Stderr = w

				defer func() {
					w.Close()
					os.Stderr = oldStderr
				}()

				defer func() {
					if r := recover(); r != nil {
						if !exitCalled {
							t.Errorf("Expected os.Exit to be called, but it was not")
						}
						if r != 1 {
							t.Errorf("Expected exit code 1, got %v", r)
						}
					} else {
						t.Errorf("Expected panic for non-nil error, but did not get one")
					}
				}()

				WriteStringAndCheck(writer, tt.input)
				w.Close()
				var buf bytes.Buffer
				io.Copy(&buf, r)
				if !strings.Contains(buf.String(), "Error: error writing") {
					t.Errorf("Expected 'Error: error writing' in output, got %v", buf.String())
				}
			} else {
				WriteStringAndCheck(writer, tt.input)
			}
		})
	}
}

type mockStringWriter struct {
	err error
}

func (m *mockStringWriter) WriteString(s string) (n int, err error) {
	return len(s), m.err
}

func TestEqFalseCondition(t *testing.T) {
	tests := []struct {
		a        interface{}
		b        interface{}
		expected bool
	}{
		{3, 2, false},               // Int type, false condition
		{"hello", "world", false},   // String type, false condition
		{3.14, 2.71, false},         // Unsupported type, should return false
	}

	for _, tt := range tests {
		if got := Eq(tt.a, tt.b); got != tt.expected {
			t.Errorf("Eq(%v, %v) = %v; want %v", tt.a, tt.b, got, tt.expected)
		}
	}
}

