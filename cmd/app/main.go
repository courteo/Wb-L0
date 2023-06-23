package main

import (
	"database/sql"
	"encoding/json"
	// "io/ioutil"
	"log"
	// "os"
	// "withNats/pkg/model"
	"withNats/pkg/model/cache"
	m "withNats/pkg/model/db"

	"withNats/pkg/publisher"
	"withNats/pkg/structs"
	"withNats/pkg/subscriber"

	_ "github.com/jackc/pgx/v4/stdlib"
	"go.uber.org/zap"

	"github.com/gin-gonic/gin"
)

// docker run -p 4223:4223 -p 8223:8223 nats-streaming -p 4223 -m 8223



func main() {
	zapLogger, _ := zap.NewProduction()

	defer zapLogger.Sync()
	logger := zapLogger.Sugar()
	connStr := "postgres://postgres:qwerty@localhost:5432?sslmode=disable"
	db, err := sql.Open("pgx", connStr)
	if err != nil {
		logger.Panicf("sql open %s",err.Error())
	}

	db.SetMaxOpenConns(10)

	err = db.Ping()
	if err != nil {
		logger.Panicf("db ping  %s",err.Error())
	}
	Cache := cache.NatsDataCacheRepository{}
	ModelDb := m.NewMemoryRepo(db)

	subSettings := structs.SubcribeSettings{
		Channel: "wb",
		Claster: "wb-claster",
		Client: "client",
	}

	Subscriber := subscriber.Subscriber{
		SubscribeSt: subSettings,
		Db: ModelDb,
		Cache: &Cache,
	}
	err = Subscriber.DbToCache()
	if err != nil {
		log.Fatal("Subscriber DbToCache error ", err.Error())
	}
	go publisher.Publisher(subSettings)

	go Subscriber.Subscribe()

	// file, err := os.Open("data/model.json")
	// if err != nil {
	// 	log.Println(err)
	// }
	// byteValue, err := ioutil.ReadAll(file)
	// if err != nil {
	// 	log.Println("read all ", err.Error())
	// }
	// var data model.NatsData
	// err = json.Unmarshal(byteValue, &data)

	router := gin.Default()
	// Subscriber.Cache.Add(&data)

	router.LoadHTMLFiles("static/index.html", "static/bye_page.html")
	router.Static("static/css", "./static/css")

	var o structs.OrderJSON
	router.GET("/", func(c *gin.Context) {
		index(c)
		id := getId(c)
		log.Println("get ", id)
		data, err := Subscriber.Cache.FindNatsData(id)
		if err != nil {
			log.Println(err.Error())
			c.PureJSON(404, gin.H{
				"error": "Not found with id: " + id,
			})
			return
		} else {
			log.Printf("OK: found order with id: %v\n", id)
		}

		d, err := json.Marshal(data)
		o.DataJSON = string(d)
		o.Order_uid = id
		c.PureJSON(200, gin.H{
			"id":   o.Order_uid,
			"data": o.DataJSON,
		})
	})

	router.GET("/bye_page", func(c *gin.Context) {
		log.Println("bye_page ")
		c.PureJSON(200, gin.H{
			"id":   o.Order_uid,
			"data": o.DataJSON,
		})
	})

	router.GET("/static/css", page)

	router.Run("127.0.0.1:8090")
}

func getId(c *gin.Context) (id string) {
	id, ok := c.GetQuery("data")
	if !ok {
		log.Println("Can't get data from form")
		return ""
	}
	return id
}

func index(c *gin.Context) {
	c.HTML(200, "index.html", nil)
}

func page(c *gin.Context) {
	c.HTML(200, "page.css", nil)
}
