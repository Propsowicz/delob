package model

type TokenizedExpression struct {
	ProcessMethod ProcessMethod
	Arguments     []string
}

type ProcessMethod int8

const (
	AddPlayer ProcessMethod = iota
)
