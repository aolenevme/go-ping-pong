const canvas = document.getElementById("canvas");
const ctx = canvas.getContext("2d");
const ballRadius = 10;
const paddleWidth = 75;
const paddleHeight = 10;
const paddleColor = "#141414";
const ballColor = "#d0d0cf";

const DIRECTIONS = Object.freeze({
        RIGHT: "RIGHT",
        LEFT: "LEFT"
});

let ballX = canvas.width / 2; // Ball.X
let ballY = canvas.height - 30; // Ball.Y
let dx = 2; // Убрать на сервер
let dy = -2; // Убрать на сервер
let paddleTopX = (canvas.width - paddleWidth) / 2; // FirstCompetitor.X
let paddleBottomX = (canvas.width - paddleWidth) / 2; // SecondCompetitor.X

document.addEventListener("keydown", keyDownHandler, false);
runSse();

function draw() {
	ctx.clearRect(0, 0, canvas.width, canvas.height);
	drawBall();
	drawPaddle(paddleTopX, 0);
	drawPaddle(paddleBottomX, canvas.height - paddleHeight)

	if ((ballX > canvas.width - ballRadius) || ballX < ballRadius) {
		dx = -dx;
	}

	if (shouldRevertBallByY()) {
		dy = -dy;
	}

	if ((ballY + ballRadius > canvas.height - paddleHeight) || (ballY - ballRadius) < paddleHeight) {
		//alert("Game Over");
		//document.location.reload();
	}

	ballX += dx;
	ballY += dy;
}

function shouldRevertBallByY() {
	const isBallTouchedTopPaddle = ballX >= paddleTopX && ballX < paddleTopX + paddleWidth && ballY - ballRadius <= paddleHeight;
	const isBallTouchedBottomPaddle = ballX >= paddleBottomX && ballX < paddleBottomX + paddleWidth && ballY + ballRadius >= canvas.height - paddleHeight;

	return isBallTouchedTopPaddle || isBallTouchedBottomPaddle;
}

function drawBall() {
	ctx.beginPath();
	ctx.arc(ballX, ballY, ballRadius, 0, Math.PI * 2);
	ctx.fillStyle = ballColor;
	ctx.fill();
	ctx.closePath();
}

function drawPaddle(paddleX = 0, paddleY = 0) {
	ctx.beginPath();
	ctx.rect(paddleX, paddleY, paddleWidth, paddleHeight);
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
		console.log(data);
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
