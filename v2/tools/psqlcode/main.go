package main

import (
	"fmt"
	"os"
)

/*
- name: user
  table: test.users
  primaryKey:
    - id
  columns:
    id: int64
    name: string
    gender: xtype.Gender
    dob: string
    place: xtype.Place
    accounts: "[]string"
  jsonColumns:
    - place
- name: class_user
  table: test.class_users
  primaryKey:
    - class_id
    - user_id
  columns:
    class_id: int64
    user_id: int64
    created_time: time.Time
    score: float64
*/

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("Usage: %s {modelFilename}", os.Args[0])
		return
	}
	Generate(os.Args[1])
}
