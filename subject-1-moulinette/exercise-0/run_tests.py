#! /usr/bin/python3

if __name__ == "__main__":
    # This is utterly ugly, but I don't care
    import sys
    sys.path.append("..")
    from tests_framework import launch_tests #pylint: disable=import-error

    FILENAME = "./biggest.py"
    USAGE_TEXT = b"Usage: ./biggest.py integer1 integer2\n"
    TEST_CASES = [
        (["1", "2"], b"2\n"),
        (["1123123", "2"], b"1123123\n"),
        (["-17", "13"], b"13\n"),
        (["0.3", "13"], USAGE_TEXT),
        (["13"], USAGE_TEXT),
        (["0.1"], USAGE_TEXT),
        (["1", "2", "3", "-8"], USAGE_TEXT),
        (["mdr", "3"], USAGE_TEXT),
        (["0", "3"], b"3\n"),
        (["0.7", "1"], USAGE_TEXT),
        (["19", "19"], b"The integers are equal\n"),
        (["0.7", "0.3"], USAGE_TEXT),
    ]

    launch_tests(FILENAME, TEST_CASES)
