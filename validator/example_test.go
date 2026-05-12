package validator_test

import (
	"fmt"

	"github.com/azghr/forge/validator"
)

func ExampleValidateStruct() {
	type User struct {
		Name  string `validate:"nonzero"`
		Email string `validate:"email"`
	}

	u := User{Name: "Alice", Email: "x@x"}
	if err := validator.ValidateStruct(u); err != nil {
		fmt.Println("Invalid:", err)
	} else {
		fmt.Println("Valid")
	}
	// Output: Valid
}

func ExampleValidateStruct_errors() {
	type User struct {
		Name  string `validate:"nonzero"`
		Email string `validate:"email"`
	}

	u := User{Name: "", Email: "bad"}
	err := validator.ValidateStruct(u)
	if err != nil {
		fmt.Println(err)
	}
	// Output: Name: nonzero (and 1 more errors)
}

func ExampleValidateStruct_customTag() {
	type Config struct {
		Host string `check:"nonzero"`
		Port int    `check:"nonzero"`
	}

	cfg := Config{Host: "localhost", Port: 0}
	err := validator.ValidateStruct(cfg, validator.WithTagName("check"))
	if err != nil {
		fmt.Println(err)
	}
	// Output: Port: nonzero
}
