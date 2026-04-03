package main

import (
	"os"

	"github.com/tensorhero-cn/leetgpu-tester/internal/stages"
	tester_utils "github.com/tensorhero-cn/tester-utils"
)

func main() {
	definition := stages.GetDefinition()
	os.Exit(tester_utils.Run(os.Args[1:], definition))
}
