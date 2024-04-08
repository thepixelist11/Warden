$version = $args[0]

Set-Variable GOOS=windows
Set-Variable GOARCH=amd64
go build -o ".\bin\mcstatus-$($version)-win-amd64.exe"

Set-Variable GOOS=windows
Set-Variable GOARCH=arm64
go build -o ".\bin\mcstatus-$($version)-win-arm64.exe"

Set-Variable GOOS=linux
Set-Variable GOARCH=amd64
go build -o ".\bin\mcstatus-$($version)-linux-amd64"

Set-Variable GOOS=linux
Set-Variable GOARCH=arm64
go build -o ".\bin\mcstatus-$($version)-linux-arm64"