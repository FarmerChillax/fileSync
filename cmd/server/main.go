package main

import (
	"context"
	"fileSync/core"
	"fileSync/service"
	"fmt"
	"log"
)

func init() {

}

func main() {
	host := core.Conf.Host
	port := core.Conf.Port

	ctx, err := service.Run(context.Background(), host, port)
	if err != nil {
		log.Println(err)
	}
	<-ctx.Done()
	fmt.Println("Shutting down scanner service.")
}
