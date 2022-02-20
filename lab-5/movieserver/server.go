// Package main implements a server for movieinfo service.
package main

import (
	"context"
	"errors"
	"labs/lab-5/movieapi"
	"log"
	"net"
	"strconv"
	"strings"

	"google.golang.org/grpc"
)

const (
	port = ":50051"
)

// server is used to implement movieapi.MovieInfoServer
type server struct {
	movieapi.UnimplementedMovieInfoServer
}

// Map representing a database
var moviedb = map[string][]string{"Pulp fiction": []string{"1994", "Quentin Tarantino", "John Travolta,Samuel Jackson,Uma Thurman,Bruce Willis"}}

// GetMovieInfo implements movieapi.MovieInfoServer
func (s *server) GetMovieInfo(ctx context.Context, in *movieapi.MovieRequest) (*movieapi.MovieReply, error) {
	title := in.GetTitle()
	log.Printf("Received: %v", title)
	reply := &movieapi.MovieReply{}
	if val, ok := moviedb[title]; !ok { // Title not present in database
		return reply, nil
	} else {
		if year, err := strconv.Atoi(val[0]); err != nil {
			reply.Year = -1
		} else {
			reply.Year = int32(year)
		}
		reply.Director = val[1]
		cast := strings.Split(val[2], ",")
		reply.Cast = append(reply.Cast, cast...)

	}

	return reply, nil

}

func (s *server) SetMovieInfo(ctx context.Context, in *movieapi.MovieData) (*movieapi.Status, error) {
	var castName string
	status := &movieapi.Status{}

	title := in.GetTitle()
	year := in.GetYear()
	director := in.GetDirector()
	cast := in.GetCast()

	if (title == "") || (director == "") || (year == 0) || (len(cast) == 0) {
		status.Code = "Failed to set Movie info!!!"
		return status, errors.New("ERROR: Invalid Input!!!")
	} else {
		moviedb[title] = append(moviedb[title], strconv.FormatInt(int64(year), 10))
		moviedb[title] = append(moviedb[title], director)

		for i, name := range cast {
			if i == 0 {
				castName = name + ","
			} else if i == len(cast)-1 {
				castName = castName + name
			} else {
				castName = castName + name + ","
			}
		}

		moviedb[title] = append(moviedb[title], castName)

		status.Code = "Movie info Successfully set!!!"

		return status, nil
	}
}

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	movieapi.RegisterMovieInfoServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
