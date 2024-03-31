import requests
import json
import os
from prompt_toolkit import PromptSession
from prompt_toolkit.completion import Completer, Completion
from prompt_toolkit.document import Document

# Base URL for the REST API
BASE_URL = "http://localhost:9000"

current_directory = "/"

commands_with_args = {
    'ls': ['-l', '-a', '-h'],
    'cd': [],
    'pwd': [],
    'rmdir': [],
    'exit': []
}

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
    commands_with_args['cd'] = []
    if long_format:
        for entry in dir_entries:
            if entry["type"] == "directory":
                commands_with_args['cd'].append(entry["name"])
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

def fetch_commands_with_args():
    # This is just an example implementation.
    # You can replace it with your own logic to fetch commands and their arguments dynamically.
    return commands_with_args

class CommandCompleter(Completer):
    def __init__(self):
        self.commands_with_args = fetch_commands_with_args()

    def get_completions(self, document, complete_event):
        text_before_cursor = document.text_before_cursor
        parts = text_before_cursor.split()

        if len(parts) == 1:
            for command in self.commands_with_args.keys():
                if command.startswith(parts[0]):
                    yield Completion(command, -len(parts[0]))
        elif len(parts) == 2:
            command = parts[0]
            if command in self.commands_with_args:
                for arg in self.commands_with_args[command]:
                    if arg.startswith(parts[1]):
                        yield Completion(arg, -len(parts[1]))

def main():

    completer = CommandCompleter()
    session = PromptSession(completer=completer, complete_while_typing=True)

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
