package naming_test

import (
	"testing"

	"code.olapie.com/sugar/v2/naming"
	"code.olapie.com/sugar/v2/testutil"
)

func TestSnakeToCamel(t *testing.T) {
	tests := []struct {
		Input  string
		Output string
	}{
		{
			Input:  "hello",
			Output: "hello",
		},
		{
			Input:  "hello_world",
			Output: "helloWorld",
		},
		{
			Input:  "hello_world_",
			Output: "helloWorld",
		},
		{
			Input:  "hello_world_id",
			Output: "helloWorldID",
		},
	}

	for _, tc := range tests {
		testutil.Equal(t, tc.Output, naming.ToCamel(tc.Input))
	}
}

func TestSnakeToClass(t *testing.T) {
	tests := []struct {
		Input  string
		Output string
	}{
		{
			Input:  "hello",
			Output: "Hello",
		},
		{
			Input:  "hello_world",
			Output: "HelloWorld",
		},
		{
			Input:  "hello_world_",
			Output: "HelloWorld",
		},
		{
			Input:  "hello_world_id",
			Output: "HelloWorldID",
		},
	}

	for _, tc := range tests {
		testutil.Equal(t, tc.Output, naming.ToClassName(tc.Input))
	}
}

func TestCamelToSnake(t *testing.T) {
	tests := []struct {
		Input  string
		Output string
	}{
		{
			Input:  "hello",
			Output: "hello",
		},
		{
			Input:  "helloWorld",
			Output: "hello_world",
		},
		{
			Input:  "helloWorldID",
			Output: "hello_world_id",
		},
	}

	for _, tc := range tests {
		testutil.Equal(t, tc.Output, naming.ToSnake(tc.Input))
	}
}
