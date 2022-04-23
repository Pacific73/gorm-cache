package callback

import (
	"fmt"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func GetPrimaryKeysFromWhereClause(db *gorm.DB) []string {
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
