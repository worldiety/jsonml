# jsonml
A go/golang based implementation of [jsonML](http://www.jsonml.org/), which is used to transform XML into json and back, 
which is more or less lossless. Things which are discarded are CDATA, prologs, comments and DTDs. 
Note that this is the array form and not the object form, as you can emit by the java based
json.org library (`JSONML.toJSONObject()`). To be compatible, use `JSONML.toJSONArray()`.