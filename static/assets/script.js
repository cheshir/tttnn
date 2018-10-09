const width = 10;
const height = 10;
const neededForWin = 5;

// TODO Rewrite to class.
function Table(width, height) {
    let cells = [];

    for (let i = 0; i < width * height; i++) {
        cells.push({value: "", active: true})
    }

    this.width = width;
    this.height = height;
    this.cells = cells;
    this.free = width * height;
}

// FIXME Fast clicking.
Table.prototype.move = function(index, figure) {
    if (!this.cells[index].active) {
        return
    }

    this.cells[index].value = figure;
    this.cells[index].active = false;
    this.free--;
};

Table.prototype.isWinning = function(index) {
    const expected = this.cells[index].value;

    // Look horizontal.
    let horizontal = 1 // Center.
        + this.calculateSiblings(expected, index % this.width, (i) => i >= 0, (i) => i - 1)  // Look left.
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
};

Table.prototype.calculateSiblings = function(expected, start, conditionFn, incrementFn) {
    let count = 0;

    for (let i = start; conditionFn(i); i = incrementFn(i)) {
        if (this.cells[i].value != expected) {
            break;
        }
        count++;
    }

    return count;
};

Table.prototype.restart = function() {
    // TODO Refresh with O(1)
    for (let i in this.cells) {
        this.cells[i].value = "";
        this.cells[i].active = true;
    }
};

Table.prototype.isFreeCell = function(index) {
    if (index < 0 || index >= this.cells.length) {
        return false
    }

    return this.cells[index].value === "";
};

const AI = {
    figure: "o",
    finders: [
        // Left.
        (table, playerTurn, figure, i) => {
            let turnIndex = playerTurn - i;
            if (!table.isFreeCell(turnIndex)) {
                return -1
            }

            table.move(turnIndex, figure);
            return turnIndex
        },
        // Right.
        (table, playerTurn, figure, i) => {
            let turnIndex = playerTurn + i;
            if (!table.isFreeCell(turnIndex)) {
                return -1
            }

            table.move(turnIndex, figure);
            return turnIndex
        },
        // Top.
        (table, playerTurn, figure, i) => {
            for (let j = playerTurn - table.width - i; j < playerTurn - table.width + i; j++) {
                if (table.isFreeCell(j)) {
                    table.move(j, figure);
                    return j
                }
            }

            return -1
        },
        // Bottom.
        (table, playerTurn, figure, i) => {
            for (let j = playerTurn + table.width - i; j < playerTurn + table.width + i; j++) {
                if (table.isFreeCell(j)) {
                    table.move(j, figure);
                    return j
                }
            }

            return -1
        }
    ],
    move: function (table, playerTurn) {
        // Find nearest to player move free cell.
        for (let i = 1; i < width; i++) {
            let finders = this.finders.shuffle();

            for (let i in finders) {
                let position = finders[i](table, playerTurn, this.figure, i);
                if (position !== -1) {
                    return position
                }
            }
        }

        return -1
    }
};

const game = new Vue({
    el: '#game',
    data: {
        table: new Table(width, height),
        playerFigure: "x",
        finished: false,
        resultMessage: "",
    },
    methods: {
        restart: function() {
            this.table.restart();
            this.resultMessage = "";
            this.finished = false;
        },
        makeTurn: function(index) {
            if (this.finished) {
                return
            }

            this.table.move(index, this.playerFigure);

            if (this.table.isWinning(index)) {
                this.win();
                return
            }

            if (!this.table.free) {
                this.draw();
                return
            }

            let position = AI.move(this.table, index);
            console.debug("AI: ", position);

            if (position === -1) {
                console.error("Unexpected error: AI didn't find a move", this.table);
                alert("Ooops. Caught an unexpected error. Please, restart the game and try again.");
                this.restart();
                return
            }

            debugger;
            if (this.table.isWinning(position)) {
                this.lose();
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
        }
    }
});

// FIXME move to other file.
Array.prototype.shuffle = function() {
    let input = this;

    for (let i = input.length-1; i >=0; i--) {
        let randomIndex = Math.floor(Math.random()*(i+1));
        let itemAtIndex = input[randomIndex];

        input[randomIndex] = input[i];
        input[i] = itemAtIndex;
    }
    return input;
};
