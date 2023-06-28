package gormlike

import (
	"github.com/google/uuid"
	"github.com/ing-bank/gormtestutil"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDeepGorm_Name_ReturnsExpectedName(t *testing.T) {
	t.Parallel()
	// Arrange
	plugin := New()

	// Act
	result := plugin.Name()

	// Assert
	assert.Equal(t, "gormlike", result)
}

func TestDeepGorm_Initialize_RegistersCallback(t *testing.T) {
	t.Parallel()
	// Arrange
	db := gormtestutil.NewMemoryDatabase(t)
	plugin := New()

	// Act
	err := plugin.Initialize(db)

	// Assert
	assert.Nil(t, err)
	assert.NotNil(t, db.Callback().Query().Get("gormlike:query"))
}

type ObjectA struct {
	ID    uuid.UUID
	Name  string
	Age   int
	Other string
}

func TestDeepGorm_Initialize_TriggersLikingCorrectly(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		filter   map[string]any
		existing []ObjectA
		expected []ObjectA
		options  []Option
	}{
		"nothing": {
			expected: []ObjectA{},
		},

		// Check if everything still works
		"simple where query": {
			filter: map[string]any{
				"name": "jessica",
			},
			existing: []ObjectA{{Name: "jessica", Age: 46}, {Name: "amy", Age: 35}},
			expected: []ObjectA{{Name: "jessica", Age: 46}},
		},
		"more complex where query": {
			filter: map[string]any{
				"name": "jessica",
				"age":  53,
			},
			existing: []ObjectA{{Name: "jessica", Age: 53}, {Name: "jessica", Age: 20}},
			expected: []ObjectA{{Name: "jessica", Age: 53}},
		},
		"multi-value where query": {
			filter: map[string]any{
				"name": []string{"jessica", "amy"},
			},
			existing: []ObjectA{{Name: "jessica", Age: 53}, {Name: "amy", Age: 20}},
			expected: []ObjectA{{Name: "jessica", Age: 53}, {Name: "amy", Age: 20}},
		},
		"more complex multi-value where query": {
			filter: map[string]any{
				"name": []string{"jessica", "amy"},
				"age":  []int{53, 20},
			},
			existing: []ObjectA{{Name: "jessica", Age: 53}, {Name: "amy", Age: 20}},
			expected: []ObjectA{{Name: "jessica", Age: 53}, {Name: "amy", Age: 20}},
		},

		// On to the 'real' tests
		"simple like query": {
			filter: map[string]any{
				"name": "%a%",
			},
			existing: []ObjectA{{Name: "jessica", Age: 53}, {Name: "amy", Age: 20}, {Name: "John", Age: 25}},
			expected: []ObjectA{{Name: "jessica", Age: 53}, {Name: "amy", Age: 20}},
		},
		"more complex like query": {
			filter: map[string]any{
				"name": "%a%",
				"age":  20,
			},
			existing: []ObjectA{{Name: "jessica", Age: 53}, {Name: "amy", Age: 20}, {Name: "John", Age: 25}},
			expected: []ObjectA{{Name: "amy", Age: 20}},
		},
		"multi-value, all like queries": {
			filter: map[string]any{
				"name": []string{"%a%", "%o%"},
			},
			existing: []ObjectA{{Name: "jessica", Age: 53}, {Name: "amy", Age: 20}, {Name: "John", Age: 25}},
			expected: []ObjectA{{Name: "jessica", Age: 53}, {Name: "amy", Age: 20}, {Name: "John", Age: 25}},
		},
		"more complex multi-value, all like queries": {
			filter: map[string]any{
				"name":  []string{"%a%", "%o%"},
				"other": []string{"%ooo", "aaa%"},
			},
			existing: []ObjectA{{Name: "jessica", Age: 53, Other: "aaaooo"}, {Name: "amy", Age: 20, Other: "aaaooo"}, {Name: "John", Age: 25, Other: "aaaooo"}},
			expected: []ObjectA{{Name: "jessica", Age: 53, Other: "aaaooo"}, {Name: "amy", Age: 20, Other: "aaaooo"}, {Name: "John", Age: 25, Other: "aaaooo"}},
		},
		"multi-value, some like queries": {
			filter: map[string]any{
				"name": []string{"jessica", "%o%"},
			},
			existing: []ObjectA{{Name: "jessica", Age: 53}, {Name: "amy", Age: 20}, {Name: "John", Age: 25}},
			expected: []ObjectA{{Name: "jessica", Age: 53}, {Name: "John", Age: 25}},
		},
		"more complex multi-value, some like queries": {
			filter: map[string]any{
				"name":  []string{"jessica", "%o%"},
				"other": []string{"aa%", "bb"},
			},
			existing: []ObjectA{{Name: "jessica", Age: 53, Other: "aab"}, {Name: "amy", Age: 20}, {Name: "John", Age: 25, Other: "bb"}},
			expected: []ObjectA{{Name: "jessica", Age: 53, Other: "aab"}, {Name: "John", Age: 25, Other: "bb"}},
		},

		// With custom character
		"simple like query with 🍌": {
			filter: map[string]any{
				"name": "🍌a🍌",
			},
			existing: []ObjectA{{Name: "jessica", Age: 53}, {Name: "amy", Age: 20}, {Name: "John", Age: 25}},
			expected: []ObjectA{{Name: "jessica", Age: 53}, {Name: "amy", Age: 20}},
			options:  []Option{WithCharacter("🍌")},
		},
		"more complex like query with 🍓": {
			filter: map[string]any{
				"name": "🍓a🍓",
				"age":  20,
			},
			existing: []ObjectA{{Name: "jessica", Age: 53}, {Name: "amy", Age: 20}, {Name: "John", Age: 25}},
			expected: []ObjectA{{Name: "amy", Age: 20}},
			options:  []Option{WithCharacter("🍓")},
		},
		"multi-value, all like queries with 🍎": {
			filter: map[string]any{
				"name": []string{"🍎a🍎", "🍎o🍎"},
			},
			existing: []ObjectA{{Name: "jessica", Age: 53}, {Name: "amy", Age: 20}, {Name: "John", Age: 25}},
			expected: []ObjectA{{Name: "jessica", Age: 53}, {Name: "amy", Age: 20}, {Name: "John", Age: 25}},
			options:  []Option{WithCharacter("🍎")},
		},
		"more complex multi-value, all like queries with 🍎": {
			filter: map[string]any{
				"name":  []string{"🍎a🍎", "🍎o🍎"},
				"other": []string{"🍎ooo", "aaa🍎"},
			},
			existing: []ObjectA{{Name: "jessica", Age: 53, Other: "aaaooo"}, {Name: "amy", Age: 20, Other: "aaaooo"}, {Name: "John", Age: 25, Other: "aaaooo"}},
			expected: []ObjectA{{Name: "jessica", Age: 53, Other: "aaaooo"}, {Name: "amy", Age: 20, Other: "aaaooo"}, {Name: "John", Age: 25, Other: "aaaooo"}},
			options:  []Option{WithCharacter("🍎")},
		},
		"multi-value, some like queries with 🍐": {
			filter: map[string]any{
				"name": []string{"jessica", "🍐o🍐"},
			},
			existing: []ObjectA{{Name: "jessica", Age: 53}, {Name: "amy", Age: 20}, {Name: "John", Age: 25}},
			expected: []ObjectA{{Name: "jessica", Age: 53}, {Name: "John", Age: 25}},
			options:  []Option{WithCharacter("🍐")},
		},
		"more complex multi-value, some like queries with 🍐": {
			filter: map[string]any{
				"name":  []string{"jessica", "🍐o🍐"},
				"other": []string{"aa🍐", "bb"},
			},
			existing: []ObjectA{{Name: "jessica", Age: 53, Other: "aab"}, {Name: "amy", Age: 20}, {Name: "John", Age: 25, Other: "bb"}},
			expected: []ObjectA{{Name: "jessica", Age: 53, Other: "aab"}, {Name: "John", Age: 25, Other: "bb"}},
			options:  []Option{WithCharacter("🍐")},
		},
	}

	for name, testData := range tests {
		testData := testData
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			// Arrange
			db := gormtestutil.NewMemoryDatabase(t, gormtestutil.WithName(t.Name())).Debug()
			_ = db.AutoMigrate(&ObjectA{})
			plugin := New(testData.options...)

			if err := db.CreateInBatches(testData.existing, 10).Error; err != nil {
				t.Error(err)
				t.FailNow()
			}

			// Act
			err := db.Use(plugin)

			// Assert
			assert.Nil(t, err)

			var actual []ObjectA
			err = db.Where(testData.filter).Find(&actual).Error
			assert.Nil(t, err)

			assert.Equal(t, testData.expected, actual)
		})
	}
}
