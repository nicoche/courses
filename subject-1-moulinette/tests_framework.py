from subprocess import run, PIPE, CalledProcessError, TimeoutExpired
from typing import List, Tuple

def launch_command(filename: str, args: List[str]) -> Tuple[bool, bytes]:
    try:
        done = run([filename] + args,
                   stderr=PIPE, stdout=PIPE,
                   timeout=2, check=True)

        if done.stderr:
            return False, done.stderr
        return True, done.stdout

    except FileNotFoundError:
        return False, b"File not found. Make sure your script is correctly named"
    except CalledProcessError as err:
        return False, b"\n" + err.stderr
    except TimeoutExpired as err:
        return False, b"Took too long to execute!"
    except PermissionError:
        return False, b"Permission denied. Are you sure your script is executable?"

def colored_print(test_name: str, success: bool) -> None:
    if success:
        print("\x1b[7;32;40mTest {} - {}\x1b[0m".format(test_name, "PASSED"))
    else:
        print("\x1b[7;31;40mTest {} - {}\x1b[0m".format(test_name, "FAILED"))

def launch_tests(filename: str, test_cases: List[Tuple[List, bytes]]) -> None:
    for test_arguments, expected_output in test_cases:
        # This is a file edition check
        if isinstance(test_arguments, tuple):
            completed_okay, output = True, test_arguments[0](*test_arguments[1:])
        # This is a subprocess.run() need
        else:
            completed_okay, output = launch_command(filename, test_arguments)

        if completed_okay and output == expected_output:
            colored_print(test_arguments, True)
            continue
        colored_print(test_arguments, False)
        if completed_okay:
            try:
                print("--> Got unexpected output, obtained {} when {} was "
                      "expected\n".format(repr(output.decode()), repr(expected_output.decode())), end='')
            except UnicodeDecodeError:
                print("--> Got unexpected output, obtained {} when {} was "
                      "expected\n".format(output, expected_output), end='')
        else:
            print("--> Your script encountered an error: {}\n".format(output.decode()), end='')
        return

    print("\nAll tests passed! Good job :). Make sure you did not forget any not tested case")
