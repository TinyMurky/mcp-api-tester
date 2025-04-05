package openapi

import (
	"fmt"
	"os"
	"strings"

	"github.com/pb33f/libopenapi"
)

// ReadFromPath will read openAPI file from given path
func ReadFromPath(path string) (*OpenAPI, error) {
	openAPIFileBinary, err := os.ReadFile(path)

	if err != nil {
		return nil, fmt.Errorf("Error happened read open api file from path: %q, error: %w", path, err)
	}

	document, err := libopenapi.NewDocument(openAPIFileBinary)

	if err != nil {
		return nil, fmt.Errorf("Error happened when convert openAPI Binary to document from path: %q, error: %w", path, err)
	}

	docModel, errs := document.BuildV3Model()

	if len(errs) > 0 {
		errorMessages := make([]string, len(errs))

		for _, err := range errs {
			errorMessages = append(errorMessages, err.Error())
		}

		return nil, fmt.Errorf("Error happened when build openAPI, errors: %s", strings.Join(errorMessages, ", "))
	}

	// Check openApi/instance
	OpenAPIPointer = &OpenAPI{
		document:   document,
		docModelV3: docModel,
	}

	return OpenAPIPointer, nil
}
