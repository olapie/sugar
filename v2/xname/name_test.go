package xname_test

import (
	"testing"

	"code.olapie.com/sugar/xname"
	"code.olapie.com/sugar/xtest"
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
		xtest.Equal(t, tc.Output, xname.ToCamel(tc.Input))
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
		xtest.Equal(t, tc.Output, xname.ToClassName(tc.Input))
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
		xtest.Equal(t, tc.Output, xname.ToSnake(tc.Input))
	}
}
