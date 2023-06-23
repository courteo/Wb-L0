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

	_ "github.com/jackc/pgx/v5/stdlib"
	"go.uber.org/zap"

	"github.com/gin-gonic/gin"
)

// docker run -p 4223:4223 -p 8223:8223 nats-streaming -p 4223 -m 8223

// func main3() {
// 	urlExample := "postgres://username:password@localhost:5432/database_name"
// 	urlExample = "postgres://postgres:qwerty@localhost:5433/postgres?sslmode=disable"
// 	conn, err := pgx.Connect(context.Background(), urlExample)
// 	if err != nil {
// 		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
// 		os.Exit(1)
// 	}
// 	defer conn.Close(context.Background())

// 	var name string
// 	var weight int64
// 	err = conn.QueryRow(context.Background(), `create table if not exists test (
// 		order_id varchar(10) not null unique,
// 		data json not null
// 	);`).Scan()

// 	// err = conn.QueryRow(context.Background(), "select * from test", 42).Scan(&name, &weight)
// 	if err != nil {
// 		fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
// 		os.Exit(1)
// 	}

// 	fmt.Println(name, weight)
// }

type Test struct {
	ID   int
	Name string
}

func main() {
	zapLogger, _ := zap.NewProduction()

	defer zapLogger.Sync()
	logger := zapLogger.Sugar()
	connStr := "postgres://postgres:qwerty@localhost:5433/postgres?sslmode=disable"
	db, err := sql.Open("pgx", connStr)
	if err != nil {
		logger.Panicf("sql open %s", err.Error())
	}

	db.SetMaxOpenConns(10)

	err = db.Ping()
	if err != nil {
		logger.Panicf("db ping  %s", err.Error())
	}

	Cache := cache.NatsDataCacheRepository{}
	ModelDb := m.NewMemoryRepo(db)

	subSettings := structs.SubcribeSettings{
		Channel: "wb",
		Cluster: "test-cluster",
		Client:  "client",
	}

	Subscriber := subscriber.Subscriber{
		SubscribeSt: subSettings,
		Db:          ModelDb,
		Cache:       &Cache,
	}
	err = Subscriber.DbToCache()
	if err != nil {
		log.Fatal("Subscriber DbToCache error ", err.Error())
	}
	go publisher.Publisher(subSettings)

	go Subscriber.Subscribe()


	router := gin.Default()

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
