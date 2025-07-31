package initialize

import (
	"fmt"
	"log"

	"Taskly.com/m/global"

	"github.com/elastic/go-elasticsearch/v8"
)

func InitElasticSearch() {
	m := global.Config.ElasticSearch
	host := m.Host
	port := m.Port
	AddressesString := fmt.Sprintf(`http://%s:%v`, host, port)
	global.Logger.Sugar().Info(AddressesString)
	cfg := elasticsearch.Config{
		Addresses: []string{
			AddressesString,
		},
	}
	es, err := elasticsearch.NewClient(cfg)
	if err != nil {
		log.Fatalf("Error creating Elasticsearch client: %s0", err)
	}
	global.Elasticsearch = es
	log.Println("Connected to Elasticsearch")
}
