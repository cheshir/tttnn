package main

import (
	"log"

	tf "github.com/tensorflow/tensorflow/tensorflow/go"
)

func main() {
	model, err := tf.LoadSavedModel("../../data/models/stanford.cnn.1.h5", nil, nil)
	if err != nil {
		log.Fatal(err)
	}

	defer model.Session.Close()

	println("done")
	//println(model.Graph.Operations())
}
