#!/bin/bash
#

localstack start -d

echo "test\ntest\n\n\n" | aws configure --profile default
aws s3 mb s3://my-bucket --endpoint-url http://localhost:4566 --profile default
aws s3 mb s3://my-bucket --endpoint-url http://localhost:4566 --profile default

mkdir /tmp/s3-test -p
echo "testfiledataaaaaaaaaaaaaaa" > /tmp/s3-test/file-1
aws s3 cp /tmp/s3-test/file-1 s3://my-bucket --profile default --endpoint-url http://localhost:4566
