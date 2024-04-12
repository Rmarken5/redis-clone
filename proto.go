package main

import (
	"bytes"
	"fmt"
	"github.com/tidwall/resp"
	"io"
	"log"
)

type Command interface {
}

type SetCommand struct {
	key, val string
}

func parseCommand(msg []byte) (Command, error) {

	rd := resp.NewReader(bytes.NewReader(msg))
	for {
		v, _, err := rd.ReadValue()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Read %s\n", v.Type())
		if v.Type() == resp.Array {
			for i, value := range v.Array() {
				switch value.String() {
				case "SET":
					if len(value.Array()) != 3 {
						return nil, fmt.Errorf("invalid number of variables for SET command")
					}
					return SetCommand{
						key: value.Array()[1].String(),
						val: value.Array()[2].String(),
					}, nil
				default:
				}
				fmt.Printf("  #%d %s, value: '%s'\n", i, value.Type(), value)
			}
		}
		return nil, fmt.Errorf("unknown command")
	}

	return nil, fmt.Errorf("unknown command")
}
