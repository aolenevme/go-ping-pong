const canvas = document.getElementById("canvas");
const ctx = canvas.getContext("2d");
const paddleColor = "#141414";
const ballColor = "#d0d0cf";

const DIRECTIONS = Object.freeze({
        RIGHT: "RIGHT",
        LEFT: "LEFT"
});

let dx = 2; // Убрать на сервер
let dy = -2; // Убрать на сервер
let game = {};

document.addEventListener("keydown", keyDownHandler, false);
runSse();

function draw() {
	ctx.clearRect(0, 0, game.canvasWidth, game.canvasHeight);
	drawBall();
	drawPaddle(game.paddleTopX, 0);
	drawPaddle(game.paddleBottomX, game.canvasHeight - game.paddleHeight)

	if ((game.ballX > game.canvasWidth - game.ballRadius) || game.ballX < game.ballRadius) {
		dx = -dx;
	}

	if (shouldRevertBallByY()) {
		dy = -dy;
	}

	if ((game.ballY + game.ballRadius > game.canvasHeight - game.paddleHeight) || (game.ballY - game.ballRadius) < game.paddleHeight) {
		//alert("Game Over");
		//document.location.reload();
	}

	game.ballX += dx;
	game.ballY += dy;
}

function shouldRevertBallByY() {
	const isBallTouchedTopPaddle = game.ballX >= game.paddleTopX && game.ballX < game.paddleTopX + game.paddleWidth && game.ballY - game.ballRadius <= game.paddleHeight;
	const isBallTouchedBottomPaddle = game.ballX >= game.paddleBottomX && game.ballX < game.paddleBottomX + game.paddleWidth && game.ballY + game.ballRadius >= game.canvasHeight - game.paddleHeight;

	return isBallTouchedTopPaddle || isBallTouchedBottomPaddle;
}

function drawBall() {
	ctx.beginPath();
	ctx.arc(game.ballX, game.ballY, game.ballRadius, 0, Math.PI * 2);
	ctx.fillStyle = ballColor;
	ctx.fill();
	ctx.closePath();
}

function drawPaddle(paddleX = 0, paddleY = 0) {
	ctx.beginPath();
	ctx.rect(paddleX, paddleY, game.paddleWidth, game.paddleHeight);
	ctx.fillStyle = paddleColor;
	ctx.fill();
	ctx.closePath();
}

async function keyDownHandler(e = {}) {
        const key = e.key;
        let direction = DIRECTIONS.RIGHT;

        if (key === "Right" || key === "ArrowRight") {
                direction = DIRECTIONS.RIGHT;
        } else if (key === "Left" || key === "ArrowLeft") {
                direction = DIRECTIONS.LEFT;
        }

        await putData("api/v1/sse", { direction });
}

function runSse() {
	const sse = new EventSource("api/v1/sse");
	
	sse.addEventListener("open", () => {
		console.log("Stream is open");
	});

	sse.addEventListener("error", e => {
		const eventSourceState = e.target.readyState;

		switch (eventSourceState) {
			case EventSource.CONNECTING: 
				console.log("Reconnecting");
				break;
			case EventSource.CLOSED:
				console.log("Connections failed, will not reconnect");
				break;
			default:
				console.log("Unknown error");
		}
	});

	sse.addEventListener("message", e => {
		const data = JSON.parse(e.data);
		game = { ...game, ...data };
		console.dir(game);
		draw();
	});
}

async function putData(url = "", data = {}) {
  const response = await fetch(url, {
	  method: "PUT",
	  headers: {
		  "Content-Type": "application/json"
	  },
	  body: JSON.stringify(data)
  });

  return await response.json();
}
