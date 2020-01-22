module github.com/freeconf/examples

go 1.12

require (
	github.com/freeconf/bridge v0.0.0-20200122224605-7aeaa06ed11b
	github.com/freeconf/gconf v0.0.0-20180113115633-7483a40ddad9 // indirect
	github.com/freeconf/restconf v0.0.0-20200122125951-25e0156fb002
	github.com/freeconf/yang v0.0.0-20200122003835-a31e8a9b9760
	golang.org/x/net v0.0.0-20200114155413-6afb5195e5aa
)

// replace github.com/freeconf/yang => ../yang

// replace github.com/freeconf/restconf => ../restconf

// replace github.com/freeconf/bridge => ../bridge
