package server

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"runtime"
	"time"

	"apis_load_test/benchmark/server/apis"
	"apis_load_test/benchmark/server/ws"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
)

var all_router *gin.Engine

var SERVER_PORT string = "7000"

func openbrowser(url string) {
	var err error

	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}
	if err != nil {
		log.Fatal(err)
	}

}

func RunServer(disable_color bool,gin_mode string) {

	if disable_color {
		gin.DisableConsoleColor()
	} else {
		gin.ForceConsoleColor()
	}

	all_router = gin.Default()

	if gin_mode == "release" {
		all_router = gin.New()
		all_router.Use(gin.Recovery())
	}
	
	all_router.Use(cors.Default())
	all_router.Use(static.Serve("/", static.LocalFile("./benchmark/server/metrics_gui/build", true)))
	{
		api_router := all_router.Group("/api")
		apis.InitApis(api_router)
	}
	{
		ws_router := all_router.Group("/go_ws")
		ws.InitWS(ws_router)
	}

	if os.Getenv("SERVER_PORT") != "" {
		SERVER_PORT = os.Getenv("SERVER_PORT")
	}

	bind_to_host := fmt.Sprintf(":%s", SERVER_PORT) //formatted host string
	fmt.Printf("\nRunning server on http://localhost%s\n",bind_to_host)

	go func ()  {
		// if gin_mode	== "release" {
		// 	time.Sleep(time.Second*5)
		// 	openbrowser("http://"+bind_to_host)
		// }
		time.Sleep(time.Second*5)
		openbrowser("http://localhost"+bind_to_host)
	}()
	all_router.Run(bind_to_host)
}
