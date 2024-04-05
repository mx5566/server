.\protoc.exe --plugin=protoc-gen-go=./protoc-gen-go.exe --proto_path=. --go_out=. *.proto

.\protoc.exe --plugin=protoc-gen-go=./protoc-gen-go.exe --proto_path=..\..\rpc3 --go_out=..\..\rpc3  ..\..\rpc3\*.proto

pause