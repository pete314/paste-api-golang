//Author: Peter Nagy <https://peternagy.ie>
//Since: 06, 2017
//Description: --
package main

import (
	"./modules/v0/auth"
	"./modules/v0/common"
	"./modules/v0/paste"
	"flag"
	"github.com/buaazp/fasthttprouter"
	"github.com/valyala/fasthttp"
	"log"
	"os"
	"strconv"
)

var (
	conf = flag.String("conf", "./../config/server-config.json", "The server configuration file")
)

func main() {
	//Check if for cli params and parse them
	if len(os.Args) > 1 {
		flag.Parse()
	}

	//Load configuration
	runtimeConfig := common.LoadConfig(*conf)
	log.Println("Server configuration loaded from ", *conf)

	//Initialize the router
	router := fasthttprouter.New()

	//Init server custom changes
	common.DisableServerFeatures(router)
	common.HTTPNotFoundHandler(router)
	common.HTTPNotAllowedHandler(router)

	//Load API modules
	pc := paste.NewController()
	pc.EnableHTTPMethods(router)

	ac := auth.NewController()
	ac.EnableHTTPMethods(router)

	//Start server
	serverAddress := runtimeConfig.HTTPAddress + ":" + strconv.Itoa(runtimeConfig.HTTPPort)
	log.Println("Server is starting on ", serverAddress)
	log.Fatal(fasthttp.ListenAndServe(serverAddress, router.Handler))
}
