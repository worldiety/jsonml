// Copyright 2019 Torben Schinke
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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
	if doc == nil {
		return nil, fmt.Errorf("start of xml not found")
	}
	return *doc, nil
}
