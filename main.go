package main

import (
	"bytes"
	"errors"
	"log"
	"sync"

	"github.com/hashicorp/golang-lru"
	"github.com/valyala/fasthttp"
	"github.com/valyala/tcplisten"
)

var allUsers map[int]user
var allLocations map[int]location
var allVisits map[int]visit
var allUsersVisit map[int][]int
var allLocationsVisit map[int][]int
var allUsersMutex = &sync.Mutex{}
var allLocationsMutex = &sync.Mutex{}
var allVisitsMutex = &sync.Mutex{}

var currentTime int
var errorBadRequest = errors.New("error")
var routeUsers = []byte("/users/")
var routeLocations = []byte("/locations/")
var emptyJson = []byte("{}")
var cache, _ = lru.New(100 * 1024)

const oneYear = 31557600

func main() {
	loadFromFile()

	requestHandler := func(ctx *fasthttp.RequestCtx) {
		switch {
		case bytes.HasPrefix(ctx.Path(), routeUsers):
			usersHandler(ctx)
		case bytes.HasPrefix(ctx.Path(), routeLocations):
			locationsHandler(ctx)
		default:
			visitsHandler(ctx)
		}
	}

	listenerConfig := tcplisten.Config{
		ReusePort:   true,
		DeferAccept: true,
		FastOpen:    true,
	}

	ln, err := listenerConfig.NewListener("tcp4", ":80")

	if err != nil {
		log.Fatalf("error in reuseport listener: %s", err)
	}

	server := fasthttp.Server{
		Handler:      requestHandler,
		LogAllErrors: true,
	}

	if err = server.Serve(ln); err != nil {
		log.Fatalf("error in fasthttp Server: %s", err)

	}

}
