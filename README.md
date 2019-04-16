# jsonml [![Travis-CI](https://travis-ci.com/worldiety/jsonml.svg?branch=master)](https://travis-ci.com/worldiety/jsonml) [![Go Report Card](https://goreportcard.com/badge/github.com/worldiety/jsonml)](https://goreportcard.com/report/github.com/worldiety/jsonml) [![GoDoc](https://godoc.org/github.com/worldiety/jsonml?status.svg)](http://godoc.org/github.com/worldiety/jsonml) [![Sourcegraph](https://sourcegraph.com/github.com/worldiety/jsonml/-/badge.svg)](https://sourcegraph.com/github.com/worldiety/jsonml?badge) [![Coverage](http://gocover.io/_badge/github.com/worldiety/jsonml)](http://gocover.io/github.com/worldiety/jsonml) 
A go/golang based implementation of [jsonML](http://www.jsonml.org/), which is used to transform XML into json and back, 
which is more or less lossless. Things which are discarded are CDATA, prologs, comments and DTDs. 
Note that this is the array form and not the object form, as you can emit by the java based
json.org library (`JSONML.toJSONObject()`). To be compatible, use `JSONML.toJSONArray()`.

## jsonML example

```xml
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
```

transforms to

```json
[
   "root",
   [
      "title",
      "This is an example"
   ],
   [
      "details",
      "\n\n\t\tSomething\n\t\n\t\tmore\n\t\twith\n\t\t\n\t\ta lot\n\t\tof breaks!\n\t"
   ],
   [
      "table",
      {
         "caption":"a tablet with fruits"
      },
      [
         "tr",
         [
            "td",
            "0a"
         ],
         [
            "td",
            "0b"
         ]
      ],
      [
         "tr",
         [
            "td",
            "1a"
         ],
         [
            "td",
            "1b"
         ]
      ]
   ],
   [
      "table",
      [
         "name",
         "A table desk"
      ],
      [
         "width",
         60
      ],
      [
         "length",
         113
      ]
   ]
]
```

