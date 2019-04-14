package jsonml

import (
	"bytes"
	"encoding/json"
	"fmt"
	"testing"
)

// http://wiki.open311.org/JSON_and_XML_Conversion/
const xmlX = `
<?xml version="1.0" encoding="UTF-8" ?>
<root xmlns:h="http://www.w3.org/TR/html4/"
xmlns:f="https://www.w3schools.com/furniture">

<h:table>
  <h:tr>
    <h:td>Apples</h:td>
    <h:td>Bananas</h:td>
  </h:tr>
</h:table>

<f:table>
  <f:name>African Coffee Table</f:name>
  <f:width>80</f:width>
  <f:length>120</f:length>
</f:table>

</root>
`

const xml0 = `
<ul>
	<li style="color:red">First Item</li>
	<li title="Some hover text." style="color:green">
		Second Item
	</li>
	<li>
		<span class="code-example-third">Third</span>
		Item
	</li>
</ul>
`

const json0 = `
["ul",
["li",
{ "style" : "color:red" },
"First Item"
],
["li",
{
"title" : "Some hover text.",
"style" : "color:green"
},
"Second Item"
],
["li",
["span",
{ "class" : "code-example-third" },
"Third"
],
" Item"
]
]
`

func TestTransform(t *testing.T) {

	tests := []struct {
		name string
		args string
		want string
	}{
		{"xml0", xmlX, json0},
	}

	//TODO test me properly
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			node, err := ToJSON(true, bytes.NewReader([]byte(tt.args)))
			if err != nil {
				t.Fatal(err)
			}
			//	if got := toString(node); got != tt.want {
			//		t.Errorf("Transform() = %v, want %v", got, tt.want)
			//	}

			str, err := ToXML(node)
			if err != nil {
				t.Fatal(err)
			}
			fmt.Println(str)
		})
	}
}

func toString(arr []interface{}) string {
	tmp, err := json.Marshal(arr)
	if err != nil {
		panic(err)
	}
	return string(tmp)
}
