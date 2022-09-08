package main

import (
	"Userclient/endport"
	"Userclient/http"
	"flag"
)

var port = flag.Int("port", 8080, "port")

func main() {
	endport.InitConf()
	flag.Parse()
	http.RunServer(*port)
}
