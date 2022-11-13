package main

import (
	config3 "auth-service/config"
	config2 "auth-service/config/config"
)

func main() {
	config := config2.NewConfig()
	server := config3.NewServer(config)
	server.Start()
}
