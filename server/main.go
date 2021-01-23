package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
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
	ballDX   = 1
	ballDY   = -1
	game     Game
	players  = make(map[string]string)
	closeSse = make(chan bool)
)

func sseSendInformation(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Content-Type", "text/event-stream")

	clientIdCookie := getClientIdCookie(r)
	if clientIdCookie != "" {
		if _, ok := players["PaddleTopX"]; !ok {
			players["PaddleTopX"] = clientIdCookie
		}

		if _, ok := players["PaddleBottomX"]; !ok {
			players["PaddleBottomX"] = clientIdCookie
		}
	}

	if players["PaddleTopX"] != "" && players["PaddleBottomX"] != "" {
		game.Status = InGame
	}

	go func() {
		<-r.Context().Done()

		game = initGame()

		delete(players, "PaddleTopX")
		delete(players, "PaddleBottomX")

		closeSse <- true
	}()

	flushNewData(w)
}

func getClientIdCookie(r *http.Request) string {
	cookies := strings.Split(r.Header["Cookie"][0], ";")
	clientIdCookie := ""

	for _, cookie := range cookies {
		cookie = strings.Trim(cookie, " ")

		if strings.HasPrefix(cookie, "client-id") {
			clientIdCookie = strings.Split(cookie, "=")[1]

			break
		}
	}

	return clientIdCookie
}

func flushNewData(w http.ResponseWriter) {
	for {
		select {
		case <-closeSse:
			return
		default:
			if game.Status == InGame {

				if game.BallX >= game.CanvasWidth-game.BallRadius || game.BallX <= game.BallRadius {
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
			}
			jsonPayload, _ := json.Marshal(game)
			fmt.Fprintf(w, "data: %s\n\n", jsonPayload)
			w.(http.Flusher).Flush()
			time.Sleep(100 * time.Millisecond)
		}
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

	if game.Status == InGame {
		clientIdCookie := getClientIdCookie(r)

		paddlePositionX := &game.PaddleTopX

		if players["PaddleBottomX"] == clientIdCookie {
			paddlePositionX = &game.PaddleBottomX
		}

		if req.Direction == "RIGHT" && (*paddlePositionX < game.CanvasWidth-game.PaddleWidth) {
			*paddlePositionX += 16
		} else if req.Direction == "LEFT" && *paddlePositionX > 0 {
			*paddlePositionX -= 16
		}
	}

	w.WriteHeader(http.StatusOK)
}

func initGame() Game {
	canvasWidth := 320
	canvasHeight := 160
	paddleWidth := 80
	paddleHeight := 10

	return Game{
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
}

func main() {
	game = initGame()

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

	http.Handle("/", setCookieMiddleware(gzipMiddleware(cacheControlMiddleware(http.FileServer(http.Dir("./static"))))))

	_ = http.ListenAndServeTLS(":8080", "security/cert.pem", "security/cert.key", nil)
}
