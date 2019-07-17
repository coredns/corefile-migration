package cmd

import (
	"fmt"
	"io"
	"strings"

	"github.com/coredns/corefile-migration/migration"

	"github.com/spf13/cobra"
)

// NewValidVersionsCmd represents the validversions command
func NewValidVersionsCmd(out io.Writer) *cobra.Command {
	validversionsCmd := &cobra.Command{
		Use:   "validversions",
		Short: "Shows valid versions of CoreDNS",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Fprintln(out, "The following are valid CoreDNS versions:")
			fmt.Fprintln(out, strings.Join(migration.ValidVersions(), ", "))
		},
	}
	return validversionsCmd
}
