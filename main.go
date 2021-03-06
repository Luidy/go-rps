//go:generate swagger generate spec
package main

import (
	"fmt"
	"gee/api/controller"
	"gee/config"
	"gee/repository"
	"log"
	"net/http"
	"runtime"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

const (
	banner = "\n" +
		"	********  ******** ********   \n" +
		"	**//////**/**///// /**/////   \n" +
		"   **      // /**      /**       \n" +
		"  /**         /******* /*******  \n" +
		"  /**    *****/**////  /**////   \n" +
		"  //**  ////**/**      /**       \n" +
		"   //******** /********/******** \n" +
		"	////////  //////// ////////   \n\n"
)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	fmt.Printf("server CPU core: %d\n", runtime.NumCPU())
}

func main() {
	Gee := config.Gee
	e := echoInit(Gee)
	db := repository.InitDB(Gee)
	if err := controller.InitHandler(Gee, e, db); err != nil {
		log.Println("InitHandler err", err.Error())
	}

	startServer(Gee, e)
}

func echoInit(gee *config.ViperConfig) *echo.Echo {
	e := echo.New()
	e.Pre(middleware.RemoveTrailingSlash())
	e.Use(middleware.Recover())

	healthCheck := func(c echo.Context) error {
		return c.String(http.StatusOK, "Gee Alive.")
	}
	e.GET("/healthCheck", healthCheck)

	e.HideBanner = true
	return e
}

func startServer(gee *config.ViperConfig, e *echo.Echo) {
	address := fmt.Sprintf("0.0.0.0:%d", gee.GetInt("port"))
	fmt.Println(banner, address)
	if err := e.Start(address); err != nil {
		log.Print("End echo server", "err", err)
	}
}
