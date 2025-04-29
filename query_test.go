package gormlike

import (
	"testing"

	"github.com/google/uuid"
	"github.com/ing-bank/gormtestutil"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

// nolint:maintidx // Acceptable
func TestGormLike_Initialize_TriggersLikingCorrectly(t *testing.T) {
	t.Parallel()

	type ObjectA struct {
		ID    uuid.UUID
		Name  string
		Age   int
		Other string
	}

	jessica1 := ObjectA{
		ID:   uuid.MustParse("30611aa6-6fdc-4eb1-b6e2-13485d6c86da"),
		Name: "jessica",
		Age:  53,
	}
	jessica2 := ObjectA{
		ID:   uuid.MustParse("90c80a13-9c5a-415f-a1da-ba4b4359262f"),
		Name: "jessica",
		Age:  20,
	}
	amy := ObjectA{
		ID:   uuid.MustParse("f02b0a72-00a5-437b-a8ab-48f033ae3373"),
		Name: "amy",
		Age:  20,
	}
	john := ObjectA{
		ID:   uuid.MustParse("aa294a48-76c0-4a5c-ae0a-0422f7f7803c"),
		Name: "John",
		Age:  25,
	}

	defaultQuery := func(db *gorm.DB) *gorm.DB { return db }

	tests := map[string]struct {
		filter   map[string]any
		existing []ObjectA
		options  []Option
		query    func(*gorm.DB) *gorm.DB

		expected []ObjectA
	}{
		"nothing": {
			expected: []ObjectA{},
			query:    defaultQuery,
		},

		// Check if everything still works
		"simple where query": {
			filter: map[string]any{
				"name": "jessica",
			},
			query:    defaultQuery,
			existing: []ObjectA{jessica1, amy},
			expected: []ObjectA{jessica1},
		},
		"more complex where query": {
			filter: map[string]any{
				"name": "jessica",
				"age":  53,
			},
			query:    defaultQuery,
			existing: []ObjectA{jessica1, jessica2},
			expected: []ObjectA{jessica1},
		},
		"multi-value where query": {
			filter: map[string]any{
				"name": []string{"jessica", "amy"},
			},
			query:    defaultQuery,
			existing: []ObjectA{jessica1, amy},
			expected: []ObjectA{jessica1, amy},
		},
		"more complex multi-value where query": {
			filter: map[string]any{
				"name": []string{"jessica", "amy"},
				"age":  []int{53, 20},
			},
			query:    defaultQuery,
			existing: []ObjectA{jessica1, amy},
			expected: []ObjectA{jessica1, amy},
		},

		// On to the 'real' tests
		"simple like query": {
			filter: map[string]any{
				"name": "%a%",
			},
			query:    defaultQuery,
			existing: []ObjectA{jessica1, amy, john},
			expected: []ObjectA{jessica1, amy},
		},
		"simple like query on uuid": {
			filter: map[string]any{
				"id": "%aa%",
			},
			query:    defaultQuery,
			existing: []ObjectA{jessica1, amy, john},
			expected: []ObjectA{jessica1, john},
		},
		"more complex like query": {
			filter: map[string]any{
				"name": []string{"%a%"},
				"age":  []int{20, 25},
			},
			query:    defaultQuery,
			existing: []ObjectA{jessica1, amy, john},
			expected: []ObjectA{amy},
		},
		"multi-value, all like queries": {
			filter: map[string]any{
				"name": []string{"%a%", "%o%"},
			},
			query:    defaultQuery,
			existing: []ObjectA{jessica1, amy, john},
			expected: []ObjectA{jessica1, amy, john},
		},
		"more complex multi-value, all like queries": {
			filter: map[string]any{
				"name":  []string{"%a%", "%o%"},
				"other": []string{"%ooo", "aaa%"},
			},
			query: defaultQuery,
			existing: []ObjectA{
				{ID: jessica1.ID, Name: jessica1.Name, Age: jessica1.Age, Other: "aaaooo"},
				{ID: amy.ID, Name: amy.Name, Age: amy.Age, Other: "aaaooo"},
				{ID: john.ID, Name: john.Name, Age: john.Age, Other: "aaaooo"},
			},
			expected: []ObjectA{
				{ID: jessica1.ID, Name: jessica1.Name, Age: jessica1.Age, Other: "aaaooo"},
				{ID: amy.ID, Name: amy.Name, Age: amy.Age, Other: "aaaooo"},
				{ID: john.ID, Name: john.Name, Age: john.Age, Other: "aaaooo"},
			},
		},
		"multi-value, some like queries": {
			filter: map[string]any{
				"name": []string{"jessica", "%o%"},
			},
			query:    defaultQuery,
			existing: []ObjectA{jessica1, amy, john},
			expected: []ObjectA{jessica1, john},
		},
		"more complex multi-value, some like queries": {
			filter: map[string]any{
				"name":  []string{"jessica", "%o%"},
				"other": []string{"aa%", "bb"},
			},
			query: defaultQuery,
			existing: []ObjectA{
				{ID: jessica1.ID, Name: jessica1.Name, Age: jessica1.Age, Other: "aab"},
				{ID: amy.ID, Name: amy.Name, Age: amy.Age},
				{ID: john.ID, Name: john.Name, Age: john.Age, Other: "bb"},
			},
			expected: []ObjectA{
				{ID: jessica1.ID, Name: jessica1.Name, Age: jessica1.Age, Other: "aab"},
				{ID: john.ID, Name: john.Name, Age: john.Age, Other: "bb"},
			},
		},
		"explicitly disable liking in query": {
			filter: map[string]any{
				"name":  []string{"jessica", "%o%"},
				"other": []string{"aa%", "bb"},
			},
			query: func(db *gorm.DB) *gorm.DB {
				return db.Set(tagName, false)
			},
			existing: []ObjectA{
				{ID: jessica1.ID, Name: jessica1.Name, Age: jessica1.Age, Other: "aab"},
				{ID: amy.ID, Name: amy.Name, Age: amy.Age},
				{ID: john.ID, Name: john.Name, Age: john.Age, Other: "bb"},
			},
			expected: []ObjectA{},
		},

		// With custom character
		"simple like query with 🍌": {
			filter: map[string]any{
				"name": "🍌a🍌",
			},
			query:    defaultQuery,
			existing: []ObjectA{jessica1, amy, john},
			expected: []ObjectA{jessica1, amy},
			options:  []Option{WithCharacter("🍌")},
		},
		"more complex like query with 🍓": {
			filter: map[string]any{
				"name": "🍓a🍓",
				"age":  20,
			},
			query:    defaultQuery,
			existing: []ObjectA{jessica1, amy, john},
			expected: []ObjectA{amy},
			options:  []Option{WithCharacter("🍓")},
		},
		"multi-value, all like queries with 🍎": {
			filter: map[string]any{
				"name": []string{"🍎a🍎", "🍎o🍎"},
			},
			query:    defaultQuery,
			existing: []ObjectA{jessica1, amy, john},
			expected: []ObjectA{jessica1, amy, john},
			options:  []Option{WithCharacter("🍎")},
		},
		"more complex multi-value, all like queries with 🍎": {
			filter: map[string]any{
				"name":  []string{"🍎a🍎", "🍎o🍎"},
				"other": []string{"🍎ooo", "aaa🍎"},
			},
			query: defaultQuery,
			existing: []ObjectA{
				{ID: jessica1.ID, Name: jessica1.Name, Age: jessica1.Age, Other: "aaaooo"},
				{ID: amy.ID, Name: amy.Name, Age: amy.Age, Other: "aaaooo"},
				{ID: john.ID, Name: john.Name, Age: john.Age, Other: "aaaooo"},
			},
			expected: []ObjectA{
				{ID: jessica1.ID, Name: jessica1.Name, Age: jessica1.Age, Other: "aaaooo"},
				{ID: amy.ID, Name: amy.Name, Age: amy.Age, Other: "aaaooo"},
				{ID: john.ID, Name: john.Name, Age: john.Age, Other: "aaaooo"},
			},
			options: []Option{WithCharacter("🍎")},
		},
		"multi-value, some like queries with 🍐": {
			filter: map[string]any{
				"name": []string{"jessica", "🍐o🍐"},
			},
			query:    defaultQuery,
			existing: []ObjectA{jessica1, amy, john},
			expected: []ObjectA{jessica1, john},
			options:  []Option{WithCharacter("🍐")},
		},
		"more complex multi-value, some like queries with 🍐": {
			filter: map[string]any{
				"name":  []string{"jessica", "🍐o🍐"},
				"other": []string{"aa🍐", "bc"},
			},
			query: defaultQuery,
			existing: []ObjectA{
				{ID: jessica1.ID, Name: jessica1.Name, Age: jessica1.Age, Other: "aab"},
				{ID: amy.ID, Name: amy.Name, Age: amy.Age, Other: "bc"},
				{ID: john.ID, Name: john.Name, Age: john.Age, Other: "bb"},
			},
			expected: []ObjectA{
				{ID: jessica1.ID, Name: jessica1.Name, Age: jessica1.Age, Other: "aab"},
			},
			options: []Option{WithCharacter("🍐")},
		},

		// With existing query
		"simple like query with existing calls": {
			filter: map[string]any{
				"name": "%a%",
			},
			query: func(db *gorm.DB) *gorm.DB {
				return db.Where("other = ?", "goodbye")
			},
			existing: []ObjectA{
				{ID: jessica1.ID, Name: jessica1.Name, Age: jessica1.Age, Other: "hello"},
				{ID: amy.ID, Name: amy.Name, Age: amy.Age, Other: "goodbye"},
				{ID: john.ID, Name: john.Name, Age: john.Age},
			},
			expected: []ObjectA{
				{ID: amy.ID, Name: amy.Name, Age: amy.Age, Other: "goodbye"},
			},
		},
		"more complex like query with existing calls": {
			filter: map[string]any{
				"name": "%a%",
				"age":  20,
			},
			query: func(db *gorm.DB) *gorm.DB {
				return db.Where("other = ?", "def")
			},
			existing: []ObjectA{
				{ID: jessica1.ID, Name: jessica1.Name, Age: jessica1.Age, Other: "abc"},
				{ID: amy.ID, Name: amy.Name, Age: amy.Age, Other: "def"},
				{ID: john.ID, Name: john.Name, Age: john.Age, Other: "ghi"},
			},
			expected: []ObjectA{
				{ID: amy.ID, Name: amy.Name, Age: amy.Age, Other: "def"},
			},
		},
		"multi-value, all like queries with existing calls": {
			filter: map[string]any{
				"name": []string{"%a%", "%o%", "%e%"},
			},
			query: func(db *gorm.DB) *gorm.DB {
				return db.Where("other = ?", "no").Or("other = ?", "yes").Or("other = ?", "maybe")
			},
			existing: []ObjectA{
				{ID: jessica1.ID, Name: jessica1.Name, Age: jessica1.Age, Other: "no"},
				{ID: amy.ID, Name: amy.Name, Age: amy.Age, Other: "yes"},
				{ID: john.ID, Name: john.Name, Age: john.Age, Other: "maybe"},
			},
			expected: []ObjectA{
				{ID: jessica1.ID, Name: jessica1.Name, Age: jessica1.Age, Other: "no"},
				{ID: amy.ID, Name: amy.Name, Age: amy.Age, Other: "yes"},
				{ID: john.ID, Name: john.Name, Age: john.Age, Other: "maybe"},
			},
		},
		"more complex multi-value, all like queries with existing calls": {
			filter: map[string]any{
				"other": []string{"aaa%", "%ooo"},
			},
			query: func(db *gorm.DB) *gorm.DB {
				return db.Where("name LIKE ?", "%a%").Or("name LIKE ?", "%o%")
			},
			existing: []ObjectA{
				{ID: jessica1.ID, Name: jessica1.Name, Age: jessica1.Age, Other: "aaaooo"},
				{ID: amy.ID, Name: amy.Name, Age: amy.Age, Other: "aaaooo"},
				{ID: john.ID, Name: john.Name, Age: john.Age, Other: "aaaooo"},
			},
			expected: []ObjectA{
				{ID: jessica1.ID, Name: jessica1.Name, Age: jessica1.Age, Other: "aaaooo"},
				{ID: amy.ID, Name: amy.Name, Age: amy.Age, Other: "aaaooo"},
				{ID: john.ID, Name: john.Name, Age: john.Age, Other: "aaaooo"},
			},
		},
		"multi-value, some like queries with existing calls": {
			filter: map[string]any{
				"name": []string{"%essica"},
			},
			query: func(db *gorm.DB) *gorm.DB {
				return db.Where("name = ?", "jessica")
			},
			existing: []ObjectA{jessica1, amy, john},
			expected: []ObjectA{jessica1},
		},
		"more complex multi-value, some like queries with existing calls": {
			filter: map[string]any{
				"other": []string{"aa%", "bb"},
			},
			query: func(db *gorm.DB) *gorm.DB {
				return db.Where("name = ?", "jessica").Or("name LIKE ?", "%o%")
			},
			existing: []ObjectA{
				{ID: jessica1.ID, Name: jessica1.Name, Age: jessica1.Age, Other: "aab"},
				{ID: amy.ID, Name: amy.Name, Age: amy.Age, Other: "bc"},
				{ID: john.ID, Name: john.Name, Age: john.Age, Other: "bb"},
			},
			expected: []ObjectA{
				{ID: jessica1.ID, Name: jessica1.Name, Age: jessica1.Age, Other: "aab"},
				{ID: john.ID, Name: john.Name, Age: john.Age, Other: "bb"},
			},
		},
	}

	for name, testData := range tests {
		testData := testData
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			// Arrange
			db := gormtestutil.NewMemoryDatabase(t, gormtestutil.WithName(t.Name()))
			_ = db.AutoMigrate(&ObjectA{})
			plugin := New(testData.options...)

			if err := db.CreateInBatches(testData.existing, 10).Error; err != nil {
				t.Error(err)
				t.FailNow()
			}

			// Act
			err := db.Use(plugin)

			// Assert
			assert.NoError(t, err)

			var actual []ObjectA

			err = testData.query(db).Where(testData.filter).Find(&actual).Error
			assert.NoError(t, err)

			assert.Equal(t, testData.expected, actual)
		})
	}
}

