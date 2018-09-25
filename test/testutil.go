package test

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"reflect"
	"testing"
)

// GetTestFileAsString gets the file as a string... Test will fail if the file cannot be read.
func GetTestFileAsString(t *testing.T, filepath string) string {
	b, err := ioutil.ReadFile(filepath)
	if err != nil {
		t.Fatal(err)
	}
	return string(b)
}

// UnmarshalJSONFromFile unmarshals from a json file to interface v
func UnmarshalJSONFromFile(t *testing.T, filepath string, v interface{}) {
	file, err := os.OpenFile(filepath, os.O_RDONLY, 0644)
	AssertNotError(t, err)

	err = json.NewDecoder(file).Decode(&v)
	AssertNotError(t, err)
}

// AssertNotError asserts if an error equals nil or fails the test
func AssertNotError(t *testing.T, err error) {
	if err != nil {
		t.Fatalf("Got unexpected error: %s", err)
	}
}

// AssertIsError asserts if an error equals nil and and fails the test if its not
func AssertIsError(t *testing.T, err error) {
	if err == nil {
		t.Fatal("Expected error. Was nil.")
	}
}

// AssertEqual asserts if two object are same type and equal value
func AssertEqual(t *testing.T, actual interface{}, expected interface{}) {
	if actual != expected {
		t.Errorf("Expected %v %v to be equal to %v %v", reflect.TypeOf(actual).Name(), actual, reflect.TypeOf(expected).Name(), expected)
	}
}

// AssertJSONEqual asserts two json strings for equality by unmarshalling them to account for formatting.
func AssertJSONEqual(t *testing.T, json1 string, json2 string) {
	var o1 interface{}
	var o2 interface{}

	var err error
	err = json.Unmarshal([]byte(json1), &o1)
	if err != nil {
		t.Errorf("Error mashalling string 1 :: %s", err.Error())
	}
	err = json.Unmarshal([]byte(json2), &o2)
	if err != nil {
		t.Errorf("Error mashalling string 2 :: %s", err.Error())
	}

	if !reflect.DeepEqual(o1, o2) {
		t.Error("Json strings are not equal:")
		t.Errorf("Actual: %s", json1)
		t.Errorf("Expected: %s", json2)
	}
}
