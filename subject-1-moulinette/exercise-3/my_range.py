#! /usr/bin/python3

import sys # This import allows us to interact with the command line arguments

# Let's try get the command line arguments
USAGE_TEXT = "Usage: {} int1 [int2] [int3]".format(sys.argv[0])

def print_my_range(begin: int, end: int, step: int) -> None:
    if step == 0:
        raise ValueError
    to_print = []

    positive_step = step > 0
    condition = lambda begin, end: begin < end if positive_step else begin > end
    while condition(begin, end):
        to_print.append(str(begin))
        begin += step
    print(" ".join(to_print))

COUNT_ARGS = len(sys.argv)
if COUNT_ARGS < 2 or COUNT_ARGS > 4:
    print(USAGE_TEXT)
else:
    try:
        if COUNT_ARGS == 2:
            BEGIN, END = 0, int(sys.argv[1])
        else:
            BEGIN, END = int(sys.argv[1]), int(sys.argv[2])
        STEP = int(sys.argv[3]) if len(sys.argv) == 4 else 1
        print_my_range(BEGIN, END, STEP)
    except ValueError:
        print(USAGE_TEXT)
