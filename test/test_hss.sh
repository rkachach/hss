set -e

log_cmd(){
    echo ""
    echo "--> $1"
}

DATA_STORE="/tmp/data-store"
DIRECTORY_NAME="mydir"
TEST_FILE="file.txt"
EXPECTED_FILE_MD5="5dd39cab1c53c2c77cd352983f9641e1"

log_cmd "removing any existing object store: $DATA_STORE"
[ -d  "$DATA_STORE" ] && rm -r "$DATA_STORE"

curl -X PUT http://localhost:9000/$DIRECTORY_NAME/\?type\=directory
[ -d "$DATA_STORE/$DIRECTORY_NAME/" ] || { echo "Directory $DATA_STORE/$DIRECTORY_NAME/ does not exist."; exit 1; }

log_cmd "Get directory info"
curl -s -X GET http://localhost:9000/$DIRECTORY_NAME/\?type\=directory | jq

log_cmd "List directory files"
curl -s -X GET http://localhost:9000/$DIRECTORY_NAME/\?type\=directory\&operation\=list | jq

log_cmd "File Upload"
curl -X PUT -H "Content-Type: application/octet-stream" --data-binary "@./$TEST_FILE" http://localhost:9000/$DIRECTORY_NAME/$TEST_FILE\?type\=file
[ -e "$DATA_STORE/$DIRECTORY_NAME/$TEST_FILE" ] || { echo "Directory $DATA_STORE/$DIRECTORY_NAME/$TEST_FILE does not exist."; exit 1; }

exit 0

log_cmd "Remove directory"
curl -X DELETE http://localhost:9000/$DIRECTORY_NAME/\?type\=directory
[ ! -d "$DATA_STORE/$DIRECTORY_NAME/" ] || { echo "Directory $DATA_STORE/$DIRECTORY_NAME/ was not deleted correctly."; exit 1; }
