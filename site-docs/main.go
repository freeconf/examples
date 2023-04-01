package main

import (
	"os"
	"path/filepath"
	"text/template"
)

var cwd string

func main() {
	tmplFile := os.Args[1]
	cwd = filepath.Dir(tmplFile)
	t := template.New("sitedoc")
	t = t.Delims("[[", "]]")
	t.Funcs(template.FuncMap{
		"cp": cp,
	})
	tmplBytes, err := os.ReadFile(tmplFile)
	chkerr(err)
	t, err = t.Parse(string(tmplBytes))
	chkerr(err)
	vars := struct{}{}
	t.Execute(os.Stdout, vars)
}

func cp(fname string) string {
	content, err := os.ReadFile(filepath.Join(cwd, fname))
	chkerr(err)
	return string(content)
}

func chkerr(err error) {
	if err != nil {
		panic(err)
	}
}
