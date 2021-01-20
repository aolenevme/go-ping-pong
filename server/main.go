package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type GameStatus string

const (
	WaitingCompetitor GameStatus = "WAITING_COMPETITOR"
	InGame            GameStatus = "IN_GAME"
	YouWon            GameStatus = "YOU_WON"
	YouLost           GameStatus = "YOU_LOST"
)

type Game struct {
	FirstCompetitor  UiElement
	SecondCompetitor UiElement
	Ball             UiElement
	Status           GameStatus
}

type UiElement struct {
	X int
	Y int
}

var game = Game{UiElement{-1, -1}, UiElement{-1, -1}, UiElement{0, 0}, WaitingCompetitor}

func sseSendInformation(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Content-Type", "text/event-stream")

	if (game.FirstCompetitor == UiElement{-1, -1}) {
		game.FirstCompetitor = UiElement{0, 0}
	} else if (game.SecondCompetitor == UiElement{-1, -1}) {
		game.SecondCompetitor = UiElement{0, 0}
		game.Status = InGame
	}

	go func() {
		<-r.Context().Done()
		game.SecondCompetitor = UiElement{-1, -1}
		game.Status = WaitingCompetitor
	}()

	for {
		b, _ := json.Marshal(game)
		fmt.Fprintf(w, "data: %s\n\n", b)
		w.(http.Flusher).Flush()
		time.Sleep(10 * time.Millisecond)
	}
}

func sseGetInformation(w http.ResponseWriter, r *http.Request) {
	game.FirstCompetitor.X += 7
	fmt.Printf("%+v", game)
	w.WriteHeader(http.StatusOK)
}

func main() {
	sse := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			sseSendInformation(w, r)
		case "PUT":
			sseGetInformation(w, r)
		default:
			http.Error(w, "Not Found", http.StatusNotFound)
		}
	})

	http.Handle("/api/v1/sse", sse)

	http.Handle("/", gzipMiddleware(cacheControlMiddleware(http.FileServer(http.Dir("./static")))))

	_ = http.ListenAndServeTLS(":8080", "security/cert.pem", "security/cert.key", nil)
}
