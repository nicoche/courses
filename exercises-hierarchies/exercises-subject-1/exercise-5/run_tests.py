#! /usr/bin/python3

if __name__ == "__main__":
    # This is utterly ugly, but I don't care
    import sys
    sys.path.append("..")
    from tests_framework import launch_tests #pylint: disable=import-error

    FILENAME = "./print_first_line_better.py"
    USAGE_TEXT = b"Usage: " + FILENAME.encode() + b" filepath\n"

    def my_range(a, b = None, c = 1):
        if b == None:
            start, end, step = 0, a, c
        else:
            start, end, step = a, b, c
        return b" ".join(str(i).encode() for i in list(range(int(start), int(end), int(step)))) + b"\n"

    TEST_CASES = [
        # Ok tests
        (["0"], b"The file does not exist\n"),
        (["etc/hostname"], b"The file does not exist\n"),
        (["/etc/shadow"], b"Unsufficient permissions\n"),
        (["resources/hosts"], b"127.0.0.1\tlocalhost\n"),
        (["resources/no_line_feed.txt"], b"This is a file which does not end with a line feed."),
        (["resources/"], b"Unexpected error\n"),
        (["resources"], b"Unexpected error\n"),
        (["run_tests.py"], b"#! /usr/bin/python3\n"),

        # Fail cases
        ([], USAGE_TEXT),
        (["/etc/shadow", "a"], USAGE_TEXT),
        (["1", "0.3"], USAGE_TEXT),
        (["1", "5", "10", "3"], USAGE_TEXT),
    ]

    launch_tests(FILENAME, TEST_CASES)
