package main

import (
	"github.com/little-angry-clouds/kubernetes-binaries-managers/internal/wrapper"
)

func main() {
	var binName = "oc"

	wrapper.Wrapper(binName)
}
