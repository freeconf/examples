module github.com/freeconf/examples

go 1.12

require (
	github.com/freeconf/bridge v0.0.0-20200122224605-7aeaa06ed11b
	github.com/freeconf/restconf v0.0.0-20200122125951-25e0156fb002
	github.com/freeconf/yang v0.0.0-20200122003835-a31e8a9b9760
	golang.org/x/net v0.7.0 // indirect
)

// replace github.com/freeconf/yang => ../yang

// replace github.com/freeconf/restconf => ../restconf

// replace github.com/freeconf/bridge => ../bridge
