#! /usr/bin/python3
import sys # This import allows us to interact with the command line arguments
from typing import List

USAGE_TEXT = "Usage: {} integer1 integer2 [integer3]...".format(sys.argv[0])

def biggest(l: List[int]) -> None:
    biggest = l[0]
    for i in l[1:]:
        if i > biggest:
            biggest = i

    print(biggest)

# Let's try get the command line arguments
try:
    if len(sys.argv) < 3: # Actually, the target value is 3 because the first element of argv is always the name of the script
        print(USAGE_TEXT)
    else:
        l = [int(arg) for arg in sys.argv[1:]]
        # Now, call the function with the arguments previously retrieved
        biggest(l)
except ValueError:
    print(USAGE_TEXT)
