let value = "o";

let game = new Vue({
    el: '#game',
    data: {
        cells: [],
    },
    methods: {
        makeTurn: function(index) {
            let nextValue = value == "x" ? "o" : "x";
            this.cells[index].value = nextValue;
            value = nextValue;
        }
    }
});

for (let i = 0; i < 100; i++) {
    game.cells.push({value: ""})
}
