package gormlike

import (
	"gorm.io/gorm"
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

// TaggedOnly makes it so that only fields with the tag `gormlike` can be turned into LIKE queries,
// useful if you don't want every field to be LIKE-able.
func TaggedOnly() Option {
	return func(like *gormLike) {
		like.conditionalTag = true
	}
}

// SettingOnly makes it so that only queries with the setting 'gormlike' set to true can be turned into LIKE queries.
// This can be configured using db.Set("gormlike", true) on the query.
func SettingOnly() Option {
	return func(like *gormLike) {
		like.conditionalSetting = true
	}
}

// New creates a new instance of the plugin that can be registered in gorm. Without any settings, all queries will be
// LIKE-d.
//
//nolint:ireturn // Acceptable
func New(opts ...Option) gorm.Plugin {
	plugin := &gormLike{}

	for _, opt := range opts {
		opt(plugin)
	}

	return plugin
}

type gormLike struct {
	replaceCharacter   string
	conditionalTag     bool
	conditionalSetting bool
}

func (d *gormLike) Name() string {
	return "gormlike"
}

func (d *gormLike) Initialize(db *gorm.DB) error {
	return db.Callback().Query().Before("gorm:query").Register("gormlike:query", d.queryCallback)
}
