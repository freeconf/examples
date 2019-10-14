module github.com/freeconf/examples

go 1.12

require (
	github.com/freeconf/bridge v0.0.0-20190910000341-add502dcae53
	github.com/freeconf/manage v0.0.0-20190928152552-c94a450e817a
	github.com/freeconf/yang v0.0.0-20191004224936-946e09ffc9b4
	golang.org/x/net v0.0.0-20190912160710-24e19bdeb0f2
)

replace github.com/freeconf/yang => ../yang

replace github.com/freeconf/manage => ../manage

replace github.com/freeconf/bridge => ../bridge
