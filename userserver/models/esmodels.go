package models

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/olivere/elastic"
	"reflect"
	"time"
)

type MemberSE struct {
	Id         int64
	Token      string     `json:"token"`
	Username   string     `json:"username"`
	CreateTime *time.Time `json:"create_time"`
	UpdateTime *time.Time `json:"update_time"`
}

func (o *Members) ToSearchEntity() *MemberSE {
	// 获取ES搜索对应的结构体
	return &MemberSE{
		Id:         o.Id,
		Token:      o.Token,
		Username:   o.Username,
		CreateTime: o.CreateTime,
		UpdateTime: o.UpdateTime,
	}
}

func UserFromSearchEntity(value *MemberSE) *Members {
	data := &Members{
		Id:         value.Id,
		Token:      value.Token,
		Username:   value.Username,
		CreateTime: value.CreateTime,
		UpdateTime: value.UpdateTime,
	}
	return data
}

// 把数据库的对象转换为ES的Entity并插入到ES
func SearchUserBuildOne(ctx context.Context, item *Members, client *elastic.Client) {
	fmt.Printf("SearchUserBuildOne: %d", item.Id)
	if item != nil {
		fmt.Printf("SearchUserBuildOne: %d", item.Id)
		se_index := "user_info"
		se_type := "doc"
		se_id := fmt.Sprintf("%d", item.Id)
		e := item.ToSearchEntity()

		// 创建index 并指定index名、id把userId作为esId、内容
		put1, err := client.Index().
			Index(se_index).
			Type(se_type).
			Id(se_id).
			BodyJson(e).
			Do(ctx)
		if err != nil {
			// Handle error
			fmt.Println(err.Error())
		}
		fmt.Printf("Indexed %s to index %s, type %s\n", put1.Id, put1.Index, put1.Type)
	}
}

// 删除一条ES中的数据
func SearchUserRemoveOne(id int64) {
	fmt.Printf("SearchUserRemoveOne: %d", id)
	client, err := elastic.NewSimpleClient(elastic.SetURL("http://192.168.33.17:9200/"))
	if err != nil {
		// Handle error
		fmt.Println(err.Error())
		return
	}
	defer client.Stop()

	ctx := context.Background()
	se_index := "user_info"
	se_type := "doc"
	se_id := fmt.Sprintf("%d", id)
	res, err := client.Delete().
		Index(se_index).
		Type(se_type).
		Id(se_id).
		Do(ctx)
	if err != nil {
		// Handle error
		fmt.Println(err.Error())
	} else {
		fmt.Printf("delete result: %d, %s", res.Status, res.Result)
	}
}

func DumpQuery(src interface{}, err error) {
	if err != nil {
		panic(err)
	}
	data, err := json.MarshalIndent(src, "", "  ")
	if err != nil {
		panic(err)
	}
	fmt.Println(string(data))

}

func UserSearch(query string) []*Members {
	fmt.Printf("UserSearch: %s\n", query)
	var err error
	client, err := elastic.NewSimpleClient(elastic.SetURL("http://192.168.33.16:9200/"))
	if err != nil {
		// Handle error
		fmt.Println(err.Error())
		return nil
	}
	defer client.Stop()

	se_index := "user_info"
	se_type := "doc"

	//条件查询
	formatS := `{
  "multi_match": {
      "query":    "%s",
      "fields":   [ "token", "username"]
  }
}`
	rawQ := elastic.NewRawStringQuery(
		fmt.Sprintf(formatS, query))
	DumpQuery(rawQ.Source())
	searchResult, err := client.Search(se_index).Type(se_type).
		Query(rawQ).
		From(0).Size(30).
		Do(context.Background())
	if err != nil {
		// Handle error
		fmt.Println(err.Error())
		return nil
	}
	fmt.Printf("Query took %d milliseconds. total: %d\n", searchResult.TookInMillis, searchResult.Hits.TotalHits)

	result := []*Members{}
	var ttyp MemberSE
	for _, item := range searchResult.Each(reflect.TypeOf(ttyp)) {
		if t, ok := item.(MemberSE); ok {
			result = append(result, UserFromSearchEntity(&t))
		}
	}
	return result
}
