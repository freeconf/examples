package demo

import (
	"testing"

	"github.com/freeconf/yang/node"
	"github.com/freeconf/yang/nodeutil"
	"github.com/freeconf/yang/parser"
	"github.com/freeconf/yang/source"
)

// Your Application
type SuperHero struct {
	Name string
}

type MarvelStudios struct {
	SuperHeros []*SuperHero
}

// Your Application Management
func manageStudio(studio *MarvelStudios) node.Node {
	return nodeutil.ReflectChild(studio)
}

// Application Management Unit Test
func TestUnitTest(t *testing.T) {
	// Just loading your module in your unit test verifies there is no syntax errors
	ypath := source.Dir("./yang")
	module := parser.RequireModule(ypath, "marvel-studios")

	// Create your app like you normally would
	studio := &MarvelStudios{}

	// Create a management into your app. This doesn't start a server, and it doesn't
	// even touch your app until you tell it to
	manage := node.NewBrowser(module, manageStudio(studio))

	// load a config you want to test or maybe load a config to
	// setup a different test
	cfg := nodeutil.ReadJSON(`{
		"super-heros":[{
			"name" : "spidey"
		}]
	}`)
	manage.Root().UpsertFrom(cfg)

	// you have access to your app, verify data is in there.
	if studio.SuperHeros[0].Name != "spidey" {
		t.Fail()
	}
}
