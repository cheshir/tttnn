const defaultWidth = 15;
const defaultHeight = 15;
const neededForWin = 5;
const aiURL = "/api/predict";

class Table {
    constructor(width, height) {
        let cells = [];

        for (let i = 0; i < width * height; i++) {
            cells.push({value: "", active: true})
        }

        this.width = width;
        this.height = height;
        this.cells = cells;
        this.history = [];
        this.free = width * height;
    }

    move(index, figure) {
        if (!this.cells[index].active) {
            throw new Error("Moved to busy cell")
        }

        this.cells[index].value = figure;
        this.cells[index].active = false;
        this.history.push(index);
        this.free--;
    }

    isWinning(index) {
        const expected = this.cells[index].value;

        // Look horizontal.
        let horizontal = 1 // Center.
            + this.calculateSiblings(expected, index - 1, (i) => i >= (index - index % this.width), (i) => i - 1)  // Look left.
            + this.calculateSiblings(expected, index + 1, (i) => i % this.width != 0, (i) => i + 1); // Look right.

        let vertical = 1 // Center.
            + this.calculateSiblings(expected, index - this.width, (i) => i >= 0, (i) => i - this.width) // Look top.
            + this.calculateSiblings(expected, index + this.width, (i) => i < this.cells.length, (i) => i + this.width); // Look bottom.

        let diagonal1 = 1 // Center.
            + this.calculateSiblings(expected, index - this.width - 1, (i) => i >= 0 && i % this.width != this.width - 1, (i) => i - this.width - 1) // Look top.
            + this.calculateSiblings(expected, index + this.width + 1, (i) => i < this.cells.length && i % this.width != 0, (i) => i + this.width + 1); // Look bottom.

        let diagonal2 = 1 // Center.
            + this.calculateSiblings(expected, index - this.width + 1, (i) => i >= 0 && i % this.width != 0, (i) => i - this.width + 1) // Look top.
            + this.calculateSiblings(expected, index + this.width - 1, (i) => i < this.cells.length && i % this.width != this.width - 1, (i) => i + this.width - 1); // Look Bottom.

        return [horizontal, vertical, diagonal1, diagonal2].some((count) => count >= neededForWin)
    }

    calculateSiblings(expected, start, conditionFn, incrementFn) {
        let count = 0;

        for (let i = start; conditionFn(i); i = incrementFn(i)) {
            if (this.cells[i].value !== expected) {
                break;
            }
            count++;
        }

        return count;
    }

    restart() {
        this.history.length = 0;
        // TODO Refresh with O(1)
        for (let i in this.cells) {
            this.resetCell(i);
        }
    }

    resetCell(index) {
        if (index < 0 || index >= this.cells.length) {
            throw new Error("Out of index")
        }

        this.cells[index].value = "";
        this.cells[index].active = true;
        this.free--;
    }

    isFreeCell(index) {
        if (index < 0 || index >= this.cells.length) {
            return false
        }

        return this.cells[index].value === "";
    }
}

class AI {
    constructor(url) {
        this.url = url
    }

    predictMove(table) {
        return new Promise((resolve, reject) => {
            POST(this.url, table.history, response => {
                let body = JSON.parse(response);

                resolve({
                    move: body.move,
                    probability: body.probability
                });
            }, response => {    });
        });
    }
}

const gameIsFinished = true;
const gameIsNotFinished = false;

const blackFigure = "x";
const whiteFigure = "o";

