module github.com/oodle-ai/grpc-go/stats/opencensus

go 1.21

require (
	github.com/google/go-cmp v0.6.0
	go.opencensus.io v0.24.0
	github.com/oodle-ai/grpc-go v1.64.0
)

require (
	github.com/golang/groupcache v0.0.0-20210331224755-41bb18bfe9da // indirect
	golang.org/x/net v0.26.0 // indirect
	golang.org/x/sys v0.21.0 // indirect
	golang.org/x/text v0.16.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20240604185151-ef581f913117 // indirect
	google.golang.org/protobuf v1.34.1 // indirect
)

replace github.com/oodle-ai/grpc-go => ../..
