// Package main imlements a client for movieinfo service
package main

import (
	"context"
	"fmt"
	"labs/lab-5/movieapi"
	"log"
	"os"
	"strings"
	"time"

	"google.golang.org/grpc"
)

const (
	address      = "localhost:50051"
	defaultTitle = "Pulp fiction"
)

var (
	title    string
	year     int32
	director string
	cast     []string
)

func set() {
	// Set up a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := movieapi.NewMovieInfoClient(conn)

	// Timeout if server doesn't respond
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	status, err := c.SetMovieInfo(ctx, &movieapi.MovieData{Title: title, Year: year, Director: director, Cast: cast})
	if err != nil {
		log.Fatalf("could not get movie info: %v", err)
	}
	log.Printf("Status Update %s", status.GetCode())
}

func get() {
	// Set up a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := movieapi.NewMovieInfoClient(conn)

	// Contact the server and print out its response.
	title := defaultTitle
	if len(os.Args) > 1 {
		title = os.Args[1]
	}
	// Timeout if server doesn't respond
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.GetMovieInfo(ctx, &movieapi.MovieRequest{Title: title})
	if err != nil {
		log.Fatalf("could not get movie info: %v", err)
	}
	log.Printf("Movie Info for %s %d %s %v", title, r.GetYear(), r.GetDirector(), r.GetCast())
}

func main() {
	var service string
	var castNames string

	fmt.Println("Enter type of service:\n 1. Get movie details \n 2. Set movie details")

	fmt.Scanln(&service)

	if service == "1" {
		get()
	} else {
		fmt.Println("Enter Movie Details:")
		fmt.Scanln(&title)
		fmt.Scanln(&year)
		fmt.Scanln(&director)
		fmt.Scanln(&castNames)

		cast = strings.Split(castNames, ",")

		set()
	}
}
