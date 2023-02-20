module examples

go 1.20

require (
	code.olapie.com/sugar/v2 v2.3.15
	github.com/gorilla/mux v1.8.0
)

require (
	github.com/golang/geo v0.0.0-20210211234256-740aa86cb551 // indirect
	github.com/google/uuid v1.3.0 // indirect
	golang.org/x/text v0.6.0 // indirect
	google.golang.org/protobuf v1.28.1 // indirect
)

replace (
	code.olapie.com/sugar/v2 => ../../
)