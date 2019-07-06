package main

import (
	"fmt"
	"github.com/Unknwon/goconfig"
	"./server"
)

func main() {
	cfg, err := goconfig.LoadConfigFile("server/conf.ini")
	if err != nil{
		fmt.Println("load config file error")
		return
	}

	ip, _ := cfg.GetValue("server", "ip")
	port, _ := cfg.Int("server", "port")

	s := server.CreateServer(
		ip,
		port,
		cfg,
	)

	fmt.Printf("server listen on %s:%d\n", ip, port)
	s.Start()
}