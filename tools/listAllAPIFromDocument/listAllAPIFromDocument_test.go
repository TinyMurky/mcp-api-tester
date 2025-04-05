package listallapifromdocument

import (
	"context"
	openapi "mcp-api-tester/openAPI"
	"testing"
)

func Test_ReadFromPath(t *testing.T) {
	_, err := openapi.ReadFromPath("/home/tinymurky/Desktop/ISunFa.yaml")

	if err != nil {
		t.Errorf("%v", err)
	}

	param := Param{}
	_, err = listAllAPIFromDocument(context.Background(), param)

	if err != nil {
		t.Errorf("%v", err)
	}
	// t.Errorf("%s", result)
}
