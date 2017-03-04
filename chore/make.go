package main

import (
	"log"
	"os/exec"
	"strings"
)

func main() {
	log.SetFlags(0)
	run("go", "generate", "-x", "./...")
	run("goimports", "-w", ".")
	run("go", "test", "-cover", "./...")
	run("gometalinter", ".")
}

func run(name string, args ...string) {
	cmd := exec.Command(name, args...)
	log.SetPrefix(cmd.Path + " " + strings.Join(cmd.Args, " ") + ":\n")
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("error: %s", err)
		log.Println(string(out))
	}
}
