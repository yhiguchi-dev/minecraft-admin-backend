package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"rcon-client/internal/handler/user/get"
	"rcon-client/internal/handler/user/post"
	"rcon-client/internal/rcon"
	"time"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	rconServerHost := os.Getenv("RCON_SERVER_HOST")
	rconServerPort := os.Getenv("RCON_SERVER_PORT")
	rconServerPassword := os.Getenv("RCON_SERVER_PASSWORD")

	conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%s", rconServerHost, rconServerPort), 3*time.Second)

	if err != nil {
		log.Fatal(err)
	}
	client, err := rcon.NewClient(conn, rconServerPassword)
	if err != nil {
		log.Fatal(err)
	}
	operator := rcon.NewOperator(client)
	get.NewGetUserListHandler(operator)
	post.NewPostUserListHandler(operator)

	http.HandleFunc("/users", get.GetUserListHandler)
	http.HandleFunc("/users", post.PostUserListHandler)

	//hostAddress := os.Getenv("HOST_ADDRESS")
	hostAddress := ":8081"
	server := &http.Server{
		Addr:    hostAddress,
		Handler: nil,
	}
	go func() {
		<-ctx.Done()
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		err := server.Shutdown(ctx)
		if err != nil {
			return
		}
	}()
	fmt.Println("call server")
	log.Fatal(server.ListenAndServe())
}
