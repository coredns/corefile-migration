package cmd

import (
	"bytes"
	"testing"
)

func TestNewValidVersionsCmd(t *testing.T) {
	var buf bytes.Buffer
	testCases := []struct {
		name           string
		expectedOutput string
		expectedError  bool
	}{
		{
			name: "Works without error",
			expectedOutput: `The following are valid CoreDNS versions:
1.1.3, 1.1.4, 1.2.0, 1.2.1, 1.2.2, 1.2.3, 1.2.4, 1.2.5, 1.2.6, 1.3.0, 1.3.1, 1.4.0, 1.5.0, 1.5.1, 1.5.2, 1.6.0, 1.6.1, 1.6.2, 1.6.3, 1.6.4, 1.6.5, 1.6.6, 1.6.7, 1.6.9, 1.7.0, 1.7.1, 1.8.0, 1.8.3, 1.8.4
`,
			expectedError: false,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			cmd := NewValidVersionsCmd(&buf)
			err := cmd.Execute()
			if err != nil && !tc.expectedError {
				t.Errorf("Cannot execute command: %v", err)
			}

			if buf.String() != tc.expectedOutput {
				t.Errorf("Expected output %v did not match %v", buf.String(), tc.expectedOutput)
			}
		})
	}
}
