package cmd

import (
	"fmt"
	"io"

	"github.com/coredns/corefile-migration/migration"

	"github.com/spf13/cobra"
)

// NewReleasedCmd represents the released command
func NewReleasedCmd(out io.Writer) *cobra.Command {
	releasedCmd := &cobra.Command{
		Use:   "released",
		Short: "Determines whether your Docker Image SHA of a CoreDNS release is valid or not",
		Run: func(cmd *cobra.Command, args []string) {
			image, _ := cmd.Flags().GetString("dockerImageSHA")
			result := migration.Released(image)

			if result {
				fmt.Fprintln(out, "The docker image SHA is valid")
			} else {
				fmt.Fprintln(out, "The docker image SHA is invalid")
			}
		},
	}

	releasedCmd.Flags().String("dockerImageSHA", "", "Required: The docker image SHA you want to check. ")
	releasedCmd.MarkFlagRequired("dockerImageSHA")

	return releasedCmd
}
