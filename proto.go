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
	key, val string
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
						key: v.Array()[1].String(),
						val: v.Array()[2].String(),
					}
					return cmd, nil
				}
			}
		}
	}
	return nil, fmt.Errorf("unknown command")
}
