import requests
import json
import os
from prompt_toolkit import PromptSession
from prompt_toolkit.completion import Completer, Completion

# Base URL for the REST API
BASE_URL = "http://localhost:9000"

current_directory = "/"

def send_query(url):
    # Send GET request
    resp = requests.get(url)
    resp.raise_for_status()
    return resp.content

def send_delete(url):
    # Send DELETE request
    resp = requests.delete(url)
    resp.raise_for_status()
    return resp.content

def directory_exists_on_server(directory_name):
    url = f"{BASE_URL}/{directory_name}?type=directory"
    resp = requests.head(url)
    return resp.status_code == requests.codes.ok

def show_entries(response, long_format):
    dir_entries = json.loads(response)
    if long_format:
        for entry in dir_entries:
            entry_type = "d" if entry["type"] == "directory" else "-"
            print(f"{entry_type} {entry['name']}")
    else:
        print(" ".join(entry["name"] for entry in dir_entries))

def list_cmd(long_format):
    url = f"{BASE_URL}/{current_directory}?type=directory&operation=list"
    response = send_query(url)
    if response != b'null':
        show_entries(response, long_format)

def change_directory(args):
    global current_directory
    dst_directory = os.path.join(current_directory, args[0])
    if directory_exists_on_server(dst_directory):
        current_directory = dst_directory
    else:
        print(f"Invalid directory: {dst_directory}")

def show_current_directory(args):
    print(current_directory)

def remove_directory(args):
    directory_name = args[0]
    url = f"{BASE_URL}/{current_directory}/{directory_name}?type=directory"
    print(f"Removing directory {directory_name} -> {url}")
    send_delete(url)

def completer(document):
    suggestions = []

    # Populate suggestions from commands
    suggestions.extend([cmd.name for cmd in [list_cmd, change_directory, remove_directory, show_current_directory]])

    # Add flags to suggestions
    flags = list_cmd.flags
    suggestions.extend([f"--{flag.name}" for flag in flags])

    # Add "exit" suggestion
    suggestions.append("exit")

    return [Completion(suggestion, start_position=-document.cursor_position) for suggestion in suggestions]


class SimpleCompleter(Completer):
    def __init__(self, words):
        self.words = words

    def get_completions(self, document, complete_event):
        word_before_cursor = document.get_word_before_cursor()
        matches = [word for word in self.words if word.startswith(word_before_cursor)]
        for m in matches:
            yield Completion(m, start_position=-len(word_before_cursor))

def main():
    completer = SimpleCompleter(['ls', 'cd', 'pwd', 'rmdir', 'exit'])
    session = PromptSession(completer=completer)

    while True:
        try:
            user_input = session.prompt('> ')

            args = user_input.split()
            if (len(args) == 0):
                continue

            cmd_name = args[0]
            cmd_args = args[1:]

            if cmd_name == "ls":
                list_cmd("-l" in cmd_args)
            elif cmd_name == "cd":
                change_directory(cmd_args)
            elif cmd_name == "pwd":
                show_current_directory(cmd_args)
            elif cmd_name == "rmdir":
                remove_directory(cmd_args)
            elif cmd_name == "exit":
                print("Exiting...")
                break
            else:
                print("Unknown command:", cmd_name)

        except KeyboardInterrupt:
            continue
        except EOFError:
            break

if __name__ == '__main__':
    main()
