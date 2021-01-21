package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type request struct {
	Direction string `json:"direction"`
}

type GameStatus string

const (
	WaitingCompetitor GameStatus = "WAITING_COMPETITOR"
	InGame            GameStatus = "IN_GAME"
	YouWon            GameStatus = "YOU_WON"
	YouLost           GameStatus = "YOU_LOST"

	canvasWidth = 320
	canvasHeight = 160
	paddleWidth = 75
	paddleHeight = 10
	paddleTopY = 0
	paddleBottomY = canvasHeight-paddleHeight
)

var (
	activeClients = 0
	ballX = canvasWidth/2
	ballY = canvasHeight - 30
	ballDX = 2
	ballDY = -2
	gameStatus = WaitingCompetitor
	paddleTopX = (canvasWidth - paddleWidth)/2
	paddleBottomX = (canvasWidth - paddleWidth)/2
)

func sseSendInformation(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Content-Type", "text/event-stream")

	activeClients++

	if activeClients == 2 {
		gameStatus = InGame
	} else if activeClients < 2 {
		gameStatus = WaitingCompetitor
	}

	go func() {
		<-r.Context().Done()
		activeClients--
	}()

	for {
		// b, _ := json.Marshal(game)
		fmt.Fprintf(w, "data: %d\n\n", paddleTopX)
		w.(http.Flusher).Flush()
		time.Sleep(10 * time.Millisecond)
	}
}

func sseGetInformation(w http.ResponseWriter, r *http.Request) {
	b, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	var req request

	err = json.Unmarshal(b, &req)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	if req.Direction == "RIGHT" && (paddleBottomX < canvasWidth - paddleWidth) {
		paddleTopX += 7
	} else if req.Direction == "LEFT" && paddleTopX > 0 {
		paddleTopX -= 7
	}

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
