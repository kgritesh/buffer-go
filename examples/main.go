package main

import "fmt"
import "github.com/kgritesh/buffer-go/buffer"

func main() {
	client := buffer.GetOauth2Client("1/020ba75e8e3c93b92516771bc42915ed")
	bufferClient := buffer.NewClient(client)
	// order := []string{"556ab34f57d5d1443027d8f7", "555716780015853b5db43b82"}

	opts := &buffer.UpdateEditOptions{Text: "Updating an update using api"}

	updates, resp, err := bufferClient.UpdateService.EditUpdate("557c379cd7a0397149359467", opts)
	fmt.Println(updates, resp, err)
}
