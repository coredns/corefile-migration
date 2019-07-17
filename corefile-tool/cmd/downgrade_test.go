package cmd

import (
	"bytes"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func TestNewDowngradeCmd(t *testing.T) {
	var buf bytes.Buffer

	testCorefile := "test-corefile"
	tmpDir, err := ioutil.TempDir("", "corefile")
	if err != nil {
		t.Errorf("Unable to create temporary directory: %v", err)
	}
	defer os.RemoveAll(tmpDir)
	corefilePath := filepath.Join(tmpDir, testCorefile)

	f, err := os.Create(corefilePath)
	if err != nil {
		t.Errorf("Unable to create test file %q: %v", corefilePath, err)
	}
	defer f.Close()

	testCases := []struct {
		name           string
		flags          map[string]string
		expectedOutput string
		corefile       string
		expectedError  bool
	}{
		{
			name:          "fails if no flags set",
			expectedError: true,
		},
		{
			name: "All flags set correct",
			flags: map[string]string{
				"from":     "1.5.0",
				"to":       "1.3.1",
				"corefile": corefilePath,
			},
			corefile: `.:53 {
    errors
    health
    ready
    kubernetes cluster.local in-addr.arpa ip6.arpa {
        pods insecure
        upstream
        fallthrough in-addr.arpa ip6.arpa
    }
    prometheus :9153
    forward . /etc/resolv.conf
    cache 30
    loop
    reload
    loadbalance
}
`,
			expectedOutput: `.:53 {
    errors
    health
    kubernetes cluster.local in-addr.arpa ip6.arpa {
        pods insecure
        upstream
        fallthrough in-addr.arpa ip6.arpa
    }
    prometheus :9153
    forward . /etc/resolv.conf
    cache 30
    loop
    reload
    loadbalance
}

`,
			expectedError: false,
		},
		{
			name: "flags set incorrect",
			flags: map[string]string{
				"from":     "1.4.0",
				"to":       "1.5.0",
				"corefile": corefilePath,
			},
			expectedError: true,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			cmd := NewDowngradeCmd(&buf)
			if _, err = f.WriteString(tc.corefile); err != nil {
				t.Errorf("Unable to write test file %q: %v", corefilePath, err)
			}

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
