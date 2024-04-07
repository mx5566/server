.\protoc.exe --plugin=protoc-gen-go=./protoc-gen-go.exe --proto_path=. --go_out=. *.proto

.\protoc.exe --plugin=protoc-gen-go=./protoc-gen-go.exe --proto_path=..\..\base\rpc3 --go_out=..\..\base\rpc3  ..\..\base\rpc3\*.proto

pause