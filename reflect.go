package config

import (
	"fmt"
	"reflect"
)

type ReflectField struct {
	Field     reflect.Value
	FieldType reflect.StructField
}

func getExportedFields(obj interface{}) ([]*ReflectField, error) {
	objValue, objType, err := getIndirect(obj)
	if err != nil {
		return nil, err
	}

	return getExportedFieldsStruct(objValue, objType)
}

func getExportedFieldsStruct(objValue reflect.Value, objType reflect.Type) ([]*ReflectField, error) {
	if objType.Kind() != reflect.Struct {
		return nil, fmt.Errorf(
			"invalid type embedded type in configuration struct",
		)
	}

	fields := []*ReflectField{}
	for i := 0; i < objType.NumField(); i++ {
		field := objValue.Field(i)
		fieldType := objType.Field(i)

		if fieldType.Anonymous {
			embeddedFields, err := getExportedFieldsStruct(field, fieldType.Type)
			if err != nil {
				return nil, err
			}

			fields = append(fields, embeddedFields...)
			continue
		}

		if !isExported(fieldType.Name) {
			continue
		}

		fields = append(fields, &ReflectField{
			Field:     field,
			FieldType: fieldType,
		})
	}

	return fields, nil
}

func getIndirect(obj interface{}) (reflect.Value, reflect.Type, error) {
	indirect := reflect.Indirect(reflect.ValueOf(obj))
	if !indirect.IsValid() {
		return reflect.Value{}, nil, fmt.Errorf("configuration target is not a pointer to struct")
	}

	indirectType := indirect.Type()
	if indirectType.Kind() != reflect.Struct {
		return reflect.Value{}, nil, fmt.Errorf("configuration target is not a pointer to struct")
	}

	return indirect, indirectType, nil
}
