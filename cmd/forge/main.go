package main

import (
	"github.com/azghr/forge/flagsub"
)

func main() {
	flagsub.AddSub("doctor", "Check Go version, config validity, env issues", runDoctor)
	flagsub.AddSub("example", "Generate working example scaffold", runExample)
	flagsub.AddSub("new", "Scaffold production-ready project", runNew)
	flagsub.AddSub("add", "Add forge package to your project", runAdd)
	flagsub.Parse()
}
