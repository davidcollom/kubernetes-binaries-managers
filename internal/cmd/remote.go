package cmd

import (
	"fmt"
	"os"

	"github.com/hashicorp/go-version"
	"github.com/little-angry-clouds/kubernetes-binaries-managers/internal/helpers"
	"github.com/little-angry-clouds/kubernetes-binaries-managers/internal/helpers/fzf"
	"github.com/little-angry-clouds/kubernetes-binaries-managers/internal/logging"
	"github.com/little-angry-clouds/kubernetes-binaries-managers/internal/versions"
	"github.com/spf13/cobra"
)

func remote(cmd *cobra.Command, args []string) {
	if len(args) != 0 {
		fmt.Println("Too many arguments.")

		_ = cmd.Help()

		os.Exit(0)
	}

	var (
		versionList []*version.Version
		err         error
		allReleases bool
		allVersions bool
	)

	logging.Debug("list-remote called", "args", args)

	versionList, err = versions.GetRemoteVersions(VersionsAPI)
	helpers.CheckGenericError(err)
	allReleases, err = cmd.Flags().GetBool("all-releases")
	helpers.CheckGenericError(err)
	allVersions, err = cmd.Flags().GetBool("all-versions")
	helpers.CheckGenericError(err)
	versionList, err = versions.SortVersions(versionList, allReleases, allVersions)
	helpers.CheckGenericError(err)

	// Interactive selection via embedded fuzzy finder. On cancel, fall back to printing all.
	items := make([]string, 0, len(versionList))
	for _, v := range versionList {
		items = append(items, v.String())
	}

	if sel, err := fzf.Select(items, "Select remote version> "); err == nil && sel != "" {
		fmt.Println(sel)
		return
	} else if err == fzf.ErrNonInteractive {
		// Items already printed to stdout for piping
		return
	}

	versions.PrintVersions(versionList)
}

// remoteCmd represents the remote command
var remoteCmd = &cobra.Command{
	Use:   "remote",
	Short: "List remote versions",
	Run:   remote,
}

func init() {
	listCmd.AddCommand(remoteCmd)
}
