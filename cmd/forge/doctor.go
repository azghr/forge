package main

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

func runDoctor() {
	ok := true
	ok = checkGoVersion() && ok
	ok = checkGOPATH() && ok
	ok = checkGoVet() && ok

	if ok {
		fmt.Println(green("All checks passed."))
	} else {
		fmt.Println(red("Some checks failed. See above for details."))
		os.Exit(1)
	}
}

func checkGoVersion() bool {
	v := runtime.Version()
	fmt.Printf("%-20s %s\n", cyan("Go version:"), v)
	if !strings.HasPrefix(v, "go1.26") && !strings.HasPrefix(v, "go1.27") {
		fmt.Printf("  %s Forge requires Go 1.26+\n", yellow("warning:"))
		return false
	}
	return true
}

func checkGOPATH() bool {
	gopath := os.Getenv("GOPATH")
	if gopath == "" {
		gopath = defaultGOPATH()
	}
	fmt.Printf("%-20s %s\n", cyan("GOPATH:"), gopath)
	if gopath == "" {
		fmt.Printf("  %s GOPATH is not set\n", red("error:"))
		return false
	}
	return true
}

func checkGoVet() bool {
	cmd := exec.Command("go", "vet", "./cmd/forge/...")
	cmd.Dir = findForgeRoot()
	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("%-20s %s\n", cyan("go vet:"), red("FAIL"))
		for _, line := range strings.Split(string(out), "\n") {
			if strings.TrimSpace(line) != "" {
				fmt.Printf("  %s\n", line)
			}
		}
		return false
	}
	fmt.Printf("%-20s %s\n", cyan("go vet:"), green("PASS"))
	return true
}

func defaultGOPATH() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	return home + "/go"
}

func findForgeRoot() string {
	dir, err := os.Getwd()
	if err != nil {
		return "."
	}
	for {
		if _, err := os.Stat(dir + "/go.work"); err == nil {
			return dir
		}
		idx := strings.LastIndex(dir, "/")
		if idx < 0 {
			return "."
		}
		parent := dir[:idx]
		if parent == dir {
			return "."
		}
		dir = parent
	}
}
