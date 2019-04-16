package jsonml

import (
	"encoding/xml"
	"fmt"
	"io"
	"strings"
)

// ToJSON takes an XML string and converts it into a json array. The trim flag indicates if whitespaces should
// be stripped. If you want to trim depends on the use case. Often there is no meaning for trailing white space,
// as defined by Unicode, so you likely want to set trim to true.
func ToJSON(trim bool, r io.Reader) ([]interface{}, error) {

	nsList := namespaces{}
	var doc *jNode
	stack := nodeStack{}

	d := xml.NewDecoder(r)
	for {
		t, tokenErr := d.Token()
		if tokenErr != nil {
			if tokenErr == io.EOF {
				break
			}
			return nil, fmt.Errorf("ToJSON failed to transform xml: %v", tokenErr)
		}

		switch t := t.(type) {
		case xml.ProcInst:

		case xml.CharData:
			//discard char data before root element
			if len(stack) > 0 {
				if trim {
					trim := strings.TrimSpace(string(t))
					if len(trim) > 0 {
						// caution: we do NOT trim real text nodes by intention
						stack.top().addPrimitive(string(t))
					}
				} else {
					stack.top().addPrimitive(string(t))
				}

			}

		case xml.StartElement:

			if len(stack) == 0 {

				// import global defined namespaces
				nsList.register("xmlns", "xmlns") //bootstrap xmlns namespace
				for _, a := range t.Attr {
					if a.Name.Space == "xmlns" {
						nsList.register(a.Name.Local, a.Value)
					}
				}

				// bootstrap document
				doc = newJNode(nsList, t.Name, t.Attr)
				stack.push(doc)
			} else {
				node := stack.top().addNode(nsList, t.Name, t.Attr)
				stack.push(node)
			}
		case xml.EndElement:
			stack.pop()
		}

	}

	return *doc, nil
}