func TestGormLike_Initialize_TriggersLikingCorrectlyWithConditionalTag(t *testing.T) {
	t.Parallel()

	type ObjectB struct {
		Name  string `gormlike:"true"`
		Other string
	}

	tests := map[string]struct {
		filter   map[string]any
		existing []ObjectB
		expected []ObjectB
	}{
		"simple filter on allowed fields": {
			filter: map[string]any{
				"name": "jes%",
			},
			existing: []ObjectB{{Name: "jessica", Other: "abc"}, {Name: "amy", Other: "def"}},
			expected: []ObjectB{{Name: "jessica", Other: "abc"}},
		},
		"simple filter on disallowed fields": {
			filter: map[string]any{
				"other": "%b%",
			},
			existing: []ObjectB{{Name: "jessica", Other: "abc"}, {Name: "jessica", Other: "abc"}},
			expected: []ObjectB{},
		},
		"multi-filter on allowed fields": {
			filter: map[string]any{
				"name": []string{"jes%", "%my"},
			},
			existing: []ObjectB{{Name: "jessica", Other: "abc"}, {Name: "amy", Other: "def"}},
			expected: []ObjectB{{Name: "jessica", Other: "abc"}, {Name: "amy", Other: "def"}},
		},
		"multi-filter on disallowed fields": {
			filter: map[string]any{
				"other": []string{"%b%", "%c%"},
			},
			existing: []ObjectB{{Name: "jessica", Other: "abc"}, {Name: "jessica", Other: "abc"}},
			expected: []ObjectB{},
		},
	}

	for name, testData := range tests {
		testData := testData
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			// Arrange
			db := gormtestutil.NewMemoryDatabase(t, gormtestutil.WithName(t.Name()))
			_ = db.AutoMigrate(&ObjectB{})
			plugin := New(TaggedOnly())

			if err := db.CreateInBatches(testData.existing, 10).Error; err != nil {
				t.Error(err)
				t.FailNow()
			}

			// Act
			err := db.Use(plugin)

			// Assert
			assert.NoError(t, err)

			var actual []ObjectB
			err = db.Where(testData.filter).Find(&actual).Error
			assert.NoError(t, err)

			assert.Equal(t, testData.expected, actual)
		})
	}
}

