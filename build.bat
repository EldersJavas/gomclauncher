SET CGO_ENABLED=0
SET GOOS=linux
SET GOARCH=amd64
go build -o gml-linux -trimpath -ldflags "-w -s"  

SET CGO_ENABLED=0
SET GOOS=darwin
SET GOARCH=amd64
go build -o gml-darwin -trimpath -ldflags "-w -s"  

SET CGO_ENABLED=0
SET GOOS=darwin
SET GOARCH=arm64
go build -o gml-arm-darwin -trimpath -ldflags "-w -s"  

SET CGO_ENABLED=0
SET GOOS=windows
SET GOARCH=amd64
go build -o gml-windows.exe -trimpath -ldflags "-w -s"  