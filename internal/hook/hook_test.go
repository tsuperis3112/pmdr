package hook

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestRunHooks(t *testing.T) {
	tests := []struct {
		name     string
		commands []string
		expected []string
	}{
		{
			name:     "no commands",
			commands: []string{},
			expected: nil,
		},
		{
			name:     "single command",
			commands: []string{"echo hook 1 called"},
			expected: []string{"hook 1 called"},
		},
		{
			name:     "multiple commands",
			commands: []string{"echo hook 1 called", "echo hook 2 called"},
			expected: []string{"hook 1 called", "hook 2 called"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tempDir := t.TempDir()
			testFile := filepath.Join(tempDir, "test_output.txt")

			// Create commands that append to the test file
			var fileCommands []string
			for _, cmd := range tt.commands {
				fileCommands = append(fileCommands, fmt.Sprintf("%s >> %s", cmd, testFile))
			}

			Run(fileCommands)

			// Give the goroutines a moment to execute and write to the file
			time.Sleep(100 * time.Millisecond)

			if len(tt.expected) == 0 {
				// If no commands were run, the file should not be created
				_, err := os.Stat(testFile)
				assert.True(t, os.IsNotExist(err), "file should not exist for no commands")
				return
			}

			content, err := os.ReadFile(testFile)
			if err != nil {
				t.Fatalf("Failed to read test output file: %v", err)
			}

			output := string(content)

			for _, expectedLine := range tt.expected {
				assert.Contains(t, output, expectedLine, "Output should contain expected line")
			}
		})
	}
}
