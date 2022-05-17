#! /usr/bin/python3

import os

def replace_file(filename: str, content: bytes) -> None:
    with open(filename, "wb") as fd: #pylint: disable=invalid-name
        fd.write(content)

def read_file(filename: str, how_many: int = -1) -> bytes:
    with open(filename, "rb") as fd: #pylint: disable=invalid-name
        return fd.read(how_many)

def reset_file_states(filename: str) -> None:
    if os.path.exists(filename):
        os.remove(filename)

def add_file_to_test(filename, test_cases, optional_int: int = None, reset_files: bool = True, specific_target: str = ""):
    target = specific_target if specific_target else filename + "_target"
    if reset_files:
        reset_file_states(filename)
    other_args = []
    if optional_int != None:
        other_args = [str(optional_int).encode()]
    test_cases.append(([filename + "_original", filename] + other_args, b"")) # This will modify the filesystem
    #if optional_int != None:
    #    test_cases.append(((read_file, filename, optional_int), read_file(target)))
    #else:
    test_cases.append(((read_file, filename), read_file(target)))

if __name__ == "__main__":
    # This is utterly ugly, but I don't care
    import sys
    sys.path.append("..")
    from tests_framework import launch_tests #pylint: disable=import-error

    FILENAME = "./print_first_bytes.py"
    USAGE_TEXT = b"Usage: " + FILENAME.encode() + b" src_filepath dst_filepath [int]\n"

    TEST_CASES = [
        # Ok tests
        (["0", "Yo"], b"The file does not exist\n"),
        (["etc/hostname", "mdr"], b"The file does not exist\n"),
        (["/etc/shadow", "a"], b"Unsufficient permissions\n"),
        (["/etc/shadow", "a", "10"], b"Unsufficient permissions\n"),
        (["resources/", "mdr2"], b"Unexpected error\n"),
        (["resources", "mdr3"], b"Unexpected error\n"),

        # Fail cases
        ([], USAGE_TEXT),
        (["/etc/shadow"], USAGE_TEXT),
        (["/etc/hostname", "my_file", "ptdr"], USAGE_TEXT),
        (["1", "5", "10", "3"], USAGE_TEXT),
    ]

    add_file_to_test("resources/ls_only_4", TEST_CASES, optional_int=4)
    add_file_to_test("resources/ls", TEST_CASES)
    add_file_to_test("resources/ls", TEST_CASES, optional_int=50, reset_files=False, specific_target="resources/ls_only_50_target")
    launch_tests(FILENAME, TEST_CASES)
