package main

import (
	"fmt"
	"micro-message-system/userserver/models"
)

func main() {
	for _, item := range models.UserSearch("vector") {
		fmt.Printf("%d, %s, %s\n", item.Id, item.Token, item.Username)
	}
	for _, item := range models.UserSearch("8e488ab4-7f1f-46d6-bd27-ece5f0673be8") {
		fmt.Printf("%d, %s, %s\n", item.Id, item.Token, item.Username)
	}
	for _, item := range models.UserSearch("xiaominig") {
		fmt.Printf("%d, %s, %s\n", item.Id, item.Token, item.Username)
	}
}
