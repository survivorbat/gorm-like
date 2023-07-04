# 👍 Gorm Auto Like Plugin

[![Go package](https://github.com/survivorbat/gorm-like/actions/workflows/test.yaml/badge.svg)](https://github.com/survivorbat/gorm-like/actions/workflows/test.yaml)

I wanted to provide a map to a WHERE query and automatically turn it into a LIKE query if wildcards were present, this
plugin does just that. You can either do it for all queries or only for specific fields using the tag `gormlike:"true"`.

```go
package main

import (
	gormlike "github.com/survivorbat/gorm-like"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// Employee is the example for normal usage
type Employee struct {
	Name string
}

// RestrictedEmployee is the example for gormlike.TaggedOnly()
type RestrictedEmployee struct {
	// Can be LIKE-d on
	Name string `gormlike:"true"`
	
	// Can NOT be LIKE-d on
	Job string
}

func main() {
	// Normal usage
	db, _ := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	
	filters := map[string]any{
		"name": "%b%",
	}

	db.Use(gormlike.New())
	db.Model(&Employee{}).Where(filters)
	
	// With custom replacement character
	db, _ = gorm.Open(sqlite.Open("test.db"), &gorm.Config{})

	filters := map[string]any{
		"name": "🍌b🍌",
	}

	db.Use(gormlike.New(gormlike.WithCharacter("🍌")))
	db.Model(&Employee{}).Where(filters)
	
	// Only uses LIKE-queries for tagged fields
	db, _ = gorm.Open(sqlite.Open("test.db"), &gorm.Config{})

	filters := map[string]any{
		"name": "🍌b🍌",
	}

	db.Use(gormlike.New(gormlike.TaggedOnly()))
	db.Model(&RestrictedEmployee{}).Where(filters)
}
```

Is automatically turned into a query that looks like this:

```sql
SELECT * FROM employees WHERE name LIKE "%b%";
```

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
