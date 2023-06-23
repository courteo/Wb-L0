package publisher

import (
	"io/ioutil"
	"log"
	"os"
	"time"
	"withNats/pkg/structs"

	"github.com/nats-io/stan.go"
)

func Publisher(settings structs.SubcribeSettings) {
	sc, err := stan.Connect(settings.Cluster, "client1", stan.NatsURL("nats://localhost:4223"))
	if err != nil {
		log.Fatal("main connect ", err.Error())
	}

	file, err := os.Open("data/model.json")
	if err != nil {
		log.Println(err)
	}
	byteValue, err := ioutil.ReadAll(file)
	if err != nil {
		log.Println("read all ", err.Error())
	}

	sc.Publish(settings.Channel, byteValue)
	time.Sleep(time.Second)
}
