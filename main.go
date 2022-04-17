package main

import (
	"flag"

	"github.com/wqvoon/cbox/pkg"
	"github.com/wqvoon/cbox/pkg/cmd"
	"github.com/wqvoon/cbox/pkg/runtime"
	runtimeUtils "github.com/wqvoon/cbox/pkg/runtime/utils"
)

func main() {
	pkg.Init()

	if runtimeUtils.IsRuntimeMode() {
		runtime.Handle()
	}

	args := flag.Args()
	cmd.Run(args[0], args[1:])
}
