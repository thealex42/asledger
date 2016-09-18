package main

import (
	"fmt"
	"os"
	"runtime"
	"strconv"

	log "github.com/Sirupsen/logrus"
	"github.com/aerospike/aerospike-client-go"
	"github.com/aerospike/aerospike-client-go/examples/shared"
	"github.com/gin-gonic/gin"
)

const DBNs = "ledger"
const DBTblAccounts = "acc"
const Listen = "127.0.0.1:8787"

var (
	Clnt   *aerospike.Client
	DBHost = "127.0.0.1"
	DBPort = 3000
)

func handleAddFunds(c *gin.Context) {

}
func handleMoveFunds(c *gin.Context) {

}
func handleBalance(c *gin.Context) {

	key, err := aerospike.NewKey(DBNs, DBTblAccounts, fmt.Sprintf("%s", c.Param("id")))
	panicOnError(err)

	fmt.Println(key)

	rec, err := Clnt.Get(nil, key)
	shared.PanicOnError(err)

	if rec == nil {
		c.JSON(404, gin.H{
			"err": "balance not found",
		})
		return
	}

	var balance float64

	if _, ok := rec.Bins["balance"]; ok && rec.Bins["balance"] != nil {
		balance = rec.Bins["balance"].(float64)
	}

	c.JSON(200, gin.H{
		"balance": balance,
	})
}

func handleControl(c *gin.Context) {

}
func panicOnError(err error) {
	if err != nil {
		panic(err)
	}
}
func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	logfd, err := os.OpenFile("log.txt", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	panicOnError(err)

	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(logfd)

	if os.Getenv("DBHOST") != "" {
		DBHost = os.Getenv("DBHOST")
	}

	if os.Getenv("DBPORT") != "" {
		DBPort, err = strconv.Atoi(os.Getenv("DBPOST"))
		panicOnError(err)
	}

	Clnt, err = aerospike.NewClient(DBHost, DBPort)
	panicOnError(err)
}

func buildRouter() *gin.Engine {
	r := gin.Default()

	r.POST("/funds/add", handleAddFunds)
	r.POST("/funds/move", handleMoveFunds)
	r.GET("/balance/:id", handleBalance)
	r.POST("/control/:machine/:op", handleControl)

	return r
}

func main() {

	r := buildRouter()
	r.Run(Listen)

}
