package main

import "fmt"
import "github.com/kgritesh/buffer-go/buffer"

func main() {
	client := buffer.GetOauth2Client("1/020ba75e8e3c93b92516771bc42915ed")
	bufferClient := buffer.NewClient(client)
	user, resp, err := bufferClient.UserService.Get()
	fmt.Println(user, resp, err)
}
