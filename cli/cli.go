package cli

import (
	"os"
)

// Operator is the type of operation to be performed on the environment
type OperatorType string

func Operator(s string) OperatorType {
	var operators = [3]OperatorType{"new", "remove", "start"}
	for _, operator := range operators {
		if s == string(operator) {
			return operator
		}
	}
	panic("Invalid operator")
}

// Lang is the type of language to be used when creating a new environment
type LangType struct {
	Name    string `pkl:"name"`
	Version string `pkl:"version"`
}

func Lang(name string, version string) LangType {
	return LangType{name, version}
}

// Accepted args for Agnostos in order
type Args struct {
	EnvOperator OperatorType
	EnvName     string
	Lang        LangType
}

func ParseArgs() Args {
	args := Args{
		EnvOperator: Operator(os.Args[1]),
		EnvName:     os.Args[2],
	}
	if args.EnvOperator == "new" {
		args.Lang = Lang(os.Args[3], os.Args[4])
	}
	return args
}
