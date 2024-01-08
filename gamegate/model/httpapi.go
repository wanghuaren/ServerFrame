package model

import (
	"baseutils/baseuts"
	"gameutils/common"
	"net/http"
	"time"

	"github.com/gin-contrib/timeout"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

func InitHttp(regisiter map[string]func(c *gin.Context), websocketCallback func(*websocket.Conn, []byte) []byte) {
	if IsDebug {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	websocketSvr = websocketCallback

	regisiter["socket"] = initWebSocket

	initHttpServer(regisiter)
	// initHttpWeb()
}

func initHttpWeb() {
	httpGin := getNewGin(true)

	httpGin.StaticFS("/", gin.Dir("./webroot", false))

	httpListenerStr := Conf.String("host")
	Log("http web", httpListenerStr)

	runGin(httpGin, httpListenerStr, true)
}

func initHttpServer(regisiter map[string]func(c *gin.Context)) {
	httpGin := getNewGin()

	for k, v := range regisiter {
		if k == "socket" {
			httpGin.GET(k, v)
		} else {
			httpGin.POST(k, v)
		}
	}

	httpListenerStr := Conf.String("host") + ":" + Conf.String("gate_port")
	Log("http server", httpListenerStr)

	runGin(httpGin, httpListenerStr)
}

func getNewGin(isWeb ...bool) *gin.Engine {
	httpGin := gin.New()
	// httpGin.Use(gin.Recovery())
	if len(isWeb) < 1 {
		httpGin.Use(timeoutMiddleware())
	}
	httpGin.Use(accessControlAllowOriginMiddleware())
	httpGin.SetTrustedProxies([]string{"127.0.0.1"})

	return httpGin
}

func runGin(httpGin *gin.Engine, httpListenerStr string, isWeb ...bool) {
	go func() {
		defer baseuts.ChkRecover()

		// if strings.Contains(LocalIP, "192.168") {
		if len(isWeb) > 0 {
			httpListenerStr += ":80"
		}
		httpGin.Run(httpListenerStr)
		// } else {
		// 	if len(isWeb) > 0 {
		// 		httpListenerStr += ":443"
		// 	}
		// 	if ServerInChina {
		// 		httpGin.RunTLS(httpListenerStr, "./gw.xhhuyu.com/gw.xhhuyu.com_bundle.pem", "./gw.xhhuyu.com/gw.xhhuyu.com.key")
		// 	} else {
		// 		httpGin.RunTLS(httpListenerStr, "./api.xhhuyu.com/api.xhhuyu.com_bundle.pem", "./api.xhhuyu.com/api.xhhuyu.com.key")
		// 	}
		// }
	}()
}

func timeoutMiddleware() gin.HandlerFunc {
	return timeout.New(
		timeout.WithTimeout(time.Duration(common.HTTP_TIMEOUT_MILLISECOND)*time.Millisecond),
		timeout.WithHandler(func(c *gin.Context) {
			c.Next()
		}),
		timeout.WithResponse(timeoutResponse),
	)
}

func timeoutResponse(c *gin.Context) {
	c.JSON(http.StatusGatewayTimeout, gin.H{
		"code": http.StatusGatewayTimeout,
		"msg":  "请重试",
	})
}

func accessControlAllowOriginMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		origin := c.Request.Header.Get("Origin")
		if origin != "" {
			c.Header("Access-Control-Allow-Origin", "*") // 可将将 * 替换为指定的域名
			c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")
			c.Header("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Authorization")
			c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Cache-Control, Content-Language, Content-Type")
			c.Header("Access-Control-Allow-Credentials", "true")
		}
		if method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
		}
		c.Next()
	}
}
