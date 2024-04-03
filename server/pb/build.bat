.\protoc.exe --plugin=protoc-gen-go=./protoc-gen-go.exe --proto_path=. --go_out=. *.proto

.\protoc.exe --plugin=protoc-gen-go=./protoc-gen-go.exe --proto_path=..\..\rpc --go_out=..\..\rpc  ..\..\rpc\*.proto

pause