#! /usr/bin/python3

import sys # This import allows us to interact with the command line arguments

# Let's try get the command line arguments
USAGE_TEXT = "Usage: {} arg1 [arg2]...".format(sys.argv[0])

if len(sys.argv) < 2: # Actually, the target value is 3 because the first element of argv is always the name of the script
    print(USAGE_TEXT)
else:
    longest_words = []
    for i in sys.argv[1:]:
        if not longest_words or len(i) > len(longest_words[0]):
            longest_words = [i]
        elif len(longest_words[0]) == len(i):
            longest_words.append(i)

    print(",".join(longest_words))
