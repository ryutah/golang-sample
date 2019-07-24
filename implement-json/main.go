package main

import (
	"encoding/json"
	"fmt"
	"reflect"
)

type foo struct {
	value string
	bar   bar
}

type bar struct {
	value string
}

var (
	_ json.Marshaler   = new(foo)
	_ json.Unmarshaler = new(foo)
)

func (f *foo) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"value": f.value,
		"bar": map[string]interface{}{
			"value": f.bar.value,
		},
	})
}

func (f *foo) UnmarshalJSON(b []byte) error {
	payload := make(map[string]interface{})
	if err := json.Unmarshal(b, &payload); err != nil {
		return err
	}
	if val, ok := payload["value"]; ok {
		strval, ok := val.(string)
		if !ok {
			return &json.UnmarshalTypeError{
				Value:  "string",
				Type:   reflect.TypeOf(val),
				Struct: "foo",
				Field:  "value",
			}
		}
		f.value = strval
	}

	if val, ok := payload["bar"]; ok {
		barval, ok := val.(map[string]interface{})
		if !ok {
			return &json.UnmarshalTypeError{
				Value:  "main.bar",
				Type:   reflect.TypeOf(val),
				Struct: "foo",
				Field:  "bar",
			}
		}
		if barvalvalue, ok := barval["value"]; ok {
			v, ok := barvalvalue.(string)
			if !ok {
				return &json.UnmarshalTypeError{
					Value:  "string",
					Type:   reflect.TypeOf(v),
					Struct: "bar",
					Field:  "value",
				}
			}
			f.bar.value = v
		}
	}
	return nil
}

func main() {
	f := &foo{
		value: "VALUE",
		bar: bar{
			value: "Bar",
		},
	}
	payload, err := json.MarshalIndent(f, "", "  ")
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s\n", payload)

	jsonValue := `{
		"value": "valval",
		"bar": {
			"value": "barbar"
		}
	}`
	var newFoo foo
	if err := json.Unmarshal([]byte(jsonValue), &newFoo); err != nil {
		panic(err)
	}
	fmt.Printf("%#v\n", newFoo)
}
