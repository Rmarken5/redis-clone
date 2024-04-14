package main

import (
	"fmt"
	"github.com/tidwall/resp"
	"io"
	"log"
	"strings"
)

type Command interface {
}

type SetCommand struct {
	key, val []byte
}
type GetCommand struct {
	key []byte
}

func parseCommand(msg string) (Command, error) {

	rd := resp.NewReader(strings.NewReader(msg))
	for {
		v, _, err := rd.ReadValue()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}

		if v.Type() == resp.Array {
			for _, value := range v.Array() {
				switch value.String() {
				case "SET":
					if len(v.Array()) != 3 {
						return nil, fmt.Errorf("invalid number of variables for SET command")
					}
					cmd := SetCommand{
						key: v.Array()[1].Bytes(),
						val: v.Array()[2].Bytes(),
					}
					return cmd, nil
				case "GET":
					if len(v.Array()) != 2 {
						return nil, fmt.Errorf("invalid number of variables for GET command")
					}
					cmd := GetCommand{
						key: v.Array()[1].Bytes(),
					}
					return cmd, nil
				}
			}
		}
	}
	return nil, fmt.Errorf("unknown command")
}
