package main

import (
	"flag"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/fatih/color"
	"github.com/gin-gonic/gin"
	"github.com/tidwall/gjson"
	"go.uber.org/ratelimit"
)

const url = "https://www.feishu.cn/api/downloads"

var (
	limit ratelimit.Limiter
	rps   = flag.Int("rps", 20, "request per second")
)

func leakBucket() gin.HandlerFunc {
	prev := time.Now()
	return func(ctx *gin.Context) {
		now := limit.Take()
		log.Print(color.CyanString("%v", now.Sub(prev)))
		prev = now
	}
}

func run(rps int) {
	limit = ratelimit.New(rps)
	app := gin.Default()
	app.Use(leakBucket())
	app.LoadHTMLFiles("tpl.html")

	app.GET("/", func(c *gin.Context) {
		resp, reqErr := http.Get(url)

		if reqErr != nil {
			log.Fatal(reqErr)
		}

		defer resp.Body.Close()

		body, readErr := ioutil.ReadAll(resp.Body)
		if readErr != nil {
			log.Fatal(readErr)
		}

		json := string(body)
		linux := gjson.Get(json, "versions.Linux")

		result := "Nope."

		if linux.Exists() {
			result = "Yes!"
		}

		c.HTML(http.StatusOK, "tpl.html", gin.H{
			"result": result,
		})
	})

	app.Run()
}

func main() {
	flag.Parse()
	run(*rps)
}
