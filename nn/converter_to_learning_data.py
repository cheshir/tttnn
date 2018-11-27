"""
Potential use:
raw_games = read_games_file("../data/prepared/full_games_renjunet.csv")
learning_data_x, learning_data_y = build_learning_data(raw_games)
"""
import csv
import numpy as np

from keras.utils import np_utils

# Game results that is used in raw data.
RESULT_BLACK_WINS = "1"
RESULT_DRAW = "0.5"
RESULT_WHITE_WINS = "0"

# Board sizes.
WIDTH = 15
HEIGHT = 15

BLACK_TURN = 0
WHITE_TURN = 1

MOVE_DIMENSION = 2
BLACK_MOVE = [1, 0]
WHITE_MOVE = [0, 1]
EMPTY_MOVE = [0, 0]

MIN_MOVES = 14

def read_games_file(filepath):
    with open(filepath) as datasource:
        reader = csv.reader(datasource, delimiter=',')
        return [(row[0], row[1].split()) for row in reader]


def build_learning_data(games):
    """
    :param games: is a cortege where first element is a game result (from black side)
    and second element are game moves.
    :return: data prepared for learning.
    """
    # Ignore known positions.
    known_positions = set()
    x = []
    y = []
    for game in games:
        result, moves = game
        if len(moves) < MIN_MOVES:
            continue
        game_x, game_y = convert_game_to_learning_data(result, moves, known_positions)
        x += game_x
        y += game_y

    # Format data.
    formatted_x = np.asarray(x)
    formatted_y = np_utils.to_categorical(y, WIDTH*HEIGHT)

    return formatted_x, formatted_y


def convert_game_to_learning_data(result, moves, known_positions):
    """
    Convert game to learning data.
    :param result: game result from black side.
    :param moves: string with space separated list of moves.
    :return: x and y where x is a game board and y is an expected result.
    """
    x = []
    y = []
    board = np.zeros([WIDTH, HEIGHT, MOVE_DIMENSION])
    index = 0
    position = ""

    for move in moves:
        active_side = BLACK_TURN if index % 2 == 0 else WHITE_TURN
        coord = convert_move_to_position(move)

        position += move
        if position not in known_positions:
          x.append(board.copy())
          y.append(vectorize_coord(coord))

        board[coord[0]][coord[1]] = BLACK_MOVE if active_side == BLACK_TURN else WHITE_MOVE
        index += 1
    return x, y


def convert_move_to_position(move):
    x = ord(move[:1]) - ord('a')
    y = int(move[1:2]) - 1

    return x, y


def vectorize_coord(coord):
    x, y = coord

    return float(x * WIDTH + y + 1) # Added 1 to make it categorizable.
