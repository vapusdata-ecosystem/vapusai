package main

import (
	server "github.com/vapusdata-ecosystem/vapusai-studio/aistudio/server"
)

func main() {
	// This is the main function to start the platform server
	server := server.GrpcServer()
	server.Run()
}
