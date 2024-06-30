package main

import (
	"cine-tool/cmd"
	"cine-tool/core"
)

func init() {
	core.InitConfig()
}

func main() {
	cmd.Execute()
}
