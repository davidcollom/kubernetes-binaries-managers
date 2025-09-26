package cmd

import (
	"fmt"
	"os"

	"github.com/hashicorp/go-version"
	. "github.com/little-angry-clouds/kubernetes-binaries-managers/internal/helpers"
	"github.com/little-angry-clouds/kubernetes-binaries-managers/internal/helpers/fzf"
	"github.com/little-angry-clouds/kubernetes-binaries-managers/internal/logging"
	. "github.com/little-angry-clouds/kubernetes-binaries-managers/internal/versions"
	"github.com/spf13/cobra"
)

func remote(cmd *cobra.Command, args []string) {
	if len(args) != 0 {
		fmt.Println("Too many arguments.")

		_ = cmd.Help()

		os.Exit(0)
	}
	var versions []*version.Version
	var err error
	var allReleases bool
	var allVersions bool
	logging.Debug("list-remote called", "args", args)

	versions, err = GetRemoteVersions(VersionsAPI)
	CheckGenericError(err)
	allReleases, err = cmd.Flags().GetBool("all-releases")
	CheckGenericError(err)
	allVersions, err = cmd.Flags().GetBool("all-versions")
	CheckGenericError(err)
	versions, err = SortVersions(versions, allReleases, allVersions)
	CheckGenericError(err)

	// Interactive selection via embedded fuzzy finder. On cancel, fall back to printing all.
	items := make([]string, 0, len(versions))
	for _, v := range versions {
		items = append(items, v.String())
	}
	if sel, err := fzf.Select(items, "Select remote version> "); err == nil && sel != "" {
		fmt.Println(sel)
		return
	} else if err == fzf.ErrNonInteractive {
		// Items already printed to stdout for piping
		return
	}

	PrintVersions(versions)
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
