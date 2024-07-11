package cobra

import (
	"testing"
	//	flag "github.com/spf13/pflag"
	//
	// included this even tho not tested against
)

func TestGPTMarkFlagsRequiredTogether(t *testing.T) {
	c := &Command{Use: "testcmd"}
	c.Flags().String("a", "", "flag a")
	c.Flags().String("b", "", "flag b")
	c.Flags().String("c", "", "flag c")

	c.MarkFlagsRequiredTogether("a", "b", "c")

	if got, want := c.Flags().Lookup("a").Annotations[requiredAsGroupAnnotation], []string{"a b c"}; !equal(got, want) {
		t.Errorf("Expected %v, got %v", want, got)
	}
	if got, want := c.Flags().Lookup("b").Annotations[requiredAsGroupAnnotation], []string{"a b c"}; !equal(got, want) {
		t.Errorf("Expected %v, got %v", want, got)
	}
	if got, want := c.Flags().Lookup("c").Annotations[requiredAsGroupAnnotation], []string{"a b c"}; !equal(got, want) {
		t.Errorf("Expected %v, got %v", want, got)
	}
}

func TestGPTMarkFlagsOneRequired(t *testing.T) {
	c := &Command{Use: "testcmd"}
	c.Flags().String("a", "", "flag a")
	c.Flags().String("b", "", "flag b")
	c.Flags().String("c", "", "flag c")

	c.MarkFlagsOneRequired("a", "b", "c")

	if got, want := c.Flags().Lookup("a").Annotations[oneRequiredAnnotation], []string{"a b c"}; !equal(got, want) {
		t.Errorf("Expected %v, got %v", want, got)
	}
	if got, want := c.Flags().Lookup("b").Annotations[oneRequiredAnnotation], []string{"a b c"}; !equal(got, want) {
		t.Errorf("Expected %v, got %v", want, got)
	}
	if got, want := c.Flags().Lookup("c").Annotations[oneRequiredAnnotation], []string{"a b c"}; !equal(got, want) {
		t.Errorf("Expected %v, got %v", want, got)
	}
}

func TestGPTMarkFlagsMutuallyExclusive(t *testing.T) {
	c := &Command{Use: "testcmd"}
	c.Flags().String("a", "", "flag a")
	c.Flags().String("b", "", "flag b")
	c.Flags().String("c", "", "flag c")

	c.MarkFlagsMutuallyExclusive("a", "b", "c")

	if got, want := c.Flags().Lookup("a").Annotations[mutuallyExclusiveAnnotation], []string{"a b c"}; !equal(got, want) {
		t.Errorf("Expected %v, got %v", want, got)
	}
	if got, want := c.Flags().Lookup("b").Annotations[mutuallyExclusiveAnnotation], []string{"a b c"}; !equal(got, want) {
		t.Errorf("Expected %v, got %v", want, got)
	}
	if got, want := c.Flags().Lookup("c").Annotations[mutuallyExclusiveAnnotation], []string{"a b c"}; !equal(got, want) {
		t.Errorf("Expected %v, got %v", want, got)
	}
}

func TestGPTMarkFlagsRequiredTogether_Panic(t *testing.T) {
	c := &Command{Use: "testcmd"}

	// This should panic as the flag "d" is not defined
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Expected panic for undefined flag 'd'")
		}
	}()
	c.MarkFlagsRequiredTogether("d")
}

func TestGPTMarkFlagsOneRequired_Panic(t *testing.T) {
	c := &Command{Use: "testcmd"}

	// This should panic as the flag "d" is not defined
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Expected panic for undefined flag 'd'")
		}
	}()
	c.MarkFlagsOneRequired("d")
}

func TestGPTMarkFlagsMutuallyExclusive_Panic(t *testing.T) {
	c := &Command{Use: "testcmd"}

	// This should panic as the flag "d" is not defined
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Expected panic for undefined flag 'd'")
		}
	}()
	c.MarkFlagsMutuallyExclusive("d")
}

func TestGPTValidateFlagGroups(t *testing.T) {
	c := &Command{Use: "testcmd"}
	c.Flags().String("a", "", "flag a")
	c.Flags().String("b", "", "flag b")
	c.Flags().String("c", "", "flag c")

	c.MarkFlagsRequiredTogether("a", "b")
	c.MarkFlagsOneRequired("c")

	// No flags set, expect error for oneRequiredGroup
	err := c.ValidateFlagGroups()
	if err == nil || err.Error() == "" {
		t.Errorf("Expected error for oneRequiredGroup not being set")
	}

	// Setting flag "c" should clear the oneRequiredGroup error
	c.Flags().Set("c", "value")
	err = c.ValidateFlagGroups()
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Setting flag "a" without "b" should cause requiredAsGroup error
	c.Flags().Set("a", "value")
	err = c.ValidateFlagGroups()
	if err == nil || err.Error() == "" {
		t.Errorf("Expected error for requiredAsGroup not being fully set")
	}
}

func TestGPTValidateFlagGroups_MutuallyExclusive(t *testing.T) {
	c := &Command{Use: "testcmd"}
	c.Flags().String("a", "", "flag a")
	c.Flags().String("b", "", "flag b")

	c.MarkFlagsMutuallyExclusive("a", "b")

	// Setting both mutually exclusive flags should cause an error
	c.Flags().Set("a", "value")
	c.Flags().Set("b", "value")
	err := c.ValidateFlagGroups()
	if err == nil || err.Error() == "" {
		t.Errorf("Expected error for mutually exclusive flags being set")
	}
}

func TestGPTEnforceFlagGroupsForCompletion(t *testing.T) {
	c := &Command{Use: "testcmd"}
	c.Flags().String("a", "", "flag a")
	c.Flags().String("b", "", "flag b")
	c.Flags().String("c", "", "flag c")

	c.MarkFlagsRequiredTogether("a", "b")
	c.MarkFlagsOneRequired("c")
	c.MarkFlagsMutuallyExclusive("a", "c")

	// Test required together
	c.Flags().Set("a", "value")
	c.enforceFlagGroupsForCompletion()
	if c.Flags().Lookup("b").Annotations[requiredAsGroupAnnotation] == nil {
		t.Errorf("Expected flag 'b' to be marked as required")
	}

	// Test one required
	c = &Command{Use: "testcmd"}
	c.Flags().String("a", "", "flag a")
	c.Flags().String("b", "", "flag b")
	c.Flags().String("c", "", "flag c")

	c.MarkFlagsOneRequired("a", "b")
	c.enforceFlagGroupsForCompletion()
	if c.Flags().Lookup("a").Annotations[oneRequiredAnnotation] == nil {
		t.Errorf("Expected flag 'a' to be marked as required")
	}
	if c.Flags().Lookup("b").Annotations[oneRequiredAnnotation] == nil {
		t.Errorf("Expected flag 'b' to be marked as required")
	}

	// Test mutually exclusive
	c = &Command{Use: "testcmd"}
	c.Flags().String("a", "", "flag a")
	c.Flags().String("b", "", "flag b")
	c.Flags().String("c", "", "flag c")

	c.MarkFlagsMutuallyExclusive("a", "b")
	c.Flags().Set("a", "value")
	c.enforceFlagGroupsForCompletion()
	if !c.Flags().Lookup("b").Hidden {
		t.Errorf("Expected flag 'b' to be hidden")
	}
}

func equal(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}

