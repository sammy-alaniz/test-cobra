package cobra_test

import (
	"strings"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func TestGt(t *testing.T) {
	assert.True(t, cobra.Gt(2, 1))
	assert.False(t, cobra.Gt(1, 2))
	assert.True(t, cobra.Gt("3", "2"))
	assert.False(t, cobra.Gt("2", "3"))
	assert.True(t, cobra.Gt([]int{1, 2, 3}, []int{1, 2}))
	assert.False(t, cobra.Gt([]int{1, 2}, []int{1, 2, 3}))
}

func TestEq(t *testing.T) {
	assert.True(t, cobra.Eq(1, 1))
	assert.False(t, cobra.Eq(1, 2))
	assert.True(t, cobra.Eq("test", "test"))
	assert.False(t, cobra.Eq("test", "Test"))
	assert.Panics(t, func() { cobra.Eq([]int{1, 2}, []int{1, 2}) })
}

func TestTrimRightSpace(t *testing.T) {
	assert.Equal(t, "test", cobra.trimRightSpace("test   "))
	assert.Equal(t, "test", cobra.trimRightSpace("test\t\t\t"))
	assert.Equal(t, "test", cobra.trimRightSpace("test\n\n"))
}

func TestAppendIfNotPresent(t *testing.T) {
	assert.Equal(t, "test", cobra.appendIfNotPresent("test", ""))
	assert.Equal(t, "test append", cobra.appendIfNotPresent("test", "append"))
	assert.Equal(t, "test append", cobra.appendIfNotPresent("test append", "append"))
}

func TestRpad(t *testing.T) {
	assert.Equal(t, "test    ", cobra.rpad("test", 8))
	assert.Equal(t, "test", cobra.rpad("test", 4))
}

func TestCheckErr(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("CheckErr did not exit on error")
		}
	}()
	cobra.CheckErr("error")
}

func TestWriteStringAndCheck(t *testing.T) {
	var b strings.Builder
	cobra.WriteStringAndCheck(&b, "test")
	assert.Equal(t, "test", b.String())
}

func TestLd(t *testing.T) {
	assert.Equal(t, 0, cobra.ld("test", "test", false))
	assert.Equal(t, 1, cobra.ld("test", "Test", false))
	assert.Equal(t, 0, cobra.ld("test", "Test", true))
}

func TestStringInSlice(t *testing.T) {
	assert.True(t, cobra.stringInSlice("test", []string{"test", "example"}))
	assert.False(t, cobra.stringInSlice("notfound", []string{"test", "example"}))
}

func TestTmpl(t *testing.T) {
	var b strings.Builder
	err := cobra.tmpl(&b, "{{.}}", "test")
	assert.NoError(t, err)
	assert.Equal(t, "test", b.String())
}
