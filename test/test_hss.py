import requests
import hashlib
import os

# Replace these variables with your actual values
BASE_URL = "http://localhost:9000/"
DIRECTORY_NAME = "mydir"
TEST_FILE = "file.txt"
DATA_STORE = "/tmp/data-store"
EXPECTED_FILE_MD5 = "5dd39cab1c53c2c77cd352983f9641e1"

# Function to log commands
def log_cmd(description):
    print(f"Executing command: {description}")

# Get directory info
log_cmd("Get directory info")
response = requests.post(BASE_URL + DIRECTORY_NAME + "/?type=directory")

# List directory files
log_cmd("List directory files")
response = requests.get(BASE_URL + DIRECTORY_NAME + "/?type=directory&operation=list")
print(response)
print(response.json())  # Assuming the response is JSON

# File Upload
log_cmd("File Upload")
with open(TEST_FILE, 'rb') as file:
    response = requests.post(BASE_URL + DIRECTORY_NAME + f"/{TEST_FILE}?type=file", data=file)
if not os.path.exists(f"{DATA_STORE}/{DIRECTORY_NAME}/{TEST_FILE}"):
    print(f"Directory {DATA_STORE}/{DIRECTORY_NAME}/{TEST_FILE} does not exist.")
    exit(1)

# File Download
log_cmd("File Download")
response = requests.get(BASE_URL + DIRECTORY_NAME + f"/{TEST_FILE}?type=file")
with open("output_file", 'wb') as file:
    file.write(response.content)
file_md5 = hashlib.md5(open("output_file", 'rb').read()).hexdigest()
if file_md5 != EXPECTED_FILE_MD5:
    print(f"File checksum mismatch: expected={EXPECTED_FILE_MD5} got={file_md5}")
    exit(1)
os.remove("output_file")

# File Head
log_cmd("File Head")
response = requests.head(BASE_URL + DIRECTORY_NAME + f"/{TEST_FILE}?type=file")
print(response.headers)

# Remove file
log_cmd("Remove file")
response = requests.delete(BASE_URL + DIRECTORY_NAME + f"/{TEST_FILE}?type=file")
if os.path.exists(f"{DATA_STORE}/{DIRECTORY_NAME}/{TEST_FILE}"):
    print(f"Directory {DATA_STORE}/{DIRECTORY_NAME}/{TEST_FILE} was not deleted correctly.")
    exit(1)

# Remove directory
log_cmd("Remove directory")
response = requests.delete(BASE_URL + DIRECTORY_NAME + "/?type=directory")
if os.path.exists(f"{DATA_STORE}/{DIRECTORY_NAME}/"):
    print(f"Directory {DATA_STORE}/{DIRECTORY_NAME}/ was not deleted correctly.")
    exit(1)
print(f"Directory {DATA_STORE}/{DIRECTORY_NAME}/ was deleted.")
