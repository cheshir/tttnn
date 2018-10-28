import numpy as np
from converter_to_learning_data import read_games_file, build_learning_data

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

# Create model based on
# @link http://cs231n.stanford.edu/reports/2016/pdfs/109_Report.pdf
model = Sequential()

model.add(Convolution2D(1024, 5, strides=1, padding="same", input_shape=(15, 15, 1), activation="relu"))
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
model.add(Dense(225, activation="softmax"))

print "Compiling..."
model.compile(loss="categorical_crossentropy", optimizer="SGD", metrics=["accuracy"])

print "Started learning"
model.fit(learn_x, learn_y, batch_size=8, epochs=30, validation_split=0.1, shuffle=True)

model.save_weights("../data/models/stanford.cnn.1.h5")
scores = model.evaluate(test_x, test_y, verbose=1)
print "Test data accuracy %.2f%%" % (scores[1] * 100)
