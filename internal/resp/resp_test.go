package resp

import (
	"bytes"
	"testing"

	"github.com/dylanmccormick/go-redis/internal/util"
)

func TestParseInt(t *testing.T) {
	inputBytes := []byte(":1\r\n")

	split := bytes.Split(inputBytes, util.SeparatorBytes)

	data, elements, _ := ParseRESP(split)

	expectedData := 1
	expectedExtraCRLFs := 0

	if data != expectedData {
		t.Errorf("Got data: %v Expected Data %v\n", data, expectedData)
	}

	if elements != expectedExtraCRLFs {
		t.Errorf("Got elements: %v Expected Elements: %v\n", elements, expectedExtraCRLFs)
	}

}

func TestParseString(t *testing.T) {
	inputBytes := []byte("+hello\r\n")
	split := bytes.Split(inputBytes, util.SeparatorBytes)

	data, elements, _ := ParseRESP(split)

	expectedData := "hello"
	expectedExtraCRLFs := 0

	if data != expectedData {
		t.Errorf("Got data: %v Expected Data %v\n", data, expectedData)
	}

	if elements != expectedExtraCRLFs {
		t.Errorf("Got elements: %v Expected Elements: %v\n", elements, expectedExtraCRLFs)
	}
}

func TestParseBulkString(t *testing.T) {
	inputBytes := []byte("$5\r\nhello\r\n")
	split := bytes.Split(inputBytes, util.SeparatorBytes)

	data, elements, _ := ParseRESP(split)

	expectedData := "hello"
	expectedExtraCRLFs := 1

	if data != expectedData {
		t.Errorf("Got data: %v Expected Data %v\n", data, expectedData)
	}

	if elements != expectedExtraCRLFs {
		t.Errorf("Got elements: %v Expected Elements: %v\n", elements, expectedExtraCRLFs)
	}
}

func TestParseNullBulk(t *testing.T) {
	inputBytes := []byte("$-1\r\n")
	split := bytes.Split(inputBytes, util.SeparatorBytes)

	data, elements, _ := ParseRESP(split)

	expectedData := ""
	expectedExtraCRLFs := 0

	if data != expectedData {
		t.Errorf("Got data: %v Expected Data %v\n", data, expectedData)
	}

	if elements != expectedExtraCRLFs {
		t.Errorf("Got elements: %v Expected Elements: %v\n", elements, expectedExtraCRLFs)
	}

}

func TestParseArray(t *testing.T) {
	inputBytes := []byte("*2\r\n:1\r\n+hello\r\n")
	split := bytes.Split(inputBytes, util.SeparatorBytes)

	data, elements, _ := ParseRESP(split)

	expectedData := []any{1, "hello"}
	expectedExtraCRLFs := 2

	for i, val := range data.([]any) {
		if val != expectedData[i] {
			t.Errorf("Index: %d, Got data: %v Expected Data %v\n", i, data, expectedData)
		}
	}

	if elements != expectedExtraCRLFs {
		t.Errorf("Got elements: %v Expected Elements: %v\n", elements, expectedExtraCRLFs)
	}
}

func TestParseArrayNull(t *testing.T) {
	inputBytes := []byte("*-1\r\n")
	split := bytes.Split(inputBytes, util.SeparatorBytes)

	data, elements, _ := ParseRESP(split)

	expectedData := []any{nil}
	expectedExtraCRLFs := 0

	for i, val := range data.([]any) {
		if val != expectedData[i] {
			t.Errorf("Index: %d, Got data: %v Expected Data %v\n", i, data, expectedData)
		}
	}

	if elements != expectedExtraCRLFs {
		t.Errorf("Got elements: %v Expected Elements: %v\n", elements, expectedExtraCRLFs)
	}
}

func TestBadArray(t *testing.T) {
	inputBytes := []byte("*5\r\n:1\r\n+hello\r\n")
	split := bytes.Split(inputBytes, util.SeparatorBytes)

	_, _, err := ParseRESP(split)

	if err == nil {
		t.Errorf("Expected an error. Got nil")
	}

}

func TestBadString(t *testing.T) {
	inputBytes := []byte("$3\r\nhello\r\n")
	split := bytes.Split(inputBytes, util.SeparatorBytes)

	_, _, err := ParseRESP(split)

	if err == nil {
		t.Errorf("Expected an error. Got nil")
	}

}

func TestSerializeData(t *testing.T) {
	inputData := []any{1, "hello world", []any{1, "nested", nil}, "done"}

	response, _ := Serialize(inputData)

	expectedString := "*4\r\n:1\r\n$11\r\nhello world\r\n*3\r\n:1\r\n$6\r\nnested\r\n$-1\r\n$4\r\ndone\r\n"

	if response != expectedString {
		t.Errorf("Response does not match expected.\n %#v | %#v", response, expectedString)
	}

}

func TestSerializeDataBad(t *testing.T) {
	type badType struct{ data int }
	inputData := []any{1, badType{3}, []any{1, "nested", nil}, "done"}

	_, err := Serialize(inputData)
	if err == nil {
		t.Errorf("Expected an error for bad type. Got nil")
	}

}
