package cmd

import (
	"fmt"
	"os"

	version "github.com/hashicorp/go-version"
	"github.com/little-angry-clouds/kubernetes-binaries-managers/internal/helpers"
	vers "github.com/little-angry-clouds/kubernetes-binaries-managers/internal/versions"
	"github.com/spf13/cobra"
)

func local(cmd *cobra.Command, args []string) {
	var (
		err         error
		allReleases bool
		allVersions bool
		versions    []*version.Version
	)

	if len(args) != 0 {
		fmt.Println("Too many arguments.")

		_ = cmd.Help()

		os.Exit(0)
	}

	allReleases, err = cmd.Flags().GetBool("all-releases")

	helpers.CheckGenericError(err)

	allVersions, err = cmd.Flags().GetBool("all-versions")

	helpers.CheckGenericError(err)

	versions, err = vers.GetLocalVersions(BinaryToInstall)

	helpers.CheckGenericError(err)

	versions, err = vers.SortVersions(versions, allReleases, allVersions)

	helpers.CheckGenericError(err)

	vers.PrintVersions(versions)
}

// localCmd represents the local command
var localCmd = &cobra.Command{
	Use:   "local",
	Short: "List installed versions",
	Run:   local,
}

func init() {
	listCmd.AddCommand(localCmd)
}
