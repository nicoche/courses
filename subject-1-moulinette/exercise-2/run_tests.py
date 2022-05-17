#! /usr/bin/python3

if __name__ == "__main__":
    # This is utterly ugly, but I don't care
    import sys
    sys.path.append("..")
    from tests_framework import launch_tests #pylint: disable=import-error

    FILENAME = "./longest_words.py"
    USAGE_TEXT = b"Usage: " + FILENAME.encode() + b" arg1 [arg2]...\n"
    TEST_CASES = [
        (["bonjour", "au-revoir", "abcdefghijklmnopqrstuvwxyz"], b"abcdefghijklmnopqrstuvwxyz\n"),
        (["bonjour", "abcdefghijklmnopqrstuvwxyz", "au-revoir"], b"abcdefghijklmnopqrstuvwxyz\n"),
        (["12913213", "-2", "mdr3333"], b"12913213\n"),
        (["a", "z", "e"], b"a,z,e\n"),
        (["abcd", "efgh", "ijkl"], b"abcd,efgh,ijkl\n"),
        (["abcd", "efgh", "ijkl", "fsociety"], b"fsociety\n"),
        (["abcd"], b"abcd\n"),
        ([], USAGE_TEXT),
    ]

    launch_tests(FILENAME, TEST_CASES)
