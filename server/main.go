package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type GameStatus string

type request struct {
	Direction string `json:"direction"`
}

type Game struct {
	BallX         int        `json:"ballX"`
	BallY         int        `json:"ballY"`
	BallRadius    int        `json:"ballRadius"`
	CanvasWidth   int        `json:"canvasWidth"`
	CanvasHeight  int        `json:"canvasHeight"`
	PaddleTopX    int        `json:"paddleTopX"`
	PaddleBottomX int        `json:"paddleBottomX"`
	PaddleTopY    int        `json:"paddleTopY"`
	PaddleBottomY int        `json:"paddleBottomY"`
	PaddleWidth   int        `json:"paddleWidth"`
	PaddleHeight  int        `json:"paddleHeight"`
	Status        GameStatus `json:"status"`
}

const (
	WaitingCompetitor GameStatus = "WAITING_COMPETITOR"
	InGame            GameStatus = "IN_GAME"
	GameOver          GameStatus = "GAME_OVER"
)

var (
	activeClients = 0
	ballDX        = 1
	ballDY        = -1
	game          Game
)

func sseSendInformation(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Content-Type", "text/event-stream")

	activeClients++

	if activeClients == 2 {
		game.Status = InGame
	} else if activeClients < 2 {
		game.Status = WaitingCompetitor
	}

	go func() {
		<-r.Context().Done()
		activeClients--
	}()

	for {
		if game.BallY > game.CanvasWidth-game.BallRadius || game.BallX < game.BallRadius {
			ballDX = -ballDX
		}

		if shouldReverBallByY() {
			ballDY = -ballDY
		}

		if game.BallY+game.BallRadius > game.CanvasHeight-game.PaddleHeight || game.BallY-game.BallRadius < game.PaddleHeight {
			game.Status = GameOver
		}

		game.BallX += ballDX
		game.BallY += ballDY

		jsonPayload, _ := json.Marshal(game)
		fmt.Fprintf(w, "data: %s\n\n", jsonPayload)
		w.(http.Flusher).Flush()
		time.Sleep(10 * time.Millisecond)
	}
}

func shouldReverBallByY() bool {
	isBallTouchedTopPaddle := game.BallX >= game.PaddleTopX && game.BallX < game.PaddleTopX+game.PaddleWidth && game.BallY-game.BallRadius <= game.PaddleHeight

	isBallTouchedBottomPaddle := game.BallX >= game.PaddleBottomX && game.BallX < game.PaddleBottomX+game.PaddleWidth && game.BallY+game.BallRadius >= game.CanvasHeight-game.PaddleHeight

	return isBallTouchedTopPaddle || isBallTouchedBottomPaddle
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

	if req.Direction == "RIGHT" && (game.PaddleTopX < game.CanvasWidth-game.PaddleWidth) {
		game.PaddleTopX += 16
	} else if req.Direction == "LEFT" && game.PaddleTopX > 0 {
		game.PaddleTopX -= 16
	}

	w.WriteHeader(http.StatusOK)
}

func main() {
	canvasWidth := 320
	canvasHeight := 160
	paddleWidth := 80
	paddleHeight := 10

	game = Game{
		BallX:         canvasWidth / 2,
		BallY:         canvasHeight - 30,
		BallRadius:    10,
		CanvasWidth:   canvasWidth,
		CanvasHeight:  canvasHeight,
		PaddleTopX:    (canvasWidth - paddleWidth) / 2,
		PaddleBottomX: (canvasWidth - paddleWidth) / 2,
		PaddleTopY:    0,
		PaddleBottomY: canvasHeight - paddleHeight,
		PaddleWidth:   75,
		PaddleHeight:  10,
		Status:        WaitingCompetitor,
	}

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
