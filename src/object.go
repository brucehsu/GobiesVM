package main

type Object interface {
	getMethods() map[string]*RMethod
}
