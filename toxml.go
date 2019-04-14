package jsonml

import (
	"bytes"
	"encoding/xml"
)

// ToXML creates an XML from the given compatible array. It is not defined what happens, if the array is not
// in jsonML format. Remember the format:
//
//  Element Node => json array:
//    [0] => tag name
//    [1] => optional, json object with Attribute Nodes
//    [2...] => optional, either json primitives (Text Nodes) or arrays (more children nodes)
func ToXML(jsonML []interface{}) (string, error) {
	writer := &bytes.Buffer{}
	doc := jNode(jsonML)
	nsList := doc.namespaces()
	enc := xml.NewEncoder(writer)
	err := write(nsList, &doc, enc)
	if err != nil {
		return writer.String(), err
	}
	err = enc.Flush()
	return writer.String(), err
}

// write creates the document from scratch in a recursive manner
func write(nsList namespaces, root *jNode, enc *xml.Encoder) error {
	err := enc.EncodeToken(xml.StartElement{Name: root.tagName(nsList), Attr: root.attributes(nsList)})
	if err != nil {
		return err
	}
	for _, c := range root.children() {
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
	err = enc.EncodeToken(xml.EndElement{Name: root.tagName(nsList)})
	if err != nil {
		return err
	}
	return nil
}
