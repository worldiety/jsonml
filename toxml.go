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
	"bytes"
	"encoding/xml"
)

// set quirkyEncoder to true, to avoid bloated namespace xml by the go encoder
const quirkyEncoder = true

// ToXML creates an XML from the given compatible array. It is not defined what happens, if the array is not
// in jsonML format. Remember the format:
//
//  Element Node => json array:
//    [0] => tag name
//    [1] => optional, json object with Attribute Nodes
//    [2...] => optional, either json primitives (Text Nodes) or arrays (more children nodes)
func ToXML(jsonML []interface{}) (string, error) {
	writer := &bytes.Buffer{}
	nsList := getNamespaces(&jsonML)
	enc := xml.NewEncoder(writer)
	err := write(nsList, &jsonML, enc)
	if err != nil {
		return writer.String(), err
	}
	err = enc.Flush()
	return writer.String(), err
}

// write creates the document from scratch in a recursive manner
func write(nsList namespaces, root *jNode, enc *xml.Encoder) error {
	err := enc.EncodeToken(xml.StartElement{Name: tagName(root, quirkyEncoder, nsList), Attr: attributes(root, quirkyEncoder, nsList)})
	if err != nil {
		return err
	}
	for _, c := range children(root) {
		if cnode, ok := c.(*jNode); ok {
			err = write(nsList, cnode, enc)
			if err != nil {
				return err
			}
		} else {
			str := c.(string)
			err = enc.EncodeToken(xml.CharData(str))
			if err != nil {
				return err
			}
		}

	}
	err = enc.EncodeToken(xml.EndElement{Name: tagName(root, quirkyEncoder, nsList)})
	if err != nil {
		return err
	}
	return nil
}
