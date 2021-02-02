#! /usr/bin/python3
import sys # This import allows us to interact with the command line arguments

# First, we declare a function which will do the comparison job
def biggest(int1, int2):
    if int1 > int2:
        print(int1)
    elif int2 > int1:
        print(int2)
    else:
        print("The integers are equal")

# Let's try get the command line arguments
try:
    if len(sys.argv) != 3: # Actually, the target value is 3 because the first element of argv is always the name of the script
        print("Usage: ./biggest.py integer1 integer2")
    else:
        first_int = int(sys.argv[1])
        second_int = int(sys.argv[2])
        # Now, call the function with the arguments previously retrieved
        biggest(first_int, second_int)
except ValueError:
    print("Usage: ./biggest.py integer1 integer2")
