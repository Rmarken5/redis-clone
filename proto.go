package main

type SetCommand struct {
	key, val []byte
}
type GetCommand struct {
	key []byte
}
