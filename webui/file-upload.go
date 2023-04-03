package main

import (
	"io"

	"github.com/freeconf/yang/node"
	"github.com/freeconf/yang/nodeutil"
)

// matches rpc input of uploading a book report
type upload struct {
	BookName string
	Pdf      io.Reader
}

// handle just one thing, file upload of book report
func manageUploader() node.Node {
	return &nodeutil.Basic{
		OnAction: func(r node.ActionRequest) (node.Node, error) {
			switch r.Meta.Ident() {
			case "bookReport":

				// you should do more validation of input data
				var req upload
				r.Input.UpdateInto(nodeutil.ReflectChild(&req))
				data, err := io.ReadAll(req.Pdf)
				if err != nil {
					return nil, err
				}

				// we throw away the report, but you'd probably want to keep it

				// tell user the file size as proof we got it
				resp := struct {
					FileSize uint64
				}{
					FileSize: uint64(len(data)),
				}

				return nodeutil.ReflectChild(&resp), nil
			}
			return nil, nil
		},
	}
}
