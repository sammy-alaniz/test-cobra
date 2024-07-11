package cobra

import (
	"testing"
)

// Mock Command for testing
func newMockCommand(use string, args []string, hasSubCommands bool, validArgs []string, hasParent bool) *Command {
	cmd := &Command{
		Use:      use,
		ValidArgs: validArgs,
		Run:      emptyRun,
	}
	if hasSubCommands {
		cmd.AddCommand(&Command{Use: "sub", Run: emptyRun})
	}
	if hasParent {
		parentCmd := &Command{Use: "parent", Run: emptyRun}
		parentCmd.AddCommand(cmd)
	}
	return cmd
}

/*func emptyRun(cmd *Command, args []string) {
	// Dummy function to satisfy the Run field
}*/

func TestGPTLegacyArgs(t *testing.T) {
	tests := []struct {
		name           string
		cmd            *Command
		args           []string
		expectedOutput string
		expectError    bool
	}{
		{
			name:        "No subcommands, no error",
			cmd:         newMockCommand("test", nil, false, nil, false),
			args:        []string{"arg1"},
			expectError: false,
		},
		{
			name:        "With subcommands, no args, no error",
			cmd:         newMockCommand("test", nil, true, nil, false),
			args:        []string{},
			expectError: false,
		},
		{
			name:           "Root command with subcommands, invalid arg",
			cmd:            newMockCommand("test", nil, true, nil, false),
			args:           []string{"invalid"},
			expectedOutput: "unknown command \"invalid\" for \"test\"",
			expectError:    true,
		},
		{
			name:        "Subcommand with arbitrary args",
			cmd:         newMockCommand("test", nil, false, nil, true),
			args:        []string{"arg1"},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := legacyArgs(tt.cmd, tt.args)
			if (err != nil) != tt.expectError {
				t.Fatalf("expected error: %v, got: %v", tt.expectError, err)
			}
			if tt.expectError && err.Error() != tt.expectedOutput {
				t.Fatalf("expected output: %v, got: %v", tt.expectedOutput, err.Error())
			}
		})
	}
}

func TestGPTNoArgs(t *testing.T) {
	cmd := newMockCommand("test", nil, false, nil, false) // Ensure a non-nil Command is provided

	tests := []struct {
		name        string
		args        []string
		expectError bool
	}{
		{"No arguments", []string{}, false},
		{"With arguments", []string{"arg1"}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := NoArgs(cmd, tt.args) // Pass the command to NoArgs
			if (err != nil) != tt.expectError {
				t.Fatalf("expected error: %v, got: %v", tt.expectError, err)
			}
		})
	}
}

func TestGPTOnlyValidArgs(t *testing.T) {
	tests := []struct {
		name        string
		cmd         *Command
		args        []string
		expectError bool
	}{
		{
			name:        "Valid arguments",
			cmd:         newMockCommand("test", nil, false, []string{"arg1"}, false),
			args:        []string{"arg1"},
			expectError: false,
		},
		{
			name:        "Invalid arguments",
			cmd:         newMockCommand("test", nil, false, []string{"arg1"}, false),
			args:        []string{"invalid"},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := OnlyValidArgs(tt.cmd, tt.args)
			if (err != nil) != tt.expectError {
				t.Fatalf("expected error: %v, got: %v", tt.expectError, err)
			}
		})
	}
}

func TestGPTArbitraryArgs(t *testing.T) {
	cmd := newMockCommand("test", nil, false, nil, false)
	err := ArbitraryArgs(cmd, []string{"arg1"})
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
}

func TestGPTMinimumNArgs(t *testing.T) {
	tests := []struct {
		name        string
		minArgs     int
		args        []string
		expectError bool
	}{
		{"Minimum args met", 1, []string{"arg1"}, false},
		{"Minimum args not met", 2, []string{"arg1"}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := MinimumNArgs(tt.minArgs)(nil, tt.args)
			if (err != nil) != tt.expectError {
				t.Fatalf("expected error: %v, got: %v", tt.expectError, err)
			}
		})
	}
}

func TestGPTMaximumNArgs(t *testing.T) {
	tests := []struct {
		name        string
		maxArgs     int
		args        []string
		expectError bool
	}{
		{"Maximum args not exceeded", 2, []string{"arg1"}, false},
		{"Maximum args exceeded", 1, []string{"arg1", "arg2"}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := MaximumNArgs(tt.maxArgs)(nil, tt.args)
			if (err != nil) != tt.expectError {
				t.Fatalf("expected error: %v, got: %v", tt.expectError, err)
			}
		})
	}
}

func TestGPTExactArgs(t *testing.T) {
	tests := []struct {
		name        string
		exactArgs   int
		args        []string
		expectError bool
	}{
		{"Exact args met", 1, []string{"arg1"}, false},
		{"Exact args not met", 2, []string{"arg1"}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ExactArgs(tt.exactArgs)(nil, tt.args)
			if (err != nil) != tt.expectError {
				t.Fatalf("expected error: %v, got: %v", tt.expectError, err)
			}
		})
	}
}

func TestGPTRangeArgs(t *testing.T) {
	tests := []struct {
		name        string
		minArgs     int
		maxArgs     int
		args        []string
		expectError bool
	}{
		{"Args within range", 1, 2, []string{"arg1"}, false},
		{"Args below range", 2, 3, []string{"arg1"}, true},
		{"Args above range", 1, 2, []string{"arg1", "arg2", "arg3"}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := RangeArgs(tt.minArgs, tt.maxArgs)(nil, tt.args)
			if (err != nil) != tt.expectError {
				t.Fatalf("expected error: %v, got: %v", tt.expectError, err)
			}
		})
	}
}

func TestGPTMatchAll(t *testing.T) {
	tests := []struct {
		name        string
		cmd         *Command
		args        []string
		conditions  []PositionalArgs
		expectError bool
	}{
		{
			name:        "All conditions met",
			cmd:         newMockCommand("test", nil, false, []string{"arg1"}, false),
			args:        []string{},
			conditions:  []PositionalArgs{NoArgs, OnlyValidArgs},
			expectError: false,
		},
		{
			name:        "One condition not met",
			cmd:         newMockCommand("test", nil, false, []string{"arg1"}, false),
			args:        []string{"arg1"},
			conditions:  []PositionalArgs{ExactArgs(2)},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := MatchAll(tt.conditions...)(tt.cmd, tt.args)
			if (err != nil) != tt.expectError {
				t.Fatalf("expected error: %v, got: %v", tt.expectError, err)
			}
		})
	}
}

func TestGPTExactValidArgs(t *testing.T) {
	tests := []struct {
		name        string
		cmd         *Command
		exactArgs   int
		args        []string
		expectError bool
	}{
		{
			name:        "Exact valid args met",
			cmd:         newMockCommand("test", nil, false, []string{"arg1"}, false),
			exactArgs:   1,
			args:        []string{"arg1"},
			expectError: false,
		},
		{
			name:        "Exact valid args not met",
			cmd:         newMockCommand("test", nil, false, []string{"arg1"}, false),
			exactArgs:   2,
			args:        []string{"arg1"},
			expectError: true,
		},
		{
			name:        "Invalid args",
			cmd:         newMockCommand("test", nil, false, []string{"arg1"}, false),
			exactArgs:   1,
			args:        []string{"invalid"},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ExactValidArgs(tt.exactArgs)(tt.cmd, tt.args)
			if (err != nil) != tt.expectError {
				t.Fatalf("expected error: %v, got: %v", tt.expectError, err)
			}
		})
	}
}

