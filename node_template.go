package node_template

import (
	"bytes"
	"code.google.com/p/cascadia"
	"code.google.com/p/go.net/html"
	"container/list"
	"io"
	"io/ioutil"
)

/* Parse HTML into a NodeTemplate */
func Parse(r io.Reader) (*NodeTemplate, error) {
	node, err := html.Parse(r)
	if err != nil {
		return nil, err
	}
	nt := NodeTemplate{node}
	return &nt, nil
}

/* NodeTemplateFromFile gives you a shortcut to parse a template from a filename.

*/
func NodeTemplateFromFile(filename string) (*NodeTemplate, error) {
	f, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	r := bytes.NewReader(f)
	return Parse(r)
}

type NodeTemplate struct {
	*html.Node
}

/* Find the first node that matces the jquery-style pattern (see  code.google.com/p/cascadia docs for details) and returns a NodeTemplate to represent that. */
func (t *NodeTemplate) FindFirst(pat string) (*NodeTemplate, error) {
	sel, err := cascadia.Compile(pat)
	if err != nil {
		return nil, err
	}
	node := sel.MatchFirst(t.Node)
	if err != nil {
		return nil, err
	}
	if node == nil {
		return nil, nil
	}
	nt := NodeTemplate{node}
	return &nt, nil
}

/* Final all nodes which match the jquery-style pattern */
func (t *NodeTemplate) Find(pat string) (*NodeTemplateSet, error) {
	sel, err := cascadia.Compile(pat)
	if err != nil {
		return nil, err
	}
	nodes := sel.MatchAll(t.Node)
	if err != nil {
		return nil, err
	}
	nts := make(NodeTemplateSet, len(nodes))
	for i := range nodes {
		nts[i] = &NodeTemplate{nodes[i]}
	}
	return &nts, nil
}

/* Replace the content of the node a string. This is html safe, so &< etc are fine in the string.*/
func (t *NodeTemplate) ReplaceContentText(content string) {
	new_node := NodeTemplate{&html.Node{Type: html.TextNode, Data: content}}
	t.ReplaceContent(&new_node)
}

/* Replace the content of a node with new_node */
func (t *NodeTemplate) ReplaceContent(new_node *NodeTemplate) {
	subnode := t.FirstChild
	for subnode != nil {
		n := subnode
		subnode = n.NextSibling
		t.RemoveChild(n)
	}
	t.AppendChild(new_node.Node)
}

func (t *NodeTemplate) Copy() *NodeTemplate {
	new_n := html.Node{}

	new_n.Parent = nil
	new_n.NextSibling = nil
	new_n.PrevSibling = nil
	new_n.FirstChild = nil
	new_n.LastChild = nil
	new_n.Data = t.Data
	new_n.Type = t.Type
	new_n.DataAtom = t.DataAtom
	new_n.Attr = t.Attr
	new_n.Namespace = t.Namespace

	subnode := t.FirstChild
	for subnode != nil {
		nt_subnode := NodeTemplate{subnode}
		nt := nt_subnode.Copy()
		new_n.AppendChild(nt.Node)
		subnode = subnode.NextSibling
	}

	return &NodeTemplate{&new_n}
}

func (t *NodeTemplate) RepeatNode(l *list.List, f func(*NodeTemplate, *list.Element)) {

	if t.Parent == nil {
		// this wont work
		return
	}
	if t == nil {
		return
	}
	for e := l.Front(); e != nil; e = e.Next() {
		n := t.Copy()
		t.Parent.InsertBefore(n.Node, t.Node)
		f(n, e)
	}

	// finally kill the original!
	t.Parent.RemoveChild(t.Node)
}

/* Render the HTML to an io.Writer */
func (t *NodeTemplate) Render(w io.Writer) {
	html.Render(w, t.Node)
}

type NodeTemplateSet []*NodeTemplate

/* For each node execute this function */
func (nts NodeTemplateSet) For(f func(*NodeTemplate)) {
	for i := range nts {
		f(nts[i])
	}
}
func (nts NodeTemplateSet) Len() int {
	return len(nts)
}
func (nts NodeTemplateSet) Get(index int) *NodeTemplate {
	if index > nts.Len() {
		return nil
	}
	return nts[index]
}

/* Replace content of each node with a string */
func (nts NodeTemplateSet) ReplaceContentText(content string) {
	nts.For(func(t *NodeTemplate) { t.ReplaceContentText(content) })
}
func (nts NodeTemplateSet) ReplaceContent(content *NodeTemplate) {
	nts.For(func(t *NodeTemplate) { t.ReplaceContent(content) })
}
