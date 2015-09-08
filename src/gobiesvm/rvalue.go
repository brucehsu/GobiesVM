package main

import "math/big"

type RValue struct {
	str    string
	fixnum *big.Int
	float  float64
}
