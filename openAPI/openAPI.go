package openapi

import (
	"github.com/pb33f/libopenapi"
	v3high "github.com/pb33f/libopenapi/datamodel/high/v3"
)

// OpenAPI contain all field that need to be used in openapi package
// Only implement OpenAPI v3
type OpenAPI struct {
	document   libopenapi.Document
	docModelV3 *libopenapi.DocumentModel[v3high.Document]
}

// SimplifyAPI only list basic information about api
type SimplifyAPI struct {
	Description string           `json:"description" yaml:"description"`
	Summary     string           `json:"summary" yaml:"summary"`
	URL         string           `json:"url" yaml:"url"`
	Methods     []SimplifyMethod `json:"methods" yaml:"methods"`
}

// SimplifyMethod gathering simplify data from operation form openAPI
type SimplifyMethod struct {
	Method      string `json:"method" yaml:"method"`
	Description string `json:"description" yaml:"description"`
	Summary     string `json:"summary" yaml:"summary"`
}

// ListAllAPIFromDocument list all path and it's methods
func (o *OpenAPI) ListAllAPIFromDocument() []SimplifyAPI {
	var simplifyAPIs []SimplifyAPI

	for pathPairs := o.docModelV3.Model.Paths.PathItems.First(); pathPairs != nil; pathPairs = pathPairs.Next() {
		pathName := pathPairs.Key()
		pathItem := pathPairs.Value()
		// fmt.Printf("Path %s has %d operations\n", pathName, pathItem.GetOperations().Len())

		var simplifyMethods []SimplifyMethod

		for option := pathItem.GetOperations().First(); option != nil; option = option.Next() {
			method := SimplifyMethod{
				Method:      option.Key(), // method name like "get"
				Description: option.Value().Description,
				Summary:     option.Value().Summary,
			}
			simplifyMethods = append(simplifyMethods, method)
		}

		simplifyAPI := SimplifyAPI{
			Description: pathItem.Description,
			Summary:     pathItem.Summary,
			URL:         pathName,
			Methods:     simplifyMethods,
		}

		simplifyAPIs = append(simplifyAPIs, simplifyAPI)

	}

	return simplifyAPIs
}
