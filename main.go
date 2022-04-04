package main

import (
	"flag"

	"github.com/wqvoon/cbox/pkg/cmd"
	"github.com/wqvoon/cbox/pkg/flags"
	"github.com/wqvoon/cbox/pkg/rootdir"
	"github.com/wqvoon/cbox/pkg/runtime"
	runtimeUtils "github.com/wqvoon/cbox/pkg/runtime/utils"
)

func main() {
	flags.ParseAll()
	rootdir.Init()

	if runtimeUtils.IsRuntimeMode() {
		runtime.Handle()
	}

	args := flag.Args()
	cmd.Run(args[0], args[1:])
}
