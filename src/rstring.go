package main

type RString struct {
	RObject
	val RValue
}

func (obj *RString) getMethods() map[string]*RMethod {
	return obj.methods
}

func (obj *RString) getString() string {
	return obj.val.str
}
