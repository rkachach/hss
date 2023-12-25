set -e

log_cmd(){
    echo ""
    echo "--> $1"
}

DATA_STORE="/tmp/data-store"
DIRECTORY_NAME="mydir"

log_cmd "removing any existing object store: $DATA_STORE"
[ -d  "$DATA_STORE" ] && rm -r "$DATA_STORE"

curl -X PUT http://localhost:9000/$DIRECTORY_NAME/\?type\=directory
[ -d "$DATA_STORE/$DIRECTORY_NAME/" ] || { echo "Directory $DATA_STORE/$DIRECTORY_NAME/ does not exist."; exit 1; }

log_cmd "Get directory info"
curl -s -X GET http://localhost:9000/$DIRECTORY_NAME/\?type\=directory | jq

log_cmd "List directory files"
curl -s -X GET http://localhost:9000/$DIRECTORY_NAME/\?type\=directory\&operation\=list | jq

log_cmd "File Upload"
curl -X PUT -F "file=@./file.txt" http://localhost:9000/$DIRECTORY_NAME/file.txt\?type\=file
[ -e "$DATA_STORE/$DIRECTORY_NAME/file.txt" ] || { echo "Directory $DATA_STORE/$DIRECTORY_NAME/file.txt does not exist."; exit 1; }
