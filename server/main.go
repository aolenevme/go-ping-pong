package main

import (
	"fmt"
	"net/http"
	"time"
)

// 1. Interface Connection
// 	1.1. FirstClient Client
// 	1.2. SecondClient Client

// 2. Interface Client
// 	2.1. method send -- ballX, ballY, enemyX, enemyY
// 	2.2. method accept -- clientX, clientY

var gameStatus = "WAITING_COMPETITOR"

var game []Competitor

type Competitor struct {
	w http.ResponseWriter
	r *http.Request
	paddleX int
	paddleY int
}

func sseSendInformation(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Content-Type", "text/event-stream")

	if len(game) < 1 {
		game = append(game, Competitor{w, r, 0, 0})
	} else {
		gameStatus = "IN_GAME"
	}

	for {
		fmt.Fprintf(w, "data: %s\n\n", gameStatus)
		w.(http.Flusher).Flush()
		time.Sleep(time.Second)
	}
}

func sseGetInformation(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "PUT")
		// 1. Есть игра
		// 2. Грузится игра. Игра инициализирует handshake: /api/v1/sse. Сервер создает Competitor. В структуру Game добавляется competitor. Competitor сохраняет w и r в своей структуре. Возвращает статус WAITING_COMPETITOR игра, которая рисует соответствующий интерфейс
		// 3. Грузится второй соперник. Все то же самое, что и в пункте 2, однако возвращается статус IN_GAME обоим конкурентам
		// 4. Сервер начинает слать обоим конкурентам положение мяча и положение конкурентов из Competitor.paddleX и Competitor.paddleY
		// 5. Сервер начинает принимать положения конкурентов и записывать их в Competitor.paddleX и Competitor.paddleY
		// 6. Когда шарик улетел за границу одного из игроков, отправляется статус YOU_LOST или YOU_WON - и отрисовывается соответсвующий UI с предложением повторить. Происходит обнуление структур Competitor
		// 7. После клика на кнопки Повторить все повторяется с пункта (1) 
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
