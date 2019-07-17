package cmd

import (
	"bytes"
	"testing"
)

func TestNewReleasedCmd(t *testing.T) {
	var buf bytes.Buffer
	testCases := []struct {
		name           string
		flags          map[string]string
		expectedOutput string
		expectedError  bool
	}{
		{
			name:          "fails if no flags set",
			expectedError: true,
		},
		{
			name: "All flags set correct",
			flags: map[string]string{
				"dockerImageSHA": "12345",
			},
			expectedOutput: `The docker image SHA is invalid
`,
			expectedError: false,
		},
		{
			name:          "flags not set",
			expectedError: true,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			cmd := NewReleasedCmd(&buf)

			// Silence the usage and errors output when testing expected errors.
			cmd.SilenceUsage = true
			cmd.SilenceErrors = true

			for f, v := range tc.flags {
				cmd.Flags().Set(f, v)
			}
			err := cmd.Execute()

			if tc.expectedError {
				if err == nil {
					t.Errorf("%s wanted err, got nil", tc.name)
				}
				return
			} else if !tc.expectedError && err != nil {
				t.Errorf("Cannot execute command: %v", err)
			}

			if buf.String() != tc.expectedOutput {
				t.Errorf("Expected output %v did not match %v", buf.String(), tc.expectedOutput)
			}
		})
	}
}
