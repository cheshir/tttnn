const width = 10;
const height = 10;
const neededForWin = 5;

let value = "o";

let game = new Vue({
    el: '#game',
    data: {
        cells: [],
        finished: false,
        resultMessage: "",
    },
    methods: {
        init: function() {
            for (let i = 0; i < width * height; i++) {
                this.cells.push({value: "", active: true})
            }

        },
        refresh: function() {
            // TODO Refresh with O(1)
            for (let i in this.cells) {
                this.cells[i].value = "";
                this.cells[i].active = true;
            }

            this.resultMessage = "";
            this.finished = false;
        },
        makeTurn: function(index) {
            if (!this.cells[index].active || this.finished) {
                return
            }

            // TODO Dummy.
            let nextValue = value == "x" ? "o" : "x";
            this.cells[index].value = nextValue;
            value = nextValue;
            this.cells[index].active = false;

            if (this.isWinning(index)) {
                this.win()
            }
        },
        isWinning(index) {
            const expected = this.cells[index].value;

            // Look horizontal.
            let horizontal = 1 // Center.
                + this.calculateSiblings(expected, width % index, (i) => i >= 0, (i) => i - 1)  // Look left.
                + this.calculateSiblings(expected, index + 1, (i) => i % width != 0, (i) => i + 1); // Look right.

            let vertical = 1 // Center.
                + this.calculateSiblings(expected, index - width, (i) => i >= 0, (i) => i - width) // Look top.
                + this.calculateSiblings(expected, index + width, (i) => i < this.cells.length, (i) => i + width); // Look bottom.

            let diagonal1 = 1 // Center.
                + this.calculateSiblings(expected, index - width - 1, (i) => i >= 0 && i % width != width - 1, (i) => i - width - 1) // Look top.
                + this.calculateSiblings(expected, index + width + 1, (i) => i < this.cells.length && i % width != 0, (i) => i + width + 1); // Look bottom.

            let diagonal2 = 1 // Center.
                + this.calculateSiblings(expected, index - width + 1, (i) => i >= 0 && i % width != 0, (i) => i - width + 1) // Look top.
                + this.calculateSiblings(expected, index + width - 1, (i) => i < this.cells.length && i % width != width - 1, (i) => i + width - 1); // Look Bottom.

            return [horizontal, vertical, diagonal1, diagonal2].some((count) => count >= neededForWin)
        },
        calculateSiblings(expected, start, conditionFn, incrementFn) {
            let count = 0;

            for (let i = start; conditionFn(i); i = incrementFn(i)) {
                if (this.cells[i].value != expected) {
                    break;
                }
                count++;
            }

            return count;
        },
        win: function() {
            this.finished = true;
            this.resultMessage = "WIN!"
        },
        lose: function() {
            this.finished = true;
            this.resultMessage = "Lose.."
        }
    }
});

game.init();
