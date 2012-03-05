package main

type Source struct {
	root string
}

func NewSource(root string) *Source {
	return &Source{root: root}
}
