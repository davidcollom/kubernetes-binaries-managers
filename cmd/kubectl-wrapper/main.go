package main

import (
	"github.com/little-angry-clouds/kubernetes-binaries-managers/internal/wrapper"
)

func main() {
	var binName = "kubectl"

	wrapper.Wrapper(binName)
}
