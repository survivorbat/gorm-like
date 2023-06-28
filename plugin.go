package gormlike

import (
	"fmt"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"strings"
)

// Compile-time interface check
var _ gorm.Plugin = new(gormLike)

// Option can be given to the New() method to tweak its behaviour
type Option func(like *gormLike)

// WithCharacter allows you to specify a replacement character for the % in the LIKE queries
func WithCharacter(character string) Option {
	return func(like *gormLike) {
		like.replaceCharacter = character
	}
}

// New creates a new instance of the plugin that can be registered in gorm.
func New(opts ...Option) gorm.Plugin {
	plugin := &gormLike{}

	for _, opt := range opts {
		opt(plugin)
	}

	return plugin
}

type gormLike struct {
	replaceCharacter string
}

func (d *gormLike) Name() string {
	return "gormlike"
}

func (d *gormLike) Initialize(db *gorm.DB) error {
	return db.Callback().Query().Before("gorm:query").Register("gormlike:query", d.queryCallback)
}

func (d *gormLike) queryCallback(db *gorm.DB) {
	exp, ok := db.Statement.Clauses["WHERE"].Expression.(clause.Where)
	if !ok {
		return
	}

	for index, cond := range exp.Exprs {
		switch cond := cond.(type) {
		case clause.Eq:
			switch value := cond.Value.(type) {
			case string:
				// If there are no % AND there aren't ony replaceable characters, just skip it because it's a normal query
				if !strings.Contains(value, "%") && !(d.replaceCharacter != "" && strings.Contains(value, d.replaceCharacter)) {
					continue
				}

				condition := fmt.Sprintf("%s LIKE ?", cond.Column)

				if d.replaceCharacter != "" {
					value = strings.ReplaceAll(value, d.replaceCharacter, "%")
				}

				exp.Exprs[index] = db.Session(&gorm.Session{NewDB: true}).Where(condition, value).Statement.Clauses["WHERE"].Expression
			}
		case clause.IN:
			var likeCounter int
			var useOr bool

			query := db.Session(&gorm.Session{NewDB: true})

			for _, value := range cond.Values {
				switch value := value.(type) {
				case string:
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
			}

			// Don't alter the query if it isn't necessary
			if likeCounter == 0 {
				continue
			}

			exp.Exprs[index] = query.Statement.Clauses["WHERE"].Expression
		}
	}

	return
}
