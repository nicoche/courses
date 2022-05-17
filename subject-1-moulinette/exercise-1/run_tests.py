#! /usr/bin/python3

if __name__ == "__main__":
    # This is utterly ugly, but I don't care
    import sys
    sys.path.append("..")
    from tests_framework import launch_tests #pylint: disable=import-error

    FILENAME = "./biggest_advanced.py"
    USAGE_TEXT = b"Usage: " + FILENAME.encode() + b" integer1 integer2 [integer3]...\n"
    TEST_CASES = [
        # Same tests as biggest.py
        (["1", "2"], b"2\n"),
        (["1123123", "2"], b"1123123\n"),
        (["-17", "13"], b"13\n"),
        (["0.3", "13"], USAGE_TEXT),
        (["13"], USAGE_TEXT),
        (["0.1"], USAGE_TEXT),
        (["1", "2", "3", "-8"], b"3\n"), #This one changes, tho
        (["mdr", "3"], USAGE_TEXT),
        (["0", "3"], b"3\n"),
        (["0.7", "1"], USAGE_TEXT),
        (["19", "19"], b"19\n"), #This one changes, too
        (["0.7", "0.3"], USAGE_TEXT),

        # New tests
        (["1", "2", "3", "4", "10"], b"10\n"),
        (["-1", "2", "3", "4", "-10"], b"4\n"),
        (["12913213", "-2", "3333"], b"12913213\n"),
        ([str(i) for i in range(10000)], b"9999\n"),
        (["12913213", "-2", "mdr3333"], USAGE_TEXT),
        (["12913213", "-2", "17", "-44", "3.8"], USAGE_TEXT),
        (["12913213", "12913213", "12913213", "12913213", "12913213"], b"12913213\n"),
        (["3", "3333", "333333", "33"], b"333333\n"),
    ]

    launch_tests(FILENAME, TEST_CASES)
