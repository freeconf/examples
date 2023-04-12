package fcinflux

import (
	"regexp"

	"github.com/freeconf/yang/node"
)

// ignoreConstraint allows for filtering when walking a FreeCONF node tree
// by regular expression.  You cannot filter in `walk` because this it really
// the destination.
type ignoreConstraint struct {
	name    string
	pattern *regexp.Regexp
}

// compileIgnores array is strings to constraint that are understood by FreeCONF
// node traversal
func compileIgnores(ignores []string) ([]ignoreConstraint, error) {
	constraints := make([]ignoreConstraint, len(ignores))
	for i, ignore := range ignores {
		r, err := regexp.Compile(ignore)
		if err != nil {
			return nil, err
		}
		constraints[i] = ignoreConstraint{
			name:    ignore,
			pattern: r,
		}
	}
	return constraints, nil
}

func (i ignoreConstraint) allow(path string) bool {
	return !i.pattern.MatchString(path)
}

// implements node.ListPreConstraint
func (i ignoreConstraint) CheckListPreConstraints(r *node.ListRequest) (bool, error) {
	return i.allow(r.Path.String()), nil
}

// implements node.ContainerPreConstraint
func (i ignoreConstraint) CheckContainerPreConstraints(r *node.ChildRequest) (bool, error) {
	return i.allow(r.Path.String()), nil
}

// implements node.FieldPreConstraint
func (i ignoreConstraint) CheckFieldPreConstraints(r *node.FieldRequest, hnd *node.ValueHandle) (bool, error) {
	return i.allow(r.Path.String()), nil
}
