// Many YANG based solutions generate code that matches the YANG.  FreeCONF let's you
// do this by giving you full control of the parsed YANG and let's you be incontrol
// of what gets generated.

// This is a handy way to integrate code generation into go build system w/o making
// to create an extra script, just run `go generate .` in this directory to generate
// all code.

package codegeneration

//go:generate go test -run Example

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/freeconf/yang/parser"
	"github.com/freeconf/yang/source"
)

func Example_main() {
	yangPath := ".:../yang"
	moduleName := "car"
	templateName := "go_code.template"
	outName := "car/car.go"

	m := parser.RequireModule(source.Path(yangPath), moduleName)

	gen := &GoCodeGenerator{}
	err := gen.Build(m)
	if err != nil {
		panic(err)
	}

	in, err := os.Open(templateName)
	if err != nil {
		panic(err)
	}
	defer in.Close()

	out, err := os.Create(outName)
	if err != nil {
		panic(err)
	}
	defer out.Close()

	// Generate the source code
	if err := gen.Generate(in, out); err != nil {
		panic(err)
	}

	// format the code after so template doesn't have to worry about such things.
	formatCmd := exec.Command("gofmt", "-w", outName)
	formatCmd.Stderr = os.Stderr
	formatCmd.Stdout = os.Stdout
	if err = formatCmd.Run(); err != nil {
		panic(err)
	}

	// Output:
	// done
	fmt.Println("done")
}
