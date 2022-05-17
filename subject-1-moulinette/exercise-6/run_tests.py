#! /usr/bin/python3

import os

def replace_file(filename: str, content: bytes) -> None:
    with open(filename, "wb") as fd:
        fd.write(content)

def read_file(filename: str) -> bytes:
    with open(filename, "rb") as fd:
        return b"".join(fd.readlines())

def reset_file_states(filename: str) -> None:
    if os.path.exists(filename + "_original"):
        replace_file(filename, read_file(filename + "_original"))
    elif os.path.exists(filename):
        os.remove(filename)

def add_file_to_test(filename, test_cases):
    reset_file_states(filename)
    test_cases.append(([filename], b"")) # This will modify toto.txt
    test_cases.append(((read_file, filename), read_file(filename + "_target")))

if __name__ == "__main__":
    # This is utterly ugly, but I don't care
    import sys
    sys.path.append("..")
    from tests_framework import launch_tests #pylint: disable=import-error

    FILENAME = "./hello.py"
    USAGE_TEXT = b"Usage: " + FILENAME.encode() + b" filepath\n"

    TEST_CASES = [
        (["etc/hostname"], b"The file does not exist\n"),
        (["/etc/shadow"], b"Unsufficient permissions\n"),
        (["resources/"], b"Unexpected error\n"),
        (["resources"], b"Unexpected error\n"),

        # Fail cases
        ([], USAGE_TEXT),
        (["/etc/shadow", "a"], USAGE_TEXT),
        (["1", "0.3"], USAGE_TEXT),
        (["1", "5", "10", "3"], USAGE_TEXT),
    ]

    add_file_to_test("resources/0", TEST_CASES)
    add_file_to_test("resources/no_line_feed.txt", TEST_CASES)
    add_file_to_test("resources/toto.txt", TEST_CASES)

    launch_tests(FILENAME, TEST_CASES)
