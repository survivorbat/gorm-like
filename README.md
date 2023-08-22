# 👍 Gorm Auto Like Plugin

[![Go package](https://github.com/survivorbat/gorm-like/actions/workflows/test.yaml/badge.svg)](https://github.com/survivorbat/gorm-like/actions/workflows/test.yaml)

I wanted to provide a map to a WHERE query and automatically turn it into a LIKE query if wildcards were present, this
plugin does just that.

By default, all queries are turned into like-queries if either a % or a given character is found, if you don't want this,
you have 2 options:

- `TaggedOnly()`: Will only change queries on fields that have the `gormlike:"true"` tag
- `SettingOnly()`: Will only change queries on `*gorm.DB` objects that have `.Set("gormlike", true)` set.

If you want a particular query to not be like-able, use `.Set("gormlike", false)`. This works
regardless of configuration.

## 💡 Related Libraries

- [deepgorm](https://github.com/survivorbat/gorm-deep-filtering) turns nested maps in WHERE-calls into subqueries
- [gormlike](https://github.com/survivorbat/gorm-like) turns WHERE-calls into LIkE queries if certain tokens were found
- [gormcase](https://github.com/survivorbat/gorm-case) adds case insensitivity to WHERE queries
- [gormtestutil](https://github.com/ing-bank/gormtestutil) provides easy utility methods for unit-testing with gorm

## ⬇️ Installation

`go get github.com/survivorbat/gorm-like`

## 📋 Usage

```go
package main

import (
    "github.com/survivorbat/gorm-like"
)

func main() {
	db, _ := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	db.Use(gormlike.New())
}

```

## 🔭 Plans

Not much here.
