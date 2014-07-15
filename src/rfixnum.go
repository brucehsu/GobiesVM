package main

import "strconv"

type RFixnum struct {
	RObject
	val RValue
}

func (obj *RFixnum) getMethods() map[string]*RMethod {
	return obj.methods
}

func (obj *RFixnum) getString() string {
	return strconv.FormatInt(obj.val.fixnum, 10)
}
