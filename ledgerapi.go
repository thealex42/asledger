package main

import (
	"errors"
	"os"
	"runtime"
	"strconv"
	"sync"

	log "github.com/Sirupsen/logrus"
	"github.com/aerospike/aerospike-client-go"
	"github.com/gin-gonic/gin"
)

type ServerState struct {
	sync.Mutex
	ServerEnabled bool
}

var (
	DBHost        = "127.0.0.1"      // default database address
	DBPort        = 3000             // default database port
	Listen        = "127.0.0.1:8787" // default server
	DBNs          = "asledger"       // namespace
	DBTblAccounts = "acc"            // table with data
	DBTblStat     = "stat"           // table with stats counters
	RootBalance   = "root"           // name of system account
	MonetaryShift = 10000.00         // Money precision

	Clnt     *aerospike.Client
	CurState = &ServerState{
		ServerEnabled: true,
	}
)

func handleMoveFunds(c *gin.Context) {

	from := c.PostForm("from")
	to := c.PostForm("to")
	amount, err := strconv.ParseFloat(c.PostForm("amount"), 64)

	if err != nil {
		c.JSON(500, gin.H{
			"err": err.Error(),
		})
		return
	}

	modelFrom, _ := NewBank(from, Clnt)
	modelTo, _ := NewBank(to, Clnt)

	modelFrom.addFunds(0-amount, Clnt)
	modelTo.addFunds(amount, Clnt)

	BankSaveStats(amount, Clnt)

	// add root balance logic
}

func handleBalance(c *gin.Context) {

	bankModel, err := NewBank(c.Param("id"), Clnt)
	if err != nil {
		c.JSON(404, gin.H{
			"error": err.Error(),
		})
	}

	c.JSON(200, gin.H{
		"balance": bankModel.Balance,
	})
}

func handleChecksum(c *gin.Context) {
	CurState.Lock()
	CurState.ServerEnabled = false
	defer func() {
		CurState.ServerEnabled = true
		CurState.Unlock()
	}()

	scanPolicy := aerospike.NewScanPolicy()

	recs, err := Clnt.ScanAll(scanPolicy, DBNs, DBTblAccounts)
	if err != nil {
		c.JSON(500, gin.H{
			"error": err.Error(),
		})
	}

	var totalBalance int

	for res := range recs.Results() {
		if res.Err != nil {
			// handle error
		} else {
			if val, ok := res.Record.Bins["balance"]; ok {
				totalBalance += val.(int)
			}
		}
	}

	c.JSON(200, gin.H{
		"total": totalBalance,
	})
}

func handleStats(c *gin.Context) {

	funds, counter, err := BankGetStats(Clnt)
	if err != nil {
		c.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"counter": counter,
		"funds":   funds,
	})
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

	r.Use(func(c *gin.Context) {
		if !CurState.ServerEnabled {
			c.AbortWithError(500, errors.New("Checksum in progress"))
			return
		}
		c.Next()
	})

	r.POST("/funds/move", handleMoveFunds)
	r.GET("/balance/:id", handleBalance)
	r.GET("/checksum", handleChecksum)
	r.GET("/stats", handleStats)

	return r
}

func main() {

	r := buildRouter()
	r.Run(Listen)

}
