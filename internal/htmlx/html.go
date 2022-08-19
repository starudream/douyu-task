package htmlx

import (
	"golang.org/x/net/html"
)

type (
	Node      = html.Node
	NodeType  = html.NodeType
	Attribute = html.Attribute
)

const (
	ErrorNode    = html.ErrorNode
	TextNode     = html.TextNode
	DocumentNode = html.DocumentNode
	ElementNode  = html.ElementNode
	CommentNode  = html.CommentNode
	DoctypeNode  = html.DoctypeNode
	RawNode      = html.RawNode
)

var (
	Parse = html.Parse
)
