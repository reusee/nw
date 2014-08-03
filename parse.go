package hns

import (
	"bytes"
	"fmt"
	"io"
	"strings"

	"code.google.com/p/go.net/html"
)

func Parse(r io.Reader) (*Node, error) {
	root := &Node{
		Tag:    "ROOT",
		rawBuf: new(bytes.Buffer),
	}
	tokenizer := html.NewTokenizer(r)
	currentNode := root
	writeRaw := func() {
		raw := tokenizer.Raw()
		currentNode.rawBuf.Write(raw)
		node := currentNode
		for node.Parent != nil {
			node = node.Parent
			node.rawBuf.Write(raw)
		}
	}
parse:
	for {
		what := tokenizer.Next()
		switch what {
		case html.ErrorToken:
			break parse
		case html.TextToken:
			text := strings.TrimSpace(string(tokenizer.Text()))
			if len(text) > 0 {
				currentNode.Text += text
				currentNode.TextParts = append(currentNode.TextParts, text)
			}
			writeRaw()
		case html.StartTagToken:
			node := &Node{
				Parent: currentNode,
				rawBuf: new(bytes.Buffer),
				Attr:   make(map[string]string),
			}
			currentNode.Children = append(currentNode.Children, node)
			currentNode = node
			writeRaw()
			name, hasAttr := tokenizer.TagName()
			currentNode.Tag = string(name)
			if hasAttr {
				key, val, more := tokenizer.TagAttr()
				currentNode.Attr[string(key)] = string(val)
				for more {
					key, val, more = tokenizer.TagAttr()
					currentNode.Attr[string(key)] = string(val)
				}
			}
		case html.EndTagToken:
			name, _ := tokenizer.TagName()
			if string(name) != currentNode.Tag { // tag mismatched
				return nil, fmt.Errorf("end tag mismatched, expected %s, got %s", currentNode.Tag, name)
			}
			writeRaw()
			currentNode.Raw = string(currentNode.rawBuf.Bytes())
			currentNode = currentNode.Parent
		case html.SelfClosingTagToken:
			node := &Node{
				Parent: currentNode,
				Raw:    string(tokenizer.Raw()),
				Attr:   make(map[string]string),
			}
			name, hasAttr := tokenizer.TagName()
			node.Tag = string(name)
			if hasAttr {
				key, val, more := tokenizer.TagAttr()
				node.Attr[string(key)] = string(val)
				for more {
					key, val, more = tokenizer.TagAttr()
					node.Attr[string(key)] = string(val)
				}
			}
			currentNode.Children = append(currentNode.Children, node)
			writeRaw()
		case html.CommentToken:
			writeRaw()
		}
	}
	root.Raw = string(root.rawBuf.Bytes())

	return root, nil
}
