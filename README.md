# 🌌 Gorm Auto Like Plugin

[![Go package](https://github.com/survivorbat/gorm-like/actions/workflows/test.yaml/badge.svg)](https://github.com/survivorbat/gorm-like/actions/workflows/test.yaml)

I wanted to provide a map to a WHERE query and automatically turn it into a LIKE query if wildcards were present, this
plugin does just that.

```go
package main

import (
	gormlike "github.com/survivorbat/gorm-like"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	db, _ := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})

	// Will match anything with a b
	filters := map[string]any{
		"name": "%b%",
	}

	db.Use(gormlike.New())
	db.Where(filters)

	db, _ = gorm.Open(sqlite.Open("test.db"), &gorm.Config{})

	// Will match anything with a b
	filters := map[string]any{
		"name": "🍌b🍌",
	}

	db.Use(gormlike.New(gormlike.WithCharacter("🍌")))
	db.Where(filters)
}
```

Is automatically turned into a query that looks like this:

```sql
SELECT * FROM employees WHERE name LIKE "%b%";
```

## ⬇️ Installation

`go get github.com/survivorbat/gorm-likes`

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

Turn it into a tag to make it a per-field basis instead of all of them.
