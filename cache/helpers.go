package cache

import (
	"context"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

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
	dbName := ""
	if db.Statement.Schema == nil {
		return nil
	}
	for _, field := range db.Statement.Schema.Fields {
		if field.PrimaryKey {
			dbName = field.DBName
			break
		}
	}
	if len(dbName) == 0 {
		return nil
	}
	for _, expr := range where.Exprs {
		eqExpr, ok := expr.(clause.Eq)
		if ok {
			if getColNameFromColumn(eqExpr.Column) == dbName {
				primaryKeys = append(primaryKeys, fmt.Sprintf("%v", eqExpr.Value))
			}
			continue
		}
		inExpr, ok := expr.(clause.IN)
		if ok {
			if getColNameFromColumn(inExpr.Column) == dbName {
				for _, val := range inExpr.Values {
					primaryKeys = append(primaryKeys, fmt.Sprintf("%v", val))
				}
			}
		}
		exprStruct, ok := expr.(clause.Expr)
		if ok {
			ttype := getExprType(exprStruct)
			//fmt.Printf("expr: %+v, ttype: %s\n", exprStruct, ttype)
			if ttype == "in" || ttype == "eq" {
				fieldName := getColNameFromExpr(exprStruct, ttype)
				if fieldName == dbName {
					pKeys := getPrimaryKeysFromExpr(exprStruct, ttype)
					primaryKeys = append(primaryKeys, pKeys...)
				}
			}
		}
	}
	return uniqueStringSlice(primaryKeys)
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
	dbName := ""
	for _, field := range db.Statement.Schema.Fields {
		if field.PrimaryKey {
			dbName = field.DBName
		}
	}
	if len(dbName) == 0 {
		return true // return true to skip cache
	}
	for _, expr := range where.Exprs {
		eqExpr, ok := expr.(clause.Eq)
		if ok {
			if getColNameFromColumn(eqExpr.Column) != dbName {
				return true
			}
			continue
		}
		inExpr, ok := expr.(clause.IN)
		if ok {
			if getColNameFromColumn(inExpr.Column) != dbName {
				return true
			}
			continue
		}
		exprStruct, ok := expr.(clause.Expr)
		if ok {
			ttype := getExprType(exprStruct)
			if ttype == "in" || ttype == "eq" {
				fieldName := getColNameFromExpr(exprStruct, ttype)
				if fieldName != dbName {
					return true
				}
				continue
			}
			return true
		}
		fmt.Printf("expr: %+v\n", expr)
		return true
	}
	return false
}

func getExprType(expr clause.Expr) string {
	// delete spaces
	sql := strings.Replace(strings.ToLower(expr.SQL), " ", "", -1)

	// see if sql has more than one clause
	hasConnector := strings.Contains(sql, "and") || strings.Contains(sql, "or")

	if strings.Contains(sql, "=") && !hasConnector {
		// possibly "id=?" or "id=123"
		fields := strings.Split(sql, "=")
		if len(fields) == 2 {
			_, isNumberErr := strconv.ParseInt(fields[1], 10, 64)
			if fields[1] == "?" || isNumberErr == nil {
				return "eq"
			}
		}
	} else if strings.Contains(sql, "in") && !hasConnector {
		// possibly "idIN(?)"
		fields := strings.Split(sql, "in")
		if len(fields) == 2 {
			if len(fields[1]) > 1 && fields[1][0] == '(' && fields[1][len(fields[1])-1] == ')' {
				return "in"
			}
		}
	}
	return "other"
}

func getColNameFromExpr(expr clause.Expr, ttype string) string {
	sql := strings.Replace(strings.ToLower(expr.SQL), " ", "", -1)
	if ttype == "in" {
		fields := strings.Split(sql, "in")
		return fields[0]
	} else if ttype == "eq" {
		fields := strings.Split(sql, "=")
		return fields[0]
	}
	return ""
}

func getPrimaryKeysFromExpr(expr clause.Expr, ttype string) []string {
	sql := strings.Replace(strings.ToLower(expr.SQL), " ", "", -1)

	primaryKeys := make([]string, 0)

	if ttype == "in" {
		fields := strings.Split(sql, "in")
		if len(fields) == 2 {
			if fields[1][0] == '(' && fields[1][len(fields[1])-1] == ')' {
				idStr := fields[1][1 : len(fields[1])-1]
				ids := strings.Split(idStr, ",")
				for _, id := range ids {
					if id == "?" {
						for _, vvar := range expr.Vars {
							keys := extractStringsFromVar(vvar)
							primaryKeys = append(primaryKeys, keys...)
						}
						break
					}
					number, err := strconv.ParseInt(id, 10, 64)
					if err == nil {
						primaryKeys = append(primaryKeys, strconv.FormatInt(number, 10))
					}
				}
			} else if fields[1] == "(?)" {
				for _, val := range expr.Vars {
					primaryKeys = append(primaryKeys, fmt.Sprintf("%v", val))
				}
			}
		}
	} else if ttype == "eq" {
		fields := strings.Split(sql, "=")
		if len(fields) == 2 {
			_, err := strconv.ParseInt(fields[1], 10, 64)
			if err == nil {
				primaryKeys = append(primaryKeys, fields[1])
			} else if fields[1] == "?" {
				for _, val := range expr.Vars {
					primaryKeys = append(primaryKeys, fmt.Sprintf("%v", val))
				}
			}
		}
	}
	return primaryKeys
}

func getObjectsAfterLoad(db *gorm.DB) (primaryKeys []string, objects []interface{}) {
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

	var valueOf func(context.Context, reflect.Value) (value interface{}, zero bool) = nil
	if db.Statement.Schema != nil {
		for _, field := range db.Statement.Schema.Fields {
			if field.PrimaryKey {
				valueOf = field.ValueOf
				break
			}
		}
	}

	objects = make([]interface{}, 0, len(values))
	for _, elemValue := range values {
		if valueOf != nil {
			primaryKey, isZero := valueOf(context.Background(), elemValue)
			if isZero {
				continue
			}
			primaryKeys = append(primaryKeys, fmt.Sprintf("%v", primaryKey))
		}
		objects = append(objects, elemValue.Interface())
	}
	return primaryKeys, objects
}

func uniqueStringSlice(slice []string) []string {
	retSlice := make([]string, 0)
	mmap := make(map[string]struct{})
	for _, str := range slice {
		_, ok := mmap[str]
		if !ok {
			mmap[str] = struct{}{}
			retSlice = append(retSlice, str)
		}
	}
	return retSlice
}

func extractStringsFromVar(v interface{}) []string {
	noPtrValue := reflect.Indirect(reflect.ValueOf(v))
	switch noPtrValue.Kind() {
	case reflect.Slice, reflect.Array:
		ans := make([]string, 0)
		for i := 0; i < noPtrValue.Len(); i++ {
			obj := reflect.Indirect(noPtrValue.Index(i))
			ans = append(ans, fmt.Sprintf("%v", obj))
		}
		return ans
	case reflect.String:
		return []string{fmt.Sprintf("%s", noPtrValue.Interface())}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8,
		reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return []string{fmt.Sprintf("%d", noPtrValue.Interface())}
	}
	return nil
}
