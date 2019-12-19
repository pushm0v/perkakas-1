package elastic

import (
	"context"
	"encoding/json"
	"fmt"
	es "github.com/olivere/elastic/v7"
	"github.com/stretchr/testify/assert"
	"testing"
)

var client *MyClient

type MyClient struct {
	*Client
}

func NewMyClient(url string) *MyClient{
	c, err := NewClient(url)
	if err != nil {
		fmt.Println(err)
	}

	return &MyClient {
		Client: c,
	}
}

func init() {
	var err error

	url := "http://localhost:9200"
	client = NewMyClient(url)
	if err != nil {
		panic(err)
	}
}

type Person struct {
	ID string `json:"id"`
	Name string `json:"name"`
	Age int `json:"age"`
}

func elasticCreateIndex() *DynamicTemplate {
	mappings := NewMappings().
		AddDynamicTemplate("id_field", MatchConditions{
			Match:            "id",
			MatchMappingType: "string",
			Mapping: MatchMapping{
				Type: "text",
			},
		}).
		AddDynamicTemplate("name_field", MatchConditions{
			Match:            "name",
			MatchMappingType: "string",
			Mapping: MatchMapping{
				Type: "text",
				Fields: map[string]Field{
					"keyword": Field{
						Type:        "keyword",
						IgnoreAbove: 256,
					},
				},
			},
		}).
		AddDynamicTemplate("age", MatchConditions{
			Match:            "age",
			MatchMappingType: "double",
			Mapping: MatchMapping{
				Type: "double",
			},
		})

	return &DynamicTemplate{
		Settings: map[string]interface{}{
			"index.refresh_interval": "5s",
		},
		Mappings: mappings,
	}
}

func TestClient_Store(t *testing.T) {
	person := Person{
		ID: "1234",
		Name: "Alvin Rizki",
		Age: 17,
	}

	ctx := context.Background()
	mapping := elasticCreateIndex()
	_, err := client.Store(ctx, "campaign-coba-777", person, mapping)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	assert.Nil(t, err)
}

func TestClient_Search(t *testing.T) {
	ctx := context.Background()
	searchOption := SearchOption{
		Query: es.NewWildcardQuery("id", "*1234*"),
		Sort:  nil,
		Size:  10,
	}

	_, err := client.Search(ctx, "campaign-coba-444", searchOption)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	assert.Nil(t, err)
}

func TestClient_Delete(t *testing.T) {
	ctx := context.Background()
	_, err := client.Remove(ctx, "campaign-coba-444", "1234")
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	assert.Nil(t, err)
}

func TestClient_DeleteIndex(t *testing.T) {
	ctx := context.Background()
	_, err := client.RemoveIndex(ctx, "person-list")
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	assert.Nil(t, err)
}

func TestCreateElasticIndex(t *testing.T) {
	idx := elasticCreateIndex()
	output, _ := json.Marshal(idx)
	t.Log(string(output))
}

func TestBulkStore(t *testing.T) {
	ctx := context.Background()
	docs := []interface{}{
		Person {
			ID:   "abcde",
			Name: "Budi",
			Age: 17,
		},
		Person {
			ID:   "fghij",
			Name: "Badu",
			Age: 50,
		},
	}

	template := elasticCreateIndex()
	err := client.BulkStore(ctx, "campaign-coba-777", "default", docs, template)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	assert.Nil(t, err)
}

