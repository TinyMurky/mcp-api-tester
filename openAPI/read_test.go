package openapi

import (
	"encoding/json"
	"testing"
)

func Test_ReadFromPath(t *testing.T) {
	_, err := ReadFromPath("/home/tinymurky/Desktop/ISunFa.yaml")

	if err != nil {
		t.Errorf("%v", err)
	}

	simplifyAPIs := OpenAPIPointer.ListAllAPIFromDocument()

	_, err = json.Marshal(simplifyAPIs)

	if err != nil {
		t.Errorf("%v", err)
	}
}
