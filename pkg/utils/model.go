package utils

import (
	"reflect"
)

// ModelsToIDs 转换 model 切片为 ID 切片
// models 必须是切片，且切片元素必须是 struct 或者 struct 指针，且 struct 必须有 int 类型的 ID 字段
func ModelsToIDs(models any) (ids []int) {
	modelsType := reflect.TypeOf(models)
	modelsValue := reflect.ValueOf(models)
	for modelsType.Kind() == reflect.Ptr {
		modelsType = modelsType.Elem()
		modelsValue = modelsValue.Elem()
	}
	if modelsType.Kind() != reflect.Slice {
		panic("models must be slice")
	}

	for i := 0; i < modelsValue.Len(); i++ {
		model := modelsValue.Index(i)
		modelType := model.Type()
		for modelType.Kind() == reflect.Ptr {
			modelType = modelType.Elem()
			model = model.Elem()
		}
		if modelType.Kind() != reflect.Struct {
			panic("models must be slice of struct")
		}

		idField := model.FieldByName("ID")
		if idField.Kind() != reflect.Int {
			panic("models must be slice of struct with int field named ID")
		}

		ids = append(ids, int(idField.Int()))
	}
	return
}
