import pathlib
import stat
import subprocess
import sys


def run_go_mod_tidy() -> bool:
    try:
        subprocess.run(["go", "mod", "tidy"], capture_output=True, check=True)
        return True
    except Exception:
        return False


def set_script_permissions() -> None:
    for path in pathlib.Path(".").glob("**/*.sh"):
        path.chmod(path.stat().st_mode | stat.S_IEXEC)


if __name__ == "__main__":
    set_script_permissions()
    if not run_go_mod_tidy():
        print("ERROR: go mod tidy failed.", file=sys.stderr)
