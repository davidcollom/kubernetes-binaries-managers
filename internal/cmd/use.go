package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/little-angry-clouds/kubernetes-binaries-managers/internal/helpers"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
)

func use(cmd *cobra.Command, args []string) {
	var (
		err           error
		minArgsLength = 1
	)

	if len(args) == 0 {
		fmt.Println("You must specify a version!")
		os.Exit(0)
	} else if len(args) != minArgsLength {
		fmt.Println("Too many arguments.")

		_ = cmd.Help()

		os.Exit(0)
	}

	var version = args[0]

	home, _ := homedir.Dir()
	binPath := fmt.Sprintf("%s/.bin", home)
	defaultBin := fmt.Sprintf("%s/.%s-version", binPath, BinaryToInstall)
	defaultBin, _ = filepath.Abs(defaultBin)

	err = os.WriteFile(defaultBin, []byte(version), 0750) // nolint: gosec,mnd
	helpers.CheckGenericError(err)

	fmt.Printf("Done! Using %s version.\n", version)
}

// useCmd represents the use command
var useCmd = &cobra.Command{
	Use:   "use",
	Short: "Set the default version to use",
	Run:   use,
}

func init() {
	RootCmd.AddCommand(useCmd)
}
