package callback

import (
	"context"
	"fmt"
	"reflect"

	"gorm.io/gorm"
)

func GetPrimaryKeysAfterCreate(db *gorm.DB) []string {
	primaryKeys := make([]string, 0)

	objects := make([]reflect.Value, 0)

	destValue := reflect.Indirect(reflect.ValueOf(db.Statement.Dest))
	switch destValue.Kind() {
	case reflect.Slice, reflect.Array:
		for i := 0; i < destValue.Len(); i++ {
			elem := destValue.Index(i)
			objects = append(objects, elem)
		}
	case reflect.Struct:
		objects = append(objects, destValue)
	}

	for _, elemValue := range objects {
		for _, field := range db.Statement.Schema.Fields {
			if field.PrimaryKey {
				primaryKey, isZero := field.ValueOf(context.Background(), elemValue)
				if isZero {
					continue
				}
				primaryKeys = append(primaryKeys, fmt.Sprintf("%v", primaryKey))
			}
		}
	}
	return primaryKeys
}