func TestGormLike_Initialize_AlwaysIgnoresFieldsWithGormLikeFalse(t *testing.T) {
	t.Parallel()

	type ObjectB struct {
		Name  string
		Other string `gormlike:"false"`
	}

	tests := map[string]struct {
		filter   map[string]any
		existing []ObjectB
		expected []ObjectB
	}{
		"Normal filter works on never field": {
			filter: map[string]any{
				"other": "abc",
			},
			existing: []ObjectB{{Name: "jessica", Other: "abc"}, {Name: "jessica", Other: "abc"}},
			expected: []ObjectB{{Name: "jessica", Other: "abc"}, {Name: "jessica", Other: "abc"}},
		},
		"simple filter on disallowed fields": {
			filter: map[string]any{
				"other": "%b%",
			},
			existing: []ObjectB{{Name: "jessica", Other: "abc"}, {Name: "jessica", Other: "abc"}},
			expected: []ObjectB{},
		},
		"multi-filter on disallowed fields": {
			filter: map[string]any{
				"other": []string{"%b%", "%c%"},
			},
			existing: []ObjectB{{Name: "jessica", Other: "abc"}, {Name: "jessica", Other: "abc"}},
			expected: []ObjectB{},
		},
	}

	for name, testData := range tests {
		testData := testData
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			// Arrange
			db := gormtestutil.NewMemoryDatabase(t, gormtestutil.WithName(t.Name()))
			_ = db.AutoMigrate(&ObjectB{})
			plugin := New(TaggedOnly())

			if err := db.CreateInBatches(testData.existing, 10).Error; err != nil {
				t.Error(err)
				t.FailNow()
			}

			// Act
			err := db.Use(plugin)

			// Assert
			assert.NoError(t, err)

			var actual []ObjectB
			err = db.Where(testData.filter).Find(&actual).Error
			assert.NoError(t, err)

			assert.Equal(t, testData.expected, actual)
		})
	}
}

