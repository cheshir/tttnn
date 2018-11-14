import numpy as np
import tensorflow as tf
from converter_to_learning_data import read_games_file, build_learning_data

from keras import backend as K
from keras.models import Sequential
from keras.layers import Dense, Flatten, Activation, Dropout
from keras.layers.convolutional import Convolution2D, MaxPooling2D

TEST_DATA_PERCENT = 0.01

# Init
np.random.seed(42)

raw_games = read_games_file("../data/prepared/full_games_renjunet.csv")
learning_data_x, learning_data_y = build_learning_data(raw_games)

# Split to train and test data.
testSplitIndex = int(len(learning_data_x) * TEST_DATA_PERCENT)
learn_x = np.array(learning_data_x[testSplitIndex:])
learn_y = np.array(learning_data_y[testSplitIndex:])
test_x = np.array(learning_data_x[:testSplitIndex])
test_y = np.array(learning_data_y[:testSplitIndex])

# We should create session to save model later.
sess = tf.Session()
K.set_session(sess)

# Create model based on
# @link http://cs231n.stanford.edu/reports/2016/pdfs/109_Report.pdf
model = Sequential()

model.add(Convolution2D(1024, 5, strides=1, padding="same", input_shape=(15, 15, 1), activation="relu")) # name=conv2d_1_input
model.add(MaxPooling2D(pool_size=3, strides=3))
model.add(Dropout(0.25))

model.add(Convolution2D(256, 5, strides=1, padding="same", activation="relu"))
model.add(Convolution2D(128, 5, strides=1, padding="same", activation="relu"))
model.add(Dropout(0.25))

# Convert from matrix to vector.
model.add(Flatten())

# 3 Fully connected layers.
model.add(Dense(900, activation="relu"))
model.add(Dense(450, activation="relu"))
model.add(Dropout(0.3))
model.add(Dense(225, activation="softmax", name="output"))

print "Compiling..."
model.compile(loss="categorical_crossentropy", optimizer="SGD", metrics=["accuracy"])

print "Started learning"
model.fit(learn_x[:100], learn_y[:100], batch_size=8, epochs=1, validation_split=0.1, shuffle=True)

model.save_weights("../data/models/stanford.cnn.1.h5")
scores = model.evaluate(test_x, test_y, verbose=1)
print "Test data accuracy %.2f%%" % (scores[1] * 100)

print "Save model for go"
# Use TF to save the graph model instead of Keras save model to load it in Golang.
builder = tf.saved_model.builder.SavedModelBuilder("3tnn")
# Tag the model, required for Go.
builder.add_meta_graph_and_variables(sess, ["tag"])
builder.save()
sess.close()
