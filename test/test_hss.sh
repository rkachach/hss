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

curl -X POST \
  -H "Content-Type: application/json" \
  -H "Metadata-Fields: custom-dir-field1,custom-dir-field2" \
  -H "custom-dir-field1: value1" \
  -H "custom-dir-field2: value2" \
  http://localhost:9000/$DIRECTORY_NAME/\?type\=directory

[ -d "$DATA_STORE/$DIRECTORY_NAME/" ] || { echo "Directory $DATA_STORE/$DIRECTORY_NAME/ does not exist."; exit 1; }

log_cmd "Get directory info"
curl -s -X GET http://localhost:9000/$DIRECTORY_NAME/\?type\=directory | jq

log_cmd "List directory files"
curl -s -X GET http://localhost:9000/$DIRECTORY_NAME/\?type\=directory\&operation\=list | jq

log_cmd "File Upload"
curl -X POST -H "Content-Type: application/octet-stream" --data-binary "@./$TEST_FILE" http://localhost:9000/$DIRECTORY_NAME/$TEST_FILE\?type\=file
[ -e "$DATA_STORE/$DIRECTORY_NAME/$TEST_FILE" ] || { echo "Directory $DATA_STORE/$DIRECTORY_NAME/$TEST_FILE does not exist."; exit 1; }

log_cmd "File Download"
curl -s -X GET http://localhost:9000/$DIRECTORY_NAME/$TEST_FILE\?type\=file -o output_file
FILE_MD5=$(md5sum output_file | awk '{ print $1 }')
[ "$FILE_MD5" = "$EXPECTED_FILE_MD5" ] || { echo "Fileect checksum mismatch: exp=$EXPECTED_FILE_MD5 got=$FILE_MD5"; exit 1; }
rm -f output_file

log_cmd "File Head"
curl -s --head http://localhost:9000/$DIRECTORY_NAME/$TEST_FILE\?type\=file

log_cmd "Remove file"
curl -s -X DELETE http://localhost:9000/$DIRECTORY_NAME/$TEST_FILE\?type\=file
[ ! -e "$DATA_STORE/$DIRECTORY_NAME/$TEST_FILE" ] || { echo "Directory $DATA_STORE/$DIRECTORY_NAME/$TEST_FILE was not deleted correctly."; exit 1; }

log_cmd "Remove directory"
curl -X DELETE http://localhost:9000/$DIRECTORY_NAME/\?type\=directory
[ ! -d "$DATA_STORE/$DIRECTORY_NAME/" ] || { echo "Directory $DATA_STORE/$DIRECTORY_NAME/ was not deleted correctly."; exit 1; }
[ ! -d "$DATA_STORE/$DIRECTORY_NAME/" ] || { echo "Directory $DATA_STORE/$DIRECTORY_NAME/ still exist."; exit 1; }
