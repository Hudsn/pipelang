package main

import (
	"encoding/json"
	"fmt"
	"log"
)

func main() {

	subData := DummySub{
		SubString: "987654321",
		SubInt:    987654321,
		List: [][]string{
			{
				"asdfjdskjfkdf",
				"fdasfkdsajf",
				"2439892834",
			},
		},
	}

	data := Dummy{
		MyString: "abcdefg",
		MyInt:    12345,
		MySub:    subData,
	}

	jsonb, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		log.Fatal(err)
	}

	var read any

	json.Unmarshal(jsonb, &read)

	typeCheck(read)

	jsonb, _ = json.Marshal(uint64(5))
	json.Unmarshal(jsonb, &read)
	fmt.Printf("%T \n", read)

	typeCheck(read)
}

func typeCheck(t any) {
	switch t.(type) {
	case map[string]interface{}:
		fmt.Println("MAP")
	case interface{}:
		fmt.Println("Interface{}")
	}
}

type Dummy struct {
	MyString string   `json:"my_string"`
	MyInt    int      `json:"my_int"`
	MySub    DummySub `json:"sub"`
}

type DummySub struct {
	SubString string     `json:"sub_string"`
	SubInt    int        `json:"sub_int"`
	List      [][]string `json:"list_list"`
}

type Type string

const (
	TSTRING  Type = "STRING"
	TNUM     Type = "NUM"
	TARRAY   Type = "ARRAY"
	TBOOL    Type = "BOOL"
	TOBJECT  Type = "OBJECT"
	TINT     Type = "INT"
	TFLOAT32 Type = "F32"
	TFLOAT64 Type = "F64"
)
