protoc --proto_path ../../../ -I=./proto --go_out=plugins=grpc:./proto proto/keymapper.proto
mv proto/github.com/brotherlogic/keymapper/proto/* ./proto
