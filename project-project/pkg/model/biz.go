package model

var (
	Normal   = 1
	Personal = 1
	AESKey   = "abcfedgehjzabkmlkjjkkoew"
)

const (
	NoDeleted = iota
	Deleted
)

const (
	NoArchive = iota
	Archive
)

const (
	Open = iota
	Private
	Custom
)

const (
	Default = "default"
	Simple  = "simple"
)
