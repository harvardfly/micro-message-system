package main

import (
	"context"
	"fmt"
	"micro-message-system/userserver/models"

	"github.com/olivere/elastic"
)

func UserGetAll() {
	list, err := models.QueryUserWithCon()
	if err == nil {
		ctx := context.Background()
		client, err := elastic.NewSimpleClient(elastic.SetURL("http://192.168.33.16:9200/"))
		if err != nil {
			// Handle error
			fmt.Println(err.Error())
			return
		}
		defer client.Stop()
		for _, item := range list {
			fmt.Printf("user: %d, %s\n", item.Id, item.Username)
			models.SearchUserBuildOne(ctx, &item, client)
		}
	} else {
		fmt.Println(err.Error())
	}
}

func main() {
	UserGetAll()
}
