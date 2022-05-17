#! /usr/bin/python3

if __name__ == "__main__":
    # This is utterly ugly, but I don't care
    import sys
    sys.path.append("..")
    from tests_framework import launch_tests #pylint: disable=import-error

    FILENAME = "./my_range.py"
    USAGE_TEXT = b"Usage: " + FILENAME.encode() + b" int1 [int2] [int3]\n"

    def my_range(a, b = None, c = 1):
        if b == None:
            start, end, step = 0, a, c
        else:
            start, end, step = a, b, c
        return b" ".join(str(i).encode() for i in list(range(int(start), int(end), int(step)))) + b"\n"

    TEST_CASES = [
        # Ok tests
        (["0", "10", "1"], my_range(0, 10, 1)),
        (["0", "10", "2"], my_range(0, 10, 2)),
        (["0", "10", "9"], my_range(0, 10, 9)),
        (["0", "10", "15"], my_range(0, 10, 15)),
        (["0", "10", "-15"], my_range(0, 10, -15)),
        (["10", "5", "-1"], my_range(10, 5, -1)),
        (["0", "10", "-1"], my_range(0, 10, -1)),
        (["10", "0", "-1"], my_range(10, 0, -1)),
        (["222", "999", "33"], my_range(222, 999, 33)),
        (["-222", "-999", "-33"], my_range(-222, -999, -33)),
        (["-222", "-999", "33"], my_range(-222, -999, 33)),
        (["-999", "-222", "33"], my_range(-999, -222, 33)),
        (["-999", "-222", "-33"], my_range(-999, -222, -33)),
        (["999", "-222", "-33"], my_range(999, -222, -33)),
        (["999", "-222", "33"], my_range(999, -222, 33)),
        (["1", "-1", "1"], my_range(1, -1, 1)),

        # Fail cases
        ([], USAGE_TEXT),
        (["1", "a"], USAGE_TEXT),
        (["1", "0.3"], USAGE_TEXT),
        (["1", "5", "10", "3"], USAGE_TEXT),

        # Sneaky cases
        (["10"], my_range(10)),
        (["-10"], my_range(-10)),
    ]

    launch_tests(FILENAME, TEST_CASES)
