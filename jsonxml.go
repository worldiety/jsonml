package jsonml

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"strconv"
	"strings"
)

// namespaces just tracks url -> abbreviation/prefix
type namespaces map[string]string

// tagName resolves a jsonML conforming attribute name
func (ns namespaces) tagName(name xml.Name) string {
	if len(name.Space) == 0 {
		return name.Local
	}
	return ns.prefix(name.Space) + ":" + name.Local
}

// register is used to put globally valid
func (ns namespaces) register(prefix string, url string) {
	ns[url] = prefix
}

// prefix returns a prefix for the given url and may generate a unique namespace itself, it not yet available
func (ns namespaces) prefix(url string) string {
	if len(url) == 0 {
		return ""
	}
	val, ok := ns[url]
	if !ok {
		for i := 'a'; i <= 'z'; i++ {
			key := string(i)
			_, exists := ns[key]
			if !exists {
				ns[url] = key
				return key
			}
		}
		panic("to many namespaces")
	}
	return val
}

func (ns namespaces) reverse(prefix string) string {
	for k, v := range ns {
		if v == prefix {
			return k
		}
	}
	return ""
}

//==

// A jNode is just an array of strings or objects or arrays
type jNode []interface{}

// newJNode is a factory function
func newJNode(nsList namespaces, name xml.Name, attributes []xml.Attr) *jNode {
	tagName := nsList.tagName(name)
	node := jNode{}
	node = append(node, tagName)
	if len(attributes) > 0 {
		attrs := make(map[string]interface{})
		for _, attr := range attributes {
			attrs[nsList.tagName(attr.Name)] = attr.Value
		}
		node = append(node, attrs)
	}
	return &node
}

// namespaces parses the xmlns:x attributes
func (n *jNode) namespaces() namespaces {
	nsList := namespaces{}
	for k, v := range n.attributesAsMap() {
		if strings.HasPrefix(k, "xmlns:") {
			prefix := k[6:]
			nsList.register(prefix, v)
		}
	}
	return nsList
}

// addNode inserts a node into the actual node (which is an array in array)
func (n *jNode) addNode(nsList namespaces, name xml.Name, attributes []xml.Attr) *jNode {
	node := newJNode(nsList, name, attributes)
	*n = append(*n, node)
	return node
}

// addText is the same as addNode but just takes the string and duck types it into json
func (n *jNode) addPrimitive(text string) {
	*n = append(*n, nicefy(text))
}

// String marshals to json
func (n *jNode) String() string {
	tmp, err := json.Marshal(n)
	if err != nil {
		panic(err)
	}
	return string(tmp)
}

// tagName returns the name of the xml tag, which is always at index 0
func (n *jNode) tagName(nsList namespaces) xml.Name {
	if len(*n) == 0 {
		return xml.Name{Space: "", Local: "unknown"}
	}
	str := stringOf((*n)[0])
	return parseXMLName(nsList, str)
}

func parseXMLName(nsList namespaces, str string) xml.Name {
	if strings.Contains(str, ":") {
		tokens := strings.Split(str, ":")
		fullUrl := nsList.reverse(tokens[0])
		if fullUrl == "" {
			return xml.Name{Space: "", Local: tokens[1]}
		}
		return xml.Name{Space: fullUrl, Local: tokens[1]}
	}
	return xml.Name{Space: "", Local: str}
}

// attributesAsMap returns the attributes, which is optionally at index 1
func (n *jNode) attributesAsMap() map[string]string {
	tmp := make(map[string]string)
	if len(*n) > 1 {
		if m, ok := (*n)[1].(map[string]interface{}); ok {
			for k, v := range m {
				tmp[k] = stringOf(v)
			}
		}

	}
	return tmp
}

// attributes returns all xml attributes of the node. Optionally it is at index 1 but if so, always as a map.
// Because if the weired golang xml encoder, the xmlns attributes are discarded
func (n *jNode) attributes(nsList namespaces) []xml.Attr {
	tmp := n.attributesAsMap()
	res := make([]xml.Attr, 0)
	for k, v := range tmp {
		if strings.HasPrefix(k, "xmlns:") {
			continue
		}
		res = append(res, xml.Attr{Name: parseXMLName(nsList, k), Value: v})
	}
	return res
}

// children returns all child nodes of this node, which are string|*jNode. Optionally it either starts at offset 1 or 2
func (n *jNode) children() []interface{} {
	res := make([]interface{}, 0)
	for i := 1; i < len(*n); i++ {
		if _, isAttr := ((*n)[i]).(map[string]interface{}); isAttr {
			continue
		}
		child := cast((*n)[i])
		if child != nil {
			res = append(res, child)
		} else {
			str := stringOf((*n)[i])
			res = append(res, str)
		}
	}
	return res
}

func cast(any interface{}) *jNode {
	if node, ok := any.(*jNode); ok {
		return node
	}
	return nil
}

//==

// A nodeStack helps us to write the code without recursion
type nodeStack []*jNode

func (s *nodeStack) push(n *jNode) {
	*s = append(*s, n)
}

func (s *nodeStack) pop() *jNode {
	tmp := *s
	n := tmp[len(tmp)-1]
	*s = tmp[:len(tmp)-1]
	return n
}

func (s *nodeStack) top() *jNode {
	tmp := *s
	return tmp[len(tmp)-1]
}

//==

// nicefy takes a (primitive) string and tries to convert it into a nice type for json, like bool, number etc.
func nicefy(str string) interface{} {
	i, err := strconv.ParseInt(str, 10, 64)
	if err == nil {
		return i
	}
	f, err := strconv.ParseFloat(str, 64)
	if err == nil {
		return f
	}

	b, err := strconv.ParseBool(str)
	if err == nil {
		return b
	}
	return str
}

func stringOf(any interface{}) string {
	if s, ok := any.(string); ok {
		return s
	}
	return fmt.Sprintf("%v", any)
}
