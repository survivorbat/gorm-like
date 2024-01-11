package gormlike

import (
	"fmt"
	"strings"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const tagName = "gormlike"

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

	for index, cond := range exp.Exprs {
		switch cond := cond.(type) {
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

			condition := fmt.Sprintf("%s LIKE ?", cond.Column)

			if d.replaceCharacter != "" {
				value = strings.ReplaceAll(value, d.replaceCharacter, "%")
			}

			exp.Exprs[index] = db.Session(&gorm.Session{NewDB: true}).Where(condition, value).Statement.Clauses["WHERE"].Expression
		case clause.IN:
			columnName, columnOk := cond.Column.(string)
			if !columnOk {
				continue
			}

			// Get the `gormlike` value
			tagValue := db.Statement.Schema.FieldsByDBName[columnName].Tag.Get(tagName)

			// If the user has explicitly set this to false, ignore this field
			if tagValue == "false" {
				continue
			}

			// If tags are required and the tag is not true, ignore this field
			if d.conditionalTag && tagValue != "true" {
				continue
			}

			var likeCounter int
			var useOr bool

			query := db.Session(&gorm.Session{NewDB: true})

			for _, value := range cond.Values {
				value, ok := value.(string)
				if !ok {
					continue
				}

				condition := fmt.Sprintf("%s = ?", cond.Column)

				// If there are no % AND there aren't ony replaceable characters, just skip it because it's a normal query
				if strings.Contains(value, "%") || (d.replaceCharacter != "" && strings.Contains(value, d.replaceCharacter)) {
					condition = fmt.Sprintf("%s LIKE ?", cond.Column)

					if d.replaceCharacter != "" {
						value = strings.ReplaceAll(value, d.replaceCharacter, "%")
					}

					likeCounter++
				}

				if useOr {
					query = query.Or(condition, value)

					continue
				}

				query = query.Where(condition, value)
				useOr = true
			}

			// Don't alter the query if it isn't necessary
			if likeCounter == 0 {
				continue
			}

			exp.Exprs[index] = db.Session(&gorm.Session{NewDB: true}).Where(query).Statement.Clauses["WHERE"].Expression
		}
	}
}
