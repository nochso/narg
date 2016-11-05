package main

import (
	"log"
	"os/exec"

	"github.com/nochso/buildutil"
)

func main() {
	log.SetFlags(0)
	run("go", "generate", "-x")
	run("goimports", "-w", ".")
	run("go", "test", "-cover")
	run("gometalinter", ".")
}

func run(name string, args ...string) {
	buildutil.ExecLog(exec.Command(name, args...))
}
