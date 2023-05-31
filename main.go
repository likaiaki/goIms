package main

import "awesomeProject/server_tool"

func main() {
	server := server_tool.NewServer("127.0.0.1", 8888)
	server.Start()
}
