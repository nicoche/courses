#! /usr/bin/python3

import sys # This import allows us to interact with the command line arguments

# Let's try get the command line arguments
USAGE_TEXT = "Usage: {} filepath".format(sys.argv[0])

def print_file_content(filepath: str) -> None:
    try:
        with open(filepath, "a") as fd: #pylint: disable=invalid-name
            fd.write("hello")
    except FileNotFoundError:
        print("The file does not exist")
    except PermissionError:
        print("Unsufficient permissions")
    except OSError:
        print("Unexpected error")

COUNT_ARGS = len(sys.argv)
if COUNT_ARGS != 2 or COUNT_ARGS > 4:
    print(USAGE_TEXT)
else:
    print_file_content(sys.argv[1])
