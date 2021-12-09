package main

import (
	"github.com/hef/cacctl/cmd"
)

//efault is `-s -w -X main.version={{.Version}} -X main.commit={{.Commit}} -X main.date={{.Date}} -X main.builtBy=goreleaser`.

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func main() {
	cmd.Version = version
	cmd.Commit = commit
	cmd.Date = date
	cmd.Execute()
}
