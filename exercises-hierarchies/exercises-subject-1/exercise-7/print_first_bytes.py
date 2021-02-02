#! /usr/bin/python3

import sys # This import allows us to interact with the command line arguments

# Let's try get the command line arguments
USAGE_TEXT = "Usage: {} src_filepath dst_filepath [int]".format(sys.argv[0])

def write_file_content(src_filepath: str, dst_filepath, number_of_bytes_to_read: int = -1) -> None:
    try:
        with open(src_filepath, "rb") as fd: #pylint: disable=invalid-name
            bytes_read = fd.read(number_of_bytes_to_read)
        with open(dst_filepath, "wb") as fd: #pylint: disable=invalid-name
            fd.write(bytes_read)
    except FileNotFoundError:
        print("The file does not exist")
    except PermissionError:
        print("Unsufficient permissions")
    except OSError:
        print("Unexpected error")

COUNT_ARGS = len(sys.argv)
if COUNT_ARGS < 3 or COUNT_ARGS > 4:
    print(USAGE_TEXT)
else:
    if COUNT_ARGS == 4:
        try:
            COUNT_BYTES_TO_READ = int(sys.argv[3])
        except ValueError:
            print(USAGE_TEXT)
        else:
            write_file_content(sys.argv[1], sys.argv[2], COUNT_BYTES_TO_READ)
    else:
        write_file_content(sys.argv[1], sys.argv[2], -1)
