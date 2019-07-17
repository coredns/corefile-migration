package cmd

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"github.com/lithammer/dedent"
	"github.com/spf13/cobra"
)

// CorefileTool represents the base command for the corefile-tool.
func CorefileTool(out io.Writer) *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "corefile-tool",
		Short: "A brief description of your application",
		Long: dedent.Dedent(`

			    ┌──────────────────────────────────────────────────────────┐
			    │ CoreDNS Migration Tool                                   │
			    │ Easily Migrate your Corefile                             │
			    │                                                          │
			    │ Please give us feedback at:                              │
			    │ https://github.com/coredns/corefile-migration/issues     │
			    └──────────────────────────────────────────────────────────┘

		`),
	}
	rootCmd.AddCommand(NewMigrateCmd(out))
	rootCmd.AddCommand(NewDowngradeCmd(out))
	rootCmd.AddCommand(NewDefaultCmd(out))
	rootCmd.AddCommand(NewDeprecatedCmd(out))
	rootCmd.AddCommand(NewUnsupportedCmd(out))
	rootCmd.AddCommand(NewValidVersionsCmd(out))
	rootCmd.AddCommand(NewReleasedCmd(out))

	return rootCmd
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() {
	if err := CorefileTool(os.Stdout).Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func getCorefileFromPath(corefilePath string) ([]byte, error) {
	if _, err := os.Stat(corefilePath); os.IsNotExist(err) {
		return nil, err
	}

	fileBytes, err := ioutil.ReadFile(corefilePath)
	if err != nil {
		return nil, err
	}

	return fileBytes, nil
}
