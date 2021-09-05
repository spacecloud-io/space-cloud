rm ./space-cli.zip
rm ./space-cli
rm ./space-cli.exe

SPACE_CLI_VERSION="0.21.5"

GOOS=linux CGO_ENABLED=0 go build -a -ldflags '-s -w -extldflags "-static"' .
zip space-cli.zip ./space-cli
rm ./space-cli

gsutil cp ./space-cli.zip gs://space-cloud/linux/space-cli.zip
gsutil cp ./space-cli.zip gs://space-cloud/linux/space-cli-v$SPACE_CLI_VERSION.zip
rm ./space-cli.zip

GOOS=windows CGO_ENABLED=0 go build -a -ldflags '-s -w -extldflags "-static"' .
zip space-cli.zip ./space-cli.exe
rm ./space-cli.exe

gsutil cp ./space-cli.zip gs://space-cloud/windows/space-cli.zip
gsutil cp ./space-cli.zip gs://space-cloud/windows/space-cli-v$SPACE_CLI_VERSION.zip
rm ./space-cli.zip

GOOS=darwin CGO_ENABLED=0 go build -a -ldflags '-s -w -extldflags "-static"' .
zip space-cli.zip ./space-cli
rm ./space-cli

gsutil cp ./space-cli.zip gs://space-cloud/darwin/space-cli.zip
gsutil cp ./space-cli.zip gs://space-cloud/darwin/space-cli-v$SPACE_CLI_VERSION.zip
rm ./space-cli.zip
