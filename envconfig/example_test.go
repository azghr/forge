package envconfig_test

import (
	"fmt"
	"os"

	"github.com/azghr/forge/envconfig"
)

func ExampleLoad() {
	type Config struct {
		Port int    `env:"PORT,default=8080"`
		Host string `env:"HOST,required"`
	}

	os.Setenv("HOST", "localhost")
	defer os.Unsetenv("HOST")

	var cfg Config
	if err := envconfig.Load(&cfg); err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(cfg.Port, cfg.Host)
	// Output: 8080 localhost
}

func ExampleLoad_missingRequired() {
	type Config struct {
		Key string `env:"KEY,required"`
	}
	os.Clearenv()
	var cfg Config
	if err := envconfig.Load(&cfg); err != nil {
		fmt.Println(err)
	}
	// Output: required env var missing: KEY
}
