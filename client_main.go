package main

import (
	"fmt"
	"github.com/Unknwon/goconfig"
	"./client"
)

func main() {
	cfg, err := goconfig.LoadConfigFile("client/conf.ini")
	if err != nil{
		fmt.Println("load config file error")
		return
	}

	client_num, _ := cfg.Int("client", "client_num")
	server_ip, _ := cfg.GetValue("client", "server_ip")
	server_port, _ := cfg.Int("client", "server_port")

	cc := client.CreateClientContariner(client_num)

	cc.Start(server_ip, server_port)

	fmt.Printf("client started!")
}