set GOOS=js
set GOARCH=wasm
go build -o assets\main.wasm cmd\wasm\main.go