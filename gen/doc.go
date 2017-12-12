package gen

//go:generate go run ./cmd/main.go -yangPath ../car -module car -in car.in -out car.go

// This is a handy way to integrate code generation into go build system w/o making
// to create an extra script, just run `go generate .` in this directory to generate
// all code.
