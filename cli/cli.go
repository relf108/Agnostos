package cli

import (
	"os"
)

func ReadArgs() {
	for i, arg := range os.Args {
		println(i, arg)
	}
}
