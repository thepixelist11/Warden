version="$1"

echo "Building Windows AMD64"
export GOOS=windows
export GOARCH=amd64
go build -o "./bin/mcstatus-$version-win-amd64-test.exe" ./src/main.go

echo "Building Windows ARM64"
export GOOS=windows
export GOARCH=arm64
go build -o "./bin/mcstatus-$version-win-arm64-test.exe" ./src/main.go

echo "Building Linux AMD64"
export GOOS=linux
export GOARCH=amd64
go build -o "./bin/mcstatus-$version-linux-amd64-test" ./src/main.go

echo "Building Linux ARM64"
export GOOS=linux
export GOARCH=arm64
go build -o "./bin/mcstatus-$version-linux-arm64-test" ./src/main.go
