package main

import (
	"flag"
	"log"
	"net/http"
	"time"

	"github.com/fatih/color"
	"github.com/gin-gonic/gin"
	"go.uber.org/ratelimit"
)

var (
	limit ratelimit.Limiter
	rps   = flag.Int("rps", 200, "request per second")
)

const html = `<!DOCTYPE html>
<html lang="en">

<head>
	<meta charset="UTF-8">
	<meta http-equiv="X-UA-Compatible" content="IE=edge">
	<meta content="initial-scale=1.0, minimum-scale=1.0, maximum-scale=1.0, width=device-width" name="viewport">
	<meta name="keywords" content="视频会议软件, 电话会议软件, 在线办公软件, 办公软件, 文档软件, 打卡软件app, 在线文档, 开会软件, feishu, lark, feishu, 飞书">
	<link rel="icon"
		href="data:image/svg+xml,<svg xmlns=%22http://www.w3.org/2000/svg%22 viewBox=%220 0 100 100%22><text y=%22.9em%22 font-size=%2290%22>😑</text></svg>">
	<title>Is Lark/Feishu linux-ready now?</title>
	<style>
		body,
		html {
			padding: 0;
			margin: 0;
			font-family: -apple-system, BlinkMacSystemFont, Segoe UI, Roboto, Oxygen, Ubuntu, Cantarell, Fira Sans, Droid Sans, Helvetica Neue, sans-serif;
		}

		h1,
		p {
			text-align: center;
		}

		h1 {
			line-height: 1.15;
			font-size: 4rem;
		}

		p {
			line-height: 1.5;
			font-size: 1.5rem;
		}
	</style>
	<script async src="https://www.googletagmanager.com/gtag/js?id=G-JQJ6SKB9XQ"></script>
	<script>
		window.dataLayer = window.dataLayer || [];
		function gtag() { dataLayer.push(arguments); }
		gtag('js', new Date());
		gtag('config', 'G-JQJ6SKB9XQ');
	</script>

</head>

<body>
	<h1>Is Lark/Feishu linux-ready now?</h1>
	<p id="result">Nope.</p>
	<script>
		(function () {
			var url = "https://www.feishu.cn/api/downloads";
			var xhr = new XMLHttpRequest();
			xhr.onreadystatechange = function () {
				if (xhr.readyState === XMLHttpRequest.DONE && xhr.status === 200) {
					var result = JSON.parse(xhr.responseText);
					if (result.versions.Linux) {
						document.querySelector('#result').innerHTML = 'Yes!!!'
					}
				}
			}
			xhr.open('GET', url);
			xhr.send();
		})();
	</script>
</body>

</html>`

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

	app.GET("/", func(c *gin.Context) {
		c.Writer.WriteHeader(http.StatusOK)
		c.Writer.Write([]byte(html))
	})

	app.Run()
}

func main() {
	flag.Parse()
	run(*rps)
}