const game = new Vue({
    el: '#game',
    data: {
        ai: new AI(aiURL),
        table: new Table(defaultWidth, defaultHeight),
        playerFigure: blackFigure,
        aiFigure: whiteFigure,
        finished: false,
        resultMessage: "",
        probability: "",
        replay: {
            rawMoves: "",
            moves: [],
            currentMove: 0
        }
    },
    methods: {
        restart: function() {
            this.table.restart();
            this.resultMessage = "";
            this.finished = false;
            this.probability = "";
        },
        makeTurn: async function(index) {
            if (this.finished) {
                return
            }

            let state = this.makeSideTurn(index, this.playerFigure, this.win);
            if (state === gameIsFinished) {
                return
            }

            let prediction = await this.ai.predictMove(this.table);
            this.probability = prediction.probability;
            console.debug(`AI predicted move ${prediction.move} with probability ${prediction.probability}`);

            if (prediction.move === -1) {
                console.error("Unexpected error: AI didn't find a move", this.table);
                alert("Ooops. Caught an unexpected error. Please, restart the game and try again.");
                this.restart();
                return
            }

            this.makeSideTurn(prediction.move, this.aiFigure, this.lose);
        },
        // Makes turn and returns finish state: true – finished, false – not.
        makeSideTurn: function(moveIndex, playerFigure, onWinCb) {
            try {
                this.table.move(moveIndex, playerFigure);
            } catch (e) {
                console.error(e);
                alert(e);
                return gameIsFinished
            }

            if (this.table.isWinning(moveIndex)) {
                onWinCb();
                return gameIsFinished
            }

            if (!this.table.free) {
                this.draw();
                return gameIsFinished
            }

            return gameIsNotFinished
        },
        playGame: function() {
            this.restart();
            this.replay.moves = this.parseMoves(this.replay.rawMoves);
            this.replay.currentMove = 0;
            this.table.move(this.replay.moves[0], "x")
        },
        // TODO DRY
        toFirstMove: function() {
            this.restart();
            this.replay.currentMove = 0;
            this.table.move(this.replay.moves[0], "x");
        },
        previousMove: function() {
            this.finished = false;
            this.table.resetCell(this.replay.moves[this.replay.currentMove]);
            this.replay.currentMove--;
        },
        nextMove: function() {
            if (this.finished) {
                return
            }

            this.replay.currentMove++;
            let figure = this.replay.currentMove % 2 === 0 ? "x" : "o";

            try {
                this.table.move(this.replay.moves[this.replay.currentMove], figure);
            } catch (e) {
                console.error(e);
                alert(e);
            }

            if (this.table.isWinning(this.replay.moves[this.replay.currentMove])) {
                this.finished = true;
                this.resultMessage = (this.replay.currentMove % 2 === 0 ? "Black" : "White") + " is winning";
                return
            }

            if (!this.table.free) {
                this.draw();
                return
            }
        },
        toLastMove: function() {
            this.restart();

            for (let i = 0; i < this.replay.moves.length; i++) {
                let figure = i % 2 === 0 ? "x" : "o";
                this.table.move(this.replay.moves[i], figure);
            }

            this.replay.currentMove = this.replay.moves.length - 1;
            if (this.table.isWinning(this.replay.moves[this.replay.currentMove])) {
                this.finished = true;
                this.resultMessage = (this.replay.currentMove % 2 === 0 ? "Black" : "White") + " is winning";
                return
            }

            if (!this.table.free) {
                this.draw();
                return
            }
        },
        win: function() {
            this.finished = true;
            this.resultMessage = "WIN!"
        },
        lose: function() {
            this.finished = true;
            this.resultMessage = "Lose.."
        },
        draw: function() {
            this.finished = true;
            this.resultMessage = "Draw"
        },
        parseMoves: function(game) {
            return game.toLowerCase().split(" ").map((value) => {
                let row = value.substring(0, 1);
                let column = value.substring(1);

                return (row.charCodeAt(0) - "a".charCodeAt(0)) * this.table.width + +column - 1
            })
        }
    }
});

function GET(url, onSuccess, onError) {
    sendRequest('GET', url, null, onSuccess, onError)
}

function POST(url, payload, onSuccess, onError) {
    sendRequest('POST', url, payload, onSuccess, onError)
}

function DELETE(url, onSuccess, onError) {
    sendRequest('DELETE', url, onSuccess, onError)
}

function sendRequest(method, url, payload, onSuccess, onError) {
    let xhr = new XMLHttpRequest();
    xhr.open(method, url, true);
    xhr.setRequestHeader('Content-Type', 'application/json');
    xhr.onreadystatechange = function() {
        if (xhr.readyState !== 4) {
            return;
        }

        if (xhr.status < 400) {
            onSuccess(xhr.responseText);
        } else {
            console.error(`Error during sending request. Code: ${xhr.status}. Message: ${xhr.responseText}`);
            onError(xhr.status, xhr.responseText);
        }
    };

    let request = payload ? JSON.stringify(payload) : undefined;
    xhr.send(request);
}