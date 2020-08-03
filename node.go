package main

type Node struct {
	Uid      string   `json:"uid,omitempty"`
	Name     string   `json:"name,omitempty"`
	Type     string   `json:"type,omitempty"`
	Children []string `json:"child,omitempty"`
}

func New(name string, children []string) *Node {
	return &Node{
		Uid:      "_:" + name,
		Name:     name,
		Children: children,
	}
}
