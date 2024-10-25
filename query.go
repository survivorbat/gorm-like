package gormlike

import (
	"fmt"
	"strings"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const tagName = "gormlike"

func (d *gormLike) replaceExpressions(db *gorm.DB, expressions []clause.Expression) []clause.Expression {
	for index, cond := range expressions {
		switch cond := cond.(type) {
		case clause.AndConditions:
			// Recursively go through the expressions of AndConditions
			cond.Exprs = d.replaceExpressions(db, cond.Exprs)
			expressions[index] = cond
		case clause.OrConditions:
			// Recursively go through the expressions of OrConditions
			cond.Exprs = d.replaceExpressions(db, cond.Exprs)
			expressions[index] = cond
		case clause.Eq:
			columnName, columnOk := cond.Column.(string)
			if !columnOk {
				continue
			}

			// Get the `gormlike` value
			var tagValue string
			dbField, ok := db.Statement.Schema.FieldsByDBName[columnName]
			if ok {
				tagValue = dbField.Tag.Get(tagName)
			}

			// If the user has explicitly set this to false, ignore this field
			if tagValue == "false" {
				continue
			}

			// If tags are required and the tag is not true, ignore this field
			if d.conditionalTag && tagValue != "true" {
				continue
			}

			value, columnOk := cond.Value.(string)
			if !columnOk {
				continue
			}

			// If there are no % AND there aren't ony replaceable characters, just skip it because it's a normal query
			if !strings.Contains(value, "%") && !(d.replaceCharacter != "" && strings.Contains(value, d.replaceCharacter)) {
				continue
			}

			condition := fmt.Sprintf("CAST(%s as varchar) LIKE ?", cond.Column)

			if d.replaceCharacter != "" {
				value = strings.ReplaceAll(value, d.replaceCharacter, "%")
			}

			expressions[index] = db.Session(&gorm.Session{NewDB: true}).Where(condition, value).Statement.Clauses["WHERE"].Expression
		case clause.IN:
			columnName, columnOk := cond.Column.(string)
			if !columnOk {
				continue
			}

			// Get the `gormlike` value
			var tagValue string
			dbField, ok := db.Statement.Schema.FieldsByDBName[columnName]
			if ok {
				tagValue = dbField.Tag.Get(tagName)
			}

			// If the user has explicitly set this to false, ignore this field
			if tagValue == "false" {
				continue
			}

			// If tags are required and the tag is not true, ignore this field
			if d.conditionalTag && tagValue != "true" {
				continue
			}

			var likeCounter int

			query := db.Session(&gorm.Session{NewDB: true})

			for _, value := range cond.Values {
				value, ok := value.(string)
				if !ok {
					continue
				}

				condition := fmt.Sprintf("%s = ?", cond.Column)

				// If there are no % AND there aren't ony replaceable characters, just skip it because it's a normal query
				if strings.Contains(value, "%") || (d.replaceCharacter != "" && strings.Contains(value, d.replaceCharacter)) {
					condition = fmt.Sprintf("CAST(%s as varchar) LIKE ?", cond.Column)

					if d.replaceCharacter != "" {
						value = strings.ReplaceAll(value, d.replaceCharacter, "%")
					}

					likeCounter++
				}

				query = query.Or(condition, value)
			}

			// Don't alter the query if it isn't necessary
			if likeCounter == 0 {
				continue
			}

			// This feels a bit like a dirty hack
			// but otherwise the generated query would not be correct in case of an AND condition between multiple OR conditions
			// e.g. without this -> x = .. OR x = .. AND y = .. OR y = .. (no brackets around the OR conditions mess up the query)
			// e.g. with this -> (x = .. OR x = ..) AND (y = .. OR y = ..)
			var newExpression clause.OrConditions
			newExpression.Exprs = query.Statement.Clauses["WHERE"].Expression.(clause.Where).Exprs

			expressions[index] = newExpression
		}
	}

	return expressions
}

//nolint:gocognit,cyclop // Acceptable
func (d *gormLike) queryCallback(db *gorm.DB) {
	// If we only want to like queries that are explicitly set to true, we back out early if anything's amiss
	settingValue, settingOk := db.Get(tagName)
	if d.conditionalSetting && !settingOk {
		return
	}

	if settingOk {
		if boolValue, _ := settingValue.(bool); !boolValue {
			return
		}
	}

	exp, settingOk := db.Statement.Clauses["WHERE"].Expression.(clause.Where)
	if !settingOk {
		return
	}

	exp.Exprs = d.replaceExpressions(db, exp.Exprs)
}
