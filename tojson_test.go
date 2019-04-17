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
	"encoding/json"
	"encoding/xml"
	"fmt"
	"reflect"
	"testing"
)

// http://wiki.open311.org/JSON_and_XML_Conversion/
const xml0 = `
<?xml version="1.0" encoding="UTF-8" ?>
<root>

	<title>This is an example</title>
	
	<details>

		Something
	
		more
		with
		
		a lot
		of breaks!
	</details>

	<!-- this is an xml comment with < and > and ]] and [[ -->
	
	
	<table caption="a tablet with fruits">
	  <tr>
		<td>0a</td>
		<td>0b</td>
	  </tr>
       <tr>
		<td>1a</td>
		<td>1b</td>
	  </tr>

	</table>
	
	<table>
	  <name>A table desk</name>
	  <width>60</width>
	  <length>113</length>
	</table>

</root>
`

type XML0Doc struct {
	XMLName xml.Name `xml:"root"`
	Title   string   `xml:"title"`
	Details string   `xml:"details"`

	Tables []XML0Table `xml:"table"`
}

type XML0Table struct {
	XMLName xml.Name  `xml:"table"`
	Caption string    `xml:"caption,attr"`
	Rows    []XML0Row `xml:"tr"`
}
type XML0Row struct {
	XMLName xml.Name `xml:"tr"`
	Columns []string `xml:"td"`
}

func TestTransformXML0(t *testing.T) {
	xml0ParsedOriginal := parseXML(t, xml0)

	jsonML, err := ToJSON(true, bytes.NewReader([]byte(xml0)))
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(toString(jsonML))

	xmlFromJSONML, err := ToXML(jsonML)
	if err != nil {
		t.Fatal(err)
	}

	xml0ParsedJSONML := parseXML(t, xmlFromJSONML)

	if !reflect.DeepEqual(xml0ParsedOriginal, xml0ParsedJSONML) {
		t.Fatalf("expected\n%+v\n but got \n%+v\n", toString(xml0ParsedOriginal), toString(xml0ParsedJSONML))
	}

}

//==

const xml1 = `
<?xml version="1.0" encoding="UTF-8" ?>
<root xmlns:h="http://my.domain.com/hello/world">
	<h:title>
		This is an example
		<h:details>An example text with even more hello world tokens</h:details>
	</h:title>
</root>
`

func TestTransformXML1(t *testing.T) {

	jsonML, err := ToJSON(true, bytes.NewReader([]byte(xml1)))
	if err != nil {
		t.Fatal(err)
	}

	xmlFromJSONML, err := ToXML(jsonML)
	if err != nil {
		t.Fatal(err)
	}

	jsonML2, err := ToJSON(true, bytes.NewReader([]byte(xmlFromJSONML)))
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(jsonML, jsonML2) {
		t.Fatalf("expected\n%+v\n but got \n%+v\n", toString(jsonML), toString(jsonML2))
	}

}

func TestInvalidXML(t *testing.T) {
	_, err := ToJSON(true, bytes.NewReader([]byte{}))
	if err == nil {
		t.Fatal("expected err")
	}
	t.Log(err)

	_, err = ToJSON(true, bytes.NewReader([]byte("<root")))
	if err == nil {
		t.Fatal("expected err")
	}
	t.Log(err)

	_, err = ToJSON(true, bytes.NewReader([]byte("<root><child></root></child>")))
	if err == nil {
		t.Fatal("expected err")
	}
	t.Log(err)
}

//=

func parseXML(t *testing.T, xmlText string) *XML0Doc {
	xml0Struct := &XML0Doc{}
	err := xml.Unmarshal([]byte(xmlText), xml0Struct)
	if err != nil {
		t.Fatal(err)
	}
	return xml0Struct
}

func toString(arr interface{}) string {
	tmp, err := json.Marshal(arr)
	if err != nil {
		panic(err)
	}
	return string(tmp)
}
