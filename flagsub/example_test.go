package flagsub_test

import (
	"fmt"

	"github.com/azghr/forge/flagsub"
)

func ExampleAddSub() {
	flagsub.Reset()

	var serve *flagsub.Sub
	serve = flagsub.AddSub("serve", "Start the server", func() {
		port := serve.Flags.Int("port", 8080, "listen port")
		serve.Flags.Parse([]string{"--port=9090"})
		fmt.Println(*port)
	})

	flagsub.ParseArgs([]string{serve.Name, "--port=9090"})
	// Output: 9090
}

func ExampleParseArgs_unknown() {
	flagsub.Reset()

	flagsub.AddSub("foo", "Foo command", func() {})

	err := flagsub.ParseArgs([]string{"bogus"})
	fmt.Println(err)
	// Output: flagsub: unknown subcommand "bogus"
}
