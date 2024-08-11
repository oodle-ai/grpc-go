module github.com/oodle-ai/grpc-go/security/advancedtls/examples

go 1.21

require (
	github.com/oodle-ai/grpc-go v1.64.0
	github.com/oodle-ai/grpc-go/examples v0.0.0-20240606220939-dfcabe08c639
	github.com/oodle-ai/grpc-go/security/advancedtls v0.0.0-20240606220939-dfcabe08c639
)

require (
	golang.org/x/crypto v0.24.0 // indirect
	golang.org/x/net v0.26.0 // indirect
	golang.org/x/sys v0.21.0 // indirect
	golang.org/x/text v0.16.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20240604185151-ef581f913117 // indirect
	google.golang.org/protobuf v1.34.1 // indirect
)

replace github.com/oodle-ai/grpc-go => ../../..

replace github.com/oodle-ai/grpc-go/examples => ../../../examples

replace github.com/oodle-ai/grpc-go/security/advancedtls => ../
