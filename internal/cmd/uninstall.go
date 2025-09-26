package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/little-angry-clouds/kubernetes-binaries-managers/internal/helpers"
	"github.com/little-angry-clouds/kubernetes-binaries-managers/internal/helpers/fzf"
	"github.com/little-angry-clouds/kubernetes-binaries-managers/internal/versions"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
)

func uninstall(cmd *cobra.Command, args []string) {
	var (
		err     error
		version string
	)

	// Support interactive selection via embedded fuzzy finder when no args provided

	switch len(args) {
	case 0:
		versionList, err := versions.GetLocalVersions(BinaryToInstall)
		helpers.CheckGenericError(err)

		if len(versionList) == 0 {
			fmt.Println("No installed versions found.")
			os.Exit(0)
		}

		versionList, err = versions.SortVersions(versionList, false, false)
		helpers.CheckGenericError(err)

		// Build string slice for selection
		items := make([]string, 0, len(versionList))
		for _, v := range versionList {
			items = append(items, v.String())
		}

		sel, err := fzf.Select(items, "Uninstall version> ")
		if err == fzf.ErrNonInteractive {
			os.Exit(0)
		}

		if err != nil || sel == "" {
			fmt.Println("No version selected.")
			os.Exit(0)
		}

		version = sel
	case 1:
		version = args[0]
	default:
		fmt.Println("Too many arguments.")

		_ = cmd.Help()

		os.Exit(0)
	}

	// Set base bin directory
	home, _ := homedir.Dir()
	fileName := fmt.Sprintf("%s/.bin/%s-v%s", home, BinaryToInstall, version)
	fileName, _ = filepath.Abs(fileName)

	if runtime.GOOS == "windows" {
		fileName += windowsSuffix
	}

	// Check if binary exists locally
	if helpers.FileExists(fileName) {
		err = os.Remove(fileName)
		helpers.CheckGenericError(err)
		fmt.Printf("Done! %s version uninstalled from %s.\n", version, fileName)
		os.Exit(0)
	}

	fmt.Printf("The version %s was already uninstalled! Doing nothing.\n", version)
}

// uninstallCmd represents the uninstall command
var uninstallCmd = &cobra.Command{
	Use:   "uninstall",
	Short: "Uninstall binary",
	Run:   uninstall,
}

func init() {
	RootCmd.AddCommand(uninstallCmd)
}
