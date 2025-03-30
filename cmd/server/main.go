package main

import "YuriyMishin/metrics/server"

func main() {
	srv := server.NewServer()
	if err := srv.Start("localhost:8080"); err != nil {
		panic(err)
	}
}
