package main

type Object interface {
	getMethods() map[string]*RMethod
	getString() string
}
