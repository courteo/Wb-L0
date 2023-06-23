package subscriber

import (
	"encoding/json"
	"log"
	"withNats/pkg/model"
	r "withNats/pkg/model/db"
	"withNats/pkg/structs"

	"github.com/nats-io/stan.go"
)

type Subscriber struct {
	SubscribeSt structs.SubcribeSettings
	Cache       model.NatsDataRepo
	Db          *r.NatsDataMemoryRepository
}

func (s *Subscriber) Subscribe() {

	sc, err := stan.Connect(s.SubscribeSt.Cluster, s.SubscribeSt.Client, stan.NatsURL("nats://localhost:4223"))
	if err != nil {
		log.Fatal("sub connect ", err.Error())
	}
	defer sc.Close()
	res := make(chan []byte)
	// Simple Async Subscriber
	_, err = sc.Subscribe(s.SubscribeSt.Channel, func(m *stan.Msg) {
		log.Printf("Received a message: %s\n", string(m.Data))
		res <- m.Data
	}, stan.DeliverAllAvailable())
	if err != nil {
		log.Fatal("Subcribe ", err.Error())
	}
	for {
		var received model.NatsData
		err := json.Unmarshal(<-res, &received)
		if err != nil {
			log.Println("Json: wrong type ", err.Error())
			continue
		}
		log.Println("Add to cache ", received)
		err = s.Cache.Add(&received)
		if err != nil {
			log.Fatal("Cache Add ", err.Error())
		}
		log.Println("Add to db ", received)
		err = s.Db.Add(&received)
		if err != nil {
			log.Fatal("db Add ", err.Error())
		}
	}
}

func (s *Subscriber) DbToCache() error {
	data, err := s.Db.Data.Query("select * from orders;")
	
	if err != nil {
		log.Println("DbToCache Query error ", err.Error())
		return err
	}
	log.Println("DbToCache Query")
	for data.Next() {
		var order structs.OrderJSON
		err := data.Scan(&order.Order_uid, &order.DataJSON)
		if err != nil {
			log.Println("DbToCache Scan error ", err.Error())
			return err
		}
		var res model.NatsData
		err = json.Unmarshal([]byte(order.DataJSON), &res)
		if err != nil {
			log.Println("DbToCache Unmarshal error ", err.Error())
			return err
		}
		err = s.Cache.Add(&res)
		if err != nil {
			log.Println("DbToCache Add error ", err.Error())
			return err
		}
	}
	return nil
}
