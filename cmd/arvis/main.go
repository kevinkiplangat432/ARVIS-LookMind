package main
// review 1: uneccessary comments to be removed later in the main branch.
import (
	"github.com/kevinkiplangat432/arvis/cmd/arvis/commands"
)

func main() {
	commands.Execute()  // calls the execute function in root.go
}