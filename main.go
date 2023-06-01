package main

import "awesomeProject/serverModule"

func main() {
	server := serverModule.NewServer("127.0.0.1", 8888)
	server.Start()
}
