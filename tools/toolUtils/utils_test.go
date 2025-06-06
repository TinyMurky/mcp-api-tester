package toolutils

import (
	"testing"
)

func Test_createJSONSchemaFromToolHandlerFunc(t *testing.T) {
	type TestStruct struct {
		Name string `json:"name" jsonschema:"required,description=This is description,enum=add,enum=subtract"`
		// Tags map[string]interface{} `json:"tags,omitempty" jsonschema_extras:"a=b,foo=bar,foo=bar1"`
	}

	testFunc := func(_ any, _ TestStruct) {}

	testSchema := createJSONSchemaFromToolHandlerFunc(testFunc)

	// b, _ := json.Marshal(testSchema.Properties.Oldest())
	// fmt.Println(string(b))

	if testSchema.Properties.Len() != 1 {
		t.Errorf("Len of TestStruct should be %d, not %d", 1, testSchema.Properties.Len())
	}

	pair := testSchema.Properties.Oldest()

	if pair.Key != "name" {
		t.Errorf("pair `Key` should be %q (Because of json:\"name\"), not %s", "name", pair.Key)
	}
}

// func Test_Stuff(t *testing.T) {
// 	type TestStruct struct {
// 		Name string `json:"name" jsonschema:"required,description=This is description,enum=add,enum=subtract"`
// 	}

// 	aaa := reflect.ValueOf(&TestStruct{
// 		Name: "hi",
// 	}).Type()
// 	bbb := reflect.New(aaa)
// 	ccc := bbb.Interface() // 會把 (type, value) 的 interface variable 改成 (interface, value)

// 	fmt.Printf("aaaa\n%v\n", aaa.String())
// 	fmt.Printf("bbb\n%v\n", bbb.String())
// 	t.Error("")
// }
