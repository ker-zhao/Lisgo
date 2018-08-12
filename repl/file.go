package repl

import (
	"os"
)

func ExecFile(path string) {
	inputFile, err := os.Open(path)
	checkErrorPanic(err, "Error: ExecFile open file failed. %s")
	defer inputFile.Close()
	Exec(inputFile, "")
}
