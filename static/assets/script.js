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
            for (let i = 0; i < 100; i++) {
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
