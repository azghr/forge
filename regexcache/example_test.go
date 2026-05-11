package regexcache_test

import (
	"fmt"

	"github.com/azghr/forge/regexcache"
)

func ExampleCache_MustCompile() {
	cache := regexcache.New()
	r := cache.MustCompile("^a.*z$")
	fmt.Println(r.MatchString("abz"))
	// Output: true
}

func ExampleCache_Compile() {
	cache := regexcache.New()
	r, err := cache.Compile("^[a-z]+$")
	if err != nil {
		fmt.Println("invalid regex")
		return
	}
	fmt.Println(r.MatchString("hello"))
	// Output: true
}

func ExampleCache_Compile_invalid() {
	cache := regexcache.New()
	_, err := cache.Compile("(foo")
	if err != nil {
		fmt.Println("invalid regex")
	}
	// Output: invalid regex
}

func ExampleNew_withMaxSize() {
	cache := regexcache.New(regexcache.WithMaxSize(100))
	r := cache.MustCompile("^[0-9]+$")
	fmt.Println(r.MatchString("42"))
	// Output: true
}
