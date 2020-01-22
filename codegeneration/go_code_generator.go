// Freeconf's opinion on code generation is that the code that get's generated
// should be under the control of the developer. That's why
// there's no specific code generation utilities built into FreeCONF.  This example
// illustrates how easy it is to generate Go code from the YANG using just the information
// from the yang parser.  While this basic example generates go code, there are a lot
// of directions you can go AND decisions you will need to make about some of the
// options in YANG that do not translate perfectly into Go code (i.e. binary types),
// or might translate into multiple things (i.e. choice)
//
// A lot of projects start off immediately with code generation and may benefit
// from using just the built in node implemtations in github.com/freeeconf/yang/nodeutil
// like node.Reflect.
//
package codegeneration

import (
	"fmt"
	"io"
	"io/ioutil"
	"strings"
	"text/template"

	"github.com/freeconf/yang/meta"
	"github.com/freeconf/yang/val"
)

// CodeGen is example of using FreeCONF's YANG parser to generate Go code
// While we could generate code from the YANG AST, it's easier to do a
// first pass on the AST and assemble the information into a set of
// data structures that make the template easier to maintain
type GoCodeGenerator struct {
	all []*objectRef // What will effectively by Go structs. all of them in a YANF file
}

// objectRef is what will be struct in go
type objectRef struct {
	Meta      meta.HasDataDefinitions
	Objects   []*objectRef
	Arrays    []*arrayRef
	Fields    []*fieldRef
	Functions []*funcRef
	Events    []*eventRef
}

// arrayRef is fields that are slices to other structs
type arrayRef struct {
	Meta *meta.List
	Item *objectRef
}

// fieldRef is details about the fields on a struct that are to primative
// types (strings, int, etc,) not other data structures
type fieldRef struct {
	Meta meta.HasType
}

// IsArray is used in the template to tell is a field is an array (i.e. leaf-list)
func (f *fieldRef) IsArray() bool {
	return f.Meta.Type().Format().IsList()
}

// funcRef are YANG rpcs or actions that will ultimately call go functions
type funcRef struct {
	Meta *meta.Rpc
}

// eventRef are YANG notifications that will call listeners (pub/sub) on objects
type eventRef struct {
	Meta *meta.Notification
}

// Build does a single pass thru the YANG AST and assemble information into a format
// that will be condusive to generating go code. Generally you pass in the *meta.Module
// but you could pass in some section of a YANG tree instead.
func (self *GoCodeGenerator) Build(m meta.HasDataDefinitions) error {
	main := &objectRef{
		Meta: m,
	}
	self.all = []*objectRef{main}
	return self.build(main)
}

func (self *GoCodeGenerator) build(p *objectRef) error {
	if x, ok := p.Meta.(meta.HasActions); ok {
		for _, y := range x.Actions() {
			ref := &funcRef{
				Meta: y,
			}
			p.Functions = append(p.Functions, ref)
		}
	}
	if x, ok := p.Meta.(meta.HasNotifications); ok {
		for _, y := range x.Notifications() {
			ref := &eventRef{
				Meta: y,
			}
			p.Events = append(p.Events, ref)
		}
	}

	for _, d := range p.Meta.DataDefinitions() {
		if meta.IsList(d) {
			ref := &arrayRef{
				Meta: d.(*meta.List),
				Item: &objectRef{Meta: d.(meta.HasDataDefinitions)},
			}
			p.Arrays = append(p.Arrays, ref)
			self.all = append(self.all, ref.Item)
			if err := self.build(ref.Item); err != nil {
				return err
			}
		} else if meta.IsLeaf(d) {
			ref := &fieldRef{
				Meta: d.(meta.HasType),
			}
			p.Fields = append(p.Fields, ref)
		} else {
			ref := &objectRef{
				Meta: d.(meta.HasDataDefinitions),
			}
			p.Objects = append(p.Objects, ref)
			self.all = append(self.all, ref)
			if err := self.build(ref); err != nil {
				return err
			}
		}
	}
	return nil
}

// Generate go code give the template in Go's template syntax.
func (self *GoCodeGenerator) Generate(in io.Reader, out io.Writer) error {
	funcMap := template.FuncMap{
		"title":     strings.Title,
		"fieldType": self.fieldType,
		"val":       self.val,
	}
	buff, err := ioutil.ReadAll(in)
	if err != nil {
		return err
	}
	t := template.Must(template.New("gen").Funcs(funcMap).Parse(string(buff)))
	return t.Execute(out, struct {
		All []*objectRef
	}{
		All: self.all,
	})
}

func (self *GoCodeGenerator) fieldType(t meta.Type) string {
	switch t.Format() {
	case val.FmtBool:
		return "bool"
	case val.FmtInt32:
		return "int"
	case val.FmtInt64:
		return "int64"
	case val.FmtString:
		return "string"
	case val.FmtDecimal64:
		return "float64"
		// TODO: many more types to cover
	}
	panic(fmt.Sprintf("unhandled format type %s", t.Format()))
}

func (self *GoCodeGenerator) val(t meta.Type) string {
	switch t.Format() {
	case val.FmtBool:
		return "Bool"
	}
	return strings.Title(t.Format().String())
}
