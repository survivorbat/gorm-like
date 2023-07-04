package gormlike

import (
	"fmt"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"strings"
)

const tagName = "gormlike"

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
			if d.conditionalTag {
				value := db.Statement.Schema.FieldsByDBName[cond.Column.(string)].Tag.Get(tagName)

				// Ignore if there's no valid tag settingValue
				if value != "true" {
					continue
				}
			}

			value, ok := cond.Value.(string)
			if !ok {
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
			if d.conditionalTag {
				value := db.Statement.Schema.FieldsByDBName[cond.Column.(string)].Tag.Get(tagName)

				// Ignore if there's no valid tag settingValue
				if value != "true" {
					continue
				}
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

			// TODO: Determine whether this is efficient
			exp.Exprs[index] = db.Session(&gorm.Session{NewDB: true}).Where(query).Statement.Clauses["WHERE"].Expression
		}
	}
}
