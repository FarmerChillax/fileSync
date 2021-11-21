package main

import (
	"context"
	"fileSync/service"
	"fmt"
	"log"
)

func init() {

}

func main() {
	host := "127.0.0.1"
	port := "5000"

	ctx, err := service.Run(context.Background(), host, port)
	if err != nil {
		log.Println(err)
	}
	<-ctx.Done()
	fmt.Println("Shutting down scanner service.")
}
