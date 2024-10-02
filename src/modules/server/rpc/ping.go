package rpc

import "fmt"

func (*Server) Ping(input string, output *string) error {
	fmt.Println(input)
	*output = "收到了"
	return nil
}