func TestGormLike_Initialize_TriggersLikingCorrectlyWithSetting(t *testing.T) {
	t.Parallel()

	type ObjectB struct {
		Name  string
		Other string
	}

	tests := map[string]struct {
		filter   map[string]any
		query    func(*gorm.DB) *gorm.DB
		existing []ObjectB
		expected []ObjectB
	}{
		"like with query set to true": {
			filter: map[string]any{
				"name": "jes%",
			},
			query: func(db *gorm.DB) *gorm.DB {
				return db.Set(tagName, true)
			},
			existing: []ObjectB{{Name: "jessica", Other: "abc"}},
			expected: []ObjectB{{Name: "jessica", Other: "abc"}},
		},
		"like with query set to false": {
			filter: map[string]any{
				"name": "jes%",
			},
			query: func(db *gorm.DB) *gorm.DB {
				return db.Set(tagName, false)
			},
			existing: []ObjectB{{Name: "jessica", Other: "abc"}},
			expected: []ObjectB{},
		},
		"like with query set to random value": {
			filter: map[string]any{
				"name": "jes%",
			},
			query: func(db *gorm.DB) *gorm.DB {
				return db.Set(tagName, "yes")
			},
			existing: []ObjectB{{Name: "jessica", Other: "abc"}},
			expected: []ObjectB{},
		},
		"like with query unset": {
			filter: map[string]any{
				"name": "jes%",
			},
			query: func(db *gorm.DB) *gorm.DB {
				return db
			},
			existing: []ObjectB{{Name: "jessica", Other: "abc"}},
			expected: []ObjectB{},
		},
	}

	for name, testData := range tests {
		testData := testData
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			// Arrange
			db := gormtestutil.NewMemoryDatabase(t, gormtestutil.WithName(t.Name()))
			_ = db.AutoMigrate(&ObjectB{})
			plugin := New(SettingOnly())

			if err := db.CreateInBatches(testData.existing, 10).Error; err != nil {
				t.Error(err)
				t.FailNow()
			}

			db = testData.query(db)

			// Act
			err := db.Use(plugin)

			// Assert
			assert.NoError(t, err)

			var actual []ObjectB
			err = db.Where(testData.filter).Find(&actual).Error
			assert.NoError(t, err)

			assert.Equal(t, testData.expected, actual)
		})
	}
}

