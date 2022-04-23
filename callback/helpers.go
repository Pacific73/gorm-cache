package callback

import (
	"context"
	"errors"
	"fmt"
	"reflect"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var primaryCacheHit = errors.New("primary cache hit")
var searchCacheHit = errors.New("search cache hit")

// getPrimaryKeysFromWhereClause try to find primary keys from Eq and IN exprs in WHERE clause,
// and get objects that are being operated
func getPrimaryKeysFromWhereClause(db *gorm.DB) []string {
	primaryKeys := make([]string, 0)

	cla, ok := db.Statement.Clauses["WHERE"]
	if !ok {
		return nil
	}
	where, ok := cla.Expression.(clause.Where)
	if !ok {
		return nil
	}
	for _, field := range db.Statement.Schema.Fields {
		if field.PrimaryKey {
			for _, expr := range where.Exprs {
				eqExpr, ok := expr.(clause.Eq)
				if ok {
					if getColNameFromColumn(eqExpr.Column) == field.DBName {
						primaryKeys = append(primaryKeys, fmt.Sprintf("%v", eqExpr.Value))
					}
					continue
				}
				inExpr, ok := expr.(clause.IN)
				if ok {
					if getColNameFromColumn(inExpr.Column) == field.DBName {
						for _, val := range inExpr.Values {
							primaryKeys = append(primaryKeys, fmt.Sprintf("%v", val))
						}
					}
				}
			}
		}
	}
	return primaryKeys
}

func getColNameFromColumn(col interface{}) string {
	switch v := col.(type) {
	case string:
		return v
	case clause.Column:
		return v.Name
	default:
		return ""
	}
}

func hasOtherClauseExceptPrimaryField(db *gorm.DB) bool {
	cla, ok := db.Statement.Clauses["WHERE"]
	if !ok {
		return false
	}
	where, ok := cla.Expression.(clause.Where)
	for _, field := range db.Statement.Schema.Fields {
		if field.PrimaryKey {
			for _, expr := range where.Exprs {
				eqExpr, ok := expr.(clause.Eq)
				if ok {
					if getColNameFromColumn(eqExpr.Column) != field.DBName {
						return true
					}
					continue
				}
				inExpr, ok := expr.(clause.IN)
				if ok {
					if getColNameFromColumn(inExpr.Column) != field.DBName {
						return true
					}
					continue
				}
				return true
			}
		}
	}
	return false
}

func ContainString(target string, slice []string) bool {
	for _, s := range slice {
		if target == s {
			return true
		}
	}
	return false
}

func GetObjectsAfterLoad(db *gorm.DB) (primaryKeys []string, objects []interface{}) {
	primaryKeys = make([]string, 0)
	values := make([]reflect.Value, 0)

	destValue := reflect.Indirect(reflect.ValueOf(db.Statement.Dest))
	switch destValue.Kind() {
	case reflect.Slice, reflect.Array:
		for i := 0; i < destValue.Len(); i++ {
			elem := destValue.Index(i)
			values = append(values, elem)
		}
	case reflect.Struct:
		values = append(values, destValue)
	}

	objects = make([]interface{}, 0, len(values))
	for _, elemValue := range values {
		for _, field := range db.Statement.Schema.Fields {
			if field.PrimaryKey {
				primaryKey, isZero := field.ValueOf(context.Background(), elemValue)
				if isZero {
					continue
				}
				primaryKeys = append(primaryKeys, fmt.Sprintf("%v", primaryKey))
			}
		}
		objects = append(objects, elemValue.Interface())
	}
	return primaryKeys, objects
}