func TestGormLike_Initialize_ProcessUnknownFields(t *testing.T) {
	t.Parallel()

	type ObjectB struct {
		Name  string
		Other string
	}

	tests := map[string]struct {
		filter   map[string]any
		query    func(*gorm.DB) *gorm.DB
		existing []ObjectB
		expected []ObjectB
	}{
		"like with unknown field": {
			filter: map[string]any{
				"name":          "jes%",
				"unknown_field": false,
			},
			query: func(db *gorm.DB) *gorm.DB {
				return db.Set(tagName, true)
			},
			existing: []ObjectB{{Name: "jessica", Other: "abc"}},
			expected: []ObjectB{},
		},
	}

	for name, testData := range tests {
		testData := testData
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			// Arrange
			db := gormtestutil.NewMemoryDatabase(t, gormtestutil.WithName(t.Name()))
			_ = db.AutoMigrate(&ObjectB{})
			plugin := New(SettingOnly())

			if err := db.CreateInBatches(testData.existing, 10).Error; err != nil {
				t.Error(err)
				t.FailNow()
			}

			db = testData.query(db)

			// Act
			err := db.Use(plugin)

			// Assert
			assert.NoError(t, err)

			var actual []ObjectB
			err = db.Where(testData.filter).Find(&actual).Error
			assert.Equal(t, "no such column: unknown_field", err.Error())
			assert.Nil(t, actual)
		})
	}
}
