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

		// // Check if everything still works
		"simple where query": {
			filter: map[string]any{
				"name": "jessica",
			},
			query:    defaultQuery,
			existing: []ObjectA{{ID: uuid.MustParse("30611aa6-6fdc-4eb1-b6e2-13485d6c86da"), Name: "jessica", Age: 46}, {ID: uuid.MustParse("f02b0a72-00a5-437b-a8ab-48f033ae3373"), Name: "amy", Age: 35}},
			expected: []ObjectA{{ID: uuid.MustParse("30611aa6-6fdc-4eb1-b6e2-13485d6c86da"), Name: "jessica", Age: 46}},
		},
		"more complex where query": {
			filter: map[string]any{
				"name": "jessica",
				"age":  53,
			},
			query:    defaultQuery,
			existing: []ObjectA{{ID: uuid.MustParse("f297b063-b6aa-4d86-b127-45f59b93969a"), Name: "jessica", Age: 53}, {ID: uuid.MustParse("90c80a13-9c5a-415f-a1da-ba4b4359262f"), Name: "jessica", Age: 20}},
			expected: []ObjectA{{ID: uuid.MustParse("f297b063-b6aa-4d86-b127-45f59b93969a"), Name: "jessica", Age: 53}},
		},
		"multi-value where query": {
			filter: map[string]any{
				"name": []string{"jessica", "amy"},
			},
			query:    defaultQuery,
			existing: []ObjectA{{ID: uuid.MustParse("95ab5cc9-762a-4afc-b465-70d73d83aff3"), Name: "jessica", Age: 53}, {ID: uuid.MustParse("08b78ec4-af61-426a-895f-092494cb9d18"), Name: "amy", Age: 20}},
			expected: []ObjectA{{ID: uuid.MustParse("95ab5cc9-762a-4afc-b465-70d73d83aff3"), Name: "jessica", Age: 53}, {ID: uuid.MustParse("08b78ec4-af61-426a-895f-092494cb9d18"), Name: "amy", Age: 20}},
		},
		"more complex multi-value where query": {
			filter: map[string]any{
				"name": []string{"jessica", "amy"},
				"age":  []int{53, 20},
			},
			query:    defaultQuery,
			existing: []ObjectA{{ID: uuid.MustParse("195de3fc-d378-4c80-8afd-4df1cc11781f"), Name: "jessica", Age: 53}, {ID: uuid.MustParse("bb8b7d42-2b91-4b24-9a2f-61d94d8004b4"), Name: "amy", Age: 20}},
			expected: []ObjectA{{ID: uuid.MustParse("195de3fc-d378-4c80-8afd-4df1cc11781f"), Name: "jessica", Age: 53}, {ID: uuid.MustParse("bb8b7d42-2b91-4b24-9a2f-61d94d8004b4"), Name: "amy", Age: 20}},
		},

		// On to the 'real' tests
		"simple like query": {
			filter: map[string]any{
				"name": "%a%",
			},
			query:    defaultQuery,
			existing: []ObjectA{{ID: uuid.MustParse("3d472489-64a4-4d71-829c-a0c57fab877d"), Name: "jessica", Age: 53}, {ID: uuid.MustParse("74008c97-2bd7-4a93-a81b-0c9bd6adf662"), Name: "amy", Age: 20}, {ID: uuid.MustParse("aa294a48-76c0-4a5c-ae0a-0422f7f7803c"), Name: "John", Age: 25}},
			expected: []ObjectA{{ID: uuid.MustParse("3d472489-64a4-4d71-829c-a0c57fab877d"), Name: "jessica", Age: 53}, {ID: uuid.MustParse("74008c97-2bd7-4a93-a81b-0c9bd6adf662"), Name: "amy", Age: 20}},
		},
		"simple like query on uuid": {
			filter: map[string]any{
				"id": "%22%",
			},
			query:    defaultQuery,
			existing: []ObjectA{{ID: uuid.MustParse("573fc37d-622c-4d41-b48b-201bf195a1c9"), Name: "jessica", Age: 53}, {ID: uuid.MustParse("66df272f-5667-4472-950a-2269833fdbab"), Name: "amy", Age: 20}, {ID: uuid.MustParse("3e3b8dbc-37a3-4113-b05c-6c1d5a3ad5fa"), Name: "John", Age: 25}},
			expected: []ObjectA{{ID: uuid.MustParse("573fc37d-622c-4d41-b48b-201bf195a1c9"), Name: "jessica", Age: 53}, {ID: uuid.MustParse("66df272f-5667-4472-950a-2269833fdbab"), Name: "amy", Age: 20}},
		},
		"more complex like query": {
			filter: map[string]any{
				"name": "a%",
				"age":  20,
			},
			query:    defaultQuery,
			existing: []ObjectA{{ID: uuid.MustParse("25c9b291-951c-4213-b71e-6e3c63ecbe72"), Name: "jessica", Age: 53}, {ID: uuid.MustParse("943ce561-2ae7-40b8-8540-dc49ca84e400"), Name: "amy", Age: 20}, {ID: uuid.MustParse("2d6d5829-cb2b-4a8a-9d09-70ea46b9a62e"), Name: "John", Age: 25}},
			expected: []ObjectA{{ID: uuid.MustParse("943ce561-2ae7-40b8-8540-dc49ca84e400"), Name: "amy", Age: 20}},
		},
		"multi-value, all like queries": {
			filter: map[string]any{
				"name": []string{"%a%", "%o%"},
			},
			query:    defaultQuery,
			existing: []ObjectA{{ID: uuid.MustParse("225c43b8-4645-4603-a362-721f4af370ab"), Name: "jessica", Age: 53}, {ID: uuid.MustParse("b680ef6f-c894-4a4f-a90d-245fb3284c75"), Name: "amy", Age: 20}, {ID: uuid.MustParse("f76c692d-4991-4da2-880e-38e19ca3261e"), Name: "John", Age: 25}},
			expected: []ObjectA{{ID: uuid.MustParse("225c43b8-4645-4603-a362-721f4af370ab"), Name: "jessica", Age: 53}, {ID: uuid.MustParse("b680ef6f-c894-4a4f-a90d-245fb3284c75"), Name: "amy", Age: 20}, {ID: uuid.MustParse("f76c692d-4991-4da2-880e-38e19ca3261e"), Name: "John", Age: 25}},
		},
		"more complex multi-value, all like queries": {
			filter: map[string]any{
				"name":  []string{"%a%", "%o%"},
				"other": []string{"%ooo", "aaa%"},
			},
			query:    defaultQuery,
			existing: []ObjectA{{ID: uuid.MustParse("77c5e4f2-30a9-459b-8178-0a3b61369cdc"), Name: "jessica", Age: 53, Other: "aaaooo"}, {ID: uuid.MustParse("909ec358-ee69-4f74-8e70-069311125e8a"), Name: "amy", Age: 20, Other: "aaaooo"}, {ID: uuid.MustParse("37ac3821-c4e7-4004-b2ed-2cdc7c9ffb53"), Name: "John", Age: 25, Other: "aaaooo"}},
			expected: []ObjectA{{ID: uuid.MustParse("77c5e4f2-30a9-459b-8178-0a3b61369cdc"), Name: "jessica", Age: 53, Other: "aaaooo"}, {ID: uuid.MustParse("909ec358-ee69-4f74-8e70-069311125e8a"), Name: "amy", Age: 20, Other: "aaaooo"}, {ID: uuid.MustParse("37ac3821-c4e7-4004-b2ed-2cdc7c9ffb53"), Name: "John", Age: 25, Other: "aaaooo"}},
		},
		"multi-value, some like queries": {
			filter: map[string]any{
				"name": []string{"jessica", "%o%"},
			},
			query:    defaultQuery,
			existing: []ObjectA{{ID: uuid.MustParse("90cd529f-e172-4264-a047-8836772b3b0c"), Name: "jessica", Age: 53}, {ID: uuid.MustParse("de9db2a9-f300-43d3-8d04-f2f9df676ca3"), Name: "amy", Age: 20}, {ID: uuid.MustParse("516875fe-1dbb-47dc-910d-a5266acec0cf"), Name: "John", Age: 25}},
			expected: []ObjectA{{ID: uuid.MustParse("90cd529f-e172-4264-a047-8836772b3b0c"), Name: "jessica", Age: 53}, {ID: uuid.MustParse("516875fe-1dbb-47dc-910d-a5266acec0cf"), Name: "John", Age: 25}},
		},
		"more complex multi-value, some like queries": {
			filter: map[string]any{
				"name":  []string{"jessica", "%o%"},
				"other": []string{"aa%", "bb"},
			},
			query:    defaultQuery,
			existing: []ObjectA{{ID: uuid.MustParse("dfaa606b-36a5-480f-82ee-43988652c81e"), Name: "jessica", Age: 53, Other: "aab"}, {ID: uuid.MustParse("faae915e-47be-47d0-b6b0-7e3d155537ba"), Name: "amy", Age: 20}, {ID: uuid.MustParse("49409ba0-c165-4a18-ba75-cc1fe77d6b99"), Name: "John", Age: 25, Other: "bb"}},
			expected: []ObjectA{{ID: uuid.MustParse("dfaa606b-36a5-480f-82ee-43988652c81e"), Name: "jessica", Age: 53, Other: "aab"}, {ID: uuid.MustParse("49409ba0-c165-4a18-ba75-cc1fe77d6b99"), Name: "John", Age: 25, Other: "bb"}},
		},
		"explicitly disable liking in query": {
			filter: map[string]any{
				"name":  []string{"jessica", "%o%"},
				"other": []string{"aa%", "bb"},
			},
			query: func(db *gorm.DB) *gorm.DB {
				return db.Set(tagName, false)
			},
			existing: []ObjectA{{ID: uuid.MustParse("9ba6f4e6-b745-44c8-9602-425df7f323fd"), Name: "jessica", Age: 53, Other: "aab"}, {ID: uuid.MustParse("3d0db531-6ea0-4794-b2f6-ac43ec1a1054"), Name: "amy", Age: 20}, {ID: uuid.MustParse("d5a5484f-f0a0-427a-83f1-6e04da969c32"), Name: "John", Age: 25, Other: "bb"}},
			expected: []ObjectA{},
		},

		// With custom character
		"simple like query with 🍌": {
			filter: map[string]any{
				"name": "🍌a🍌",
			},
			query:    defaultQuery,
			existing: []ObjectA{{ID: uuid.MustParse("e8a2e524-dcb5-4436-b608-46f38ab47944"), Name: "jessica", Age: 53}, {ID: uuid.MustParse("a3367b2f-0d13-4e64-afa2-7d50844e370d"), Name: "amy", Age: 20}, {ID: uuid.MustParse("fc5e4af9-5959-4e31-8b56-d33fdb6d71c3"), Name: "John", Age: 25}},
			expected: []ObjectA{{ID: uuid.MustParse("e8a2e524-dcb5-4436-b608-46f38ab47944"), Name: "jessica", Age: 53}, {ID: uuid.MustParse("a3367b2f-0d13-4e64-afa2-7d50844e370d"), Name: "amy", Age: 20}},
			options:  []Option{WithCharacter("🍌")},
		},
		"more complex like query with 🍓": {
			filter: map[string]any{
				"name": "%a%",
				"age":  20,
			},
			query:    defaultQuery,
			existing: []ObjectA{{ID: uuid.MustParse("e9f40044-8579-4d22-8798-bdb74b597de3"), Name: "jessica", Age: 53}, {ID: uuid.MustParse("67b06733-9d2d-4336-86c4-49298da8ebf3"), Name: "amy", Age: 20}, {ID: uuid.MustParse("2e23dbf3-28bb-45f7-ad8b-c02009eed982"), Name: "John", Age: 25}},
			expected: []ObjectA{{ID: uuid.MustParse("67b06733-9d2d-4336-86c4-49298da8ebf3"), Name: "amy", Age: 20}},
		},
		"multi-value, all like queries with 🍎": {
			filter: map[string]any{
				"name": []string{"🍎a🍎", "🍎o🍎"},
			},
			query:    defaultQuery,
			existing: []ObjectA{{ID: uuid.MustParse("c1c04d44-e97a-4e0a-bcc5-7a33546654a1"), Name: "jessica", Age: 53}, {ID: uuid.MustParse("82752fa8-a166-47d3-aacd-70344441f65f"), Name: "amy", Age: 20}, {ID: uuid.MustParse("6a62f60c-b38d-4c65-9d91-0204d2dad639"), Name: "John", Age: 25}},
			expected: []ObjectA{{ID: uuid.MustParse("c1c04d44-e97a-4e0a-bcc5-7a33546654a1"), Name: "jessica", Age: 53}, {ID: uuid.MustParse("82752fa8-a166-47d3-aacd-70344441f65f"), Name: "amy", Age: 20}, {ID: uuid.MustParse("6a62f60c-b38d-4c65-9d91-0204d2dad639"), Name: "John", Age: 25}},
			options:  []Option{WithCharacter("🍎")},
		},
		"more complex multi-value, all like queries with 🍎": {
			filter: map[string]any{
				"name":  []string{"🍎a🍎", "🍎o🍎"},
				"other": []string{"🍎ooo", "aaa🍎"},
			},
			query:    defaultQuery,
			existing: []ObjectA{{ID: uuid.MustParse("48b38043-34e2-45a3-8757-64680765c201"), Name: "jessica", Age: 53, Other: "aaaooo"}, {ID: uuid.MustParse("f4c295e4-6fe4-49c0-9900-af7633951be0"), Name: "amy", Age: 20, Other: "aaaooo"}, {ID: uuid.MustParse("5075359a-01fd-4633-afc3-9d6914f0e86c"), Name: "John", Age: 25, Other: "aaaooo"}},
			expected: []ObjectA{{ID: uuid.MustParse("48b38043-34e2-45a3-8757-64680765c201"), Name: "jessica", Age: 53, Other: "aaaooo"}, {ID: uuid.MustParse("f4c295e4-6fe4-49c0-9900-af7633951be0"), Name: "amy", Age: 20, Other: "aaaooo"}, {ID: uuid.MustParse("5075359a-01fd-4633-afc3-9d6914f0e86c"), Name: "John", Age: 25, Other: "aaaooo"}},
			options:  []Option{WithCharacter("🍎")},
		},
		"multi-value, some like queries with 🍐": {
			filter: map[string]any{
				"name": []string{"jessica", "🍐o🍐"},
			},
			query:    defaultQuery,
			existing: []ObjectA{{ID: uuid.MustParse("9800cafe-a92a-4805-9196-9f17fcea42cb"), Name: "jessica", Age: 53}, {ID: uuid.MustParse("22208fb7-9685-4119-9ef3-dea673ab506b"), Name: "amy", Age: 20}, {ID: uuid.MustParse("d46ee447-5ce8-4943-9e90-0c2359d60076"), Name: "John", Age: 25}},
			expected: []ObjectA{{ID: uuid.MustParse("9800cafe-a92a-4805-9196-9f17fcea42cb"), Name: "jessica", Age: 53}, {ID: uuid.MustParse("d46ee447-5ce8-4943-9e90-0c2359d60076"), Name: "John", Age: 25}},
			options:  []Option{WithCharacter("🍐")},
		},
		"more complex multi-value, some like queries with 🍐": {
			filter: map[string]any{
				"name":  []string{"jessica", "🍐o🍐"},
				"other": []string{"aa🍐", "bc"},
			},
			query:    defaultQuery,
			existing: []ObjectA{{ID: uuid.MustParse("d72880ca-baf7-4fb4-8c2a-070e5274a89e"), Name: "jessica", Age: 53, Other: "aab"}, {ID: uuid.MustParse("c45363aa-8e91-497b-afc1-b5ba0642954f"), Name: "amy", Age: 20, Other: "bc"}, {ID: uuid.MustParse("e2befb93-030a-419f-83d5-b9b740fde262"), Name: "John", Age: 25, Other: "bb"}},
			expected: []ObjectA{{ID: uuid.MustParse("d72880ca-baf7-4fb4-8c2a-070e5274a89e"), Name: "jessica", Age: 53, Other: "aab"}},
			options:  []Option{WithCharacter("🍐")},
		},

		// With existing query
		"simple like query with existing calls": {
			filter: map[string]any{
				"name": "%a%",
			},
			query: func(db *gorm.DB) *gorm.DB {
				return db.Where("other = ?", "goodbye")
			},
			existing: []ObjectA{{ID: uuid.MustParse("bf418b07-fc50-4a74-8411-96360f73d530"), Name: "jessica", Age: 53, Other: "hello"}, {ID: uuid.MustParse("869584f6-ebeb-4a82-8c1b-701798284d69"), Name: "amy", Age: 20, Other: "goodbye"}, {ID: uuid.MustParse("dfee7228-9d8c-49a8-b92b-f2ec97f542f6"), Name: "John", Age: 25}},
			expected: []ObjectA{{ID: uuid.MustParse("869584f6-ebeb-4a82-8c1b-701798284d69"), Name: "amy", Age: 20, Other: "goodbye"}},
		},
		"more complex like query with existing calls": {
			filter: map[string]any{
				"name": "%a%",
				"age":  20,
			},
			query: func(db *gorm.DB) *gorm.DB {
				return db.Where("other = ?", "def")
			},
			existing: []ObjectA{{ID: uuid.MustParse("63af0e59-0c7f-42a5-a08d-c9b22e7b5bbf"), Name: "jessica", Age: 53, Other: "abc"}, {ID: uuid.MustParse("fc101a07-e250-4056-86dd-eb4719337538"), Name: "amy", Age: 20, Other: "def"}, {ID: uuid.MustParse("1257821f-cac6-4b55-a00d-e45e0f3ac2f4"), Name: "John", Age: 25, Other: "ghi"}},
			expected: []ObjectA{{ID: uuid.MustParse("fc101a07-e250-4056-86dd-eb4719337538"), Name: "amy", Age: 20, Other: "def"}},
		},
		"multi-value, all like queries with existing calls": {
			filter: map[string]any{
				"name": []string{"%a%", "%o%", "%e%"},
			},
			query: func(db *gorm.DB) *gorm.DB {
				return db.Where("other = ?", "no").Or("other = ?", "yes").Or("other = ?", "maybe")
			},
			existing: []ObjectA{{ID: uuid.MustParse("9ab327ad-42fa-46dd-b061-9583a297e8a6"), Name: "jessica", Age: 53, Other: "no"}, {ID: uuid.MustParse("a4ec1630-fa8e-4556-93cc-143522af072b"), Name: "amy", Age: 20, Other: "yes"}, {ID: uuid.MustParse("0e0ef81a-4afc-42ad-8387-bc005ecad5c5"), Name: "John", Age: 25, Other: "maybe"}},
			expected: []ObjectA{{ID: uuid.MustParse("9ab327ad-42fa-46dd-b061-9583a297e8a6"), Name: "jessica", Age: 53, Other: "no"}, {ID: uuid.MustParse("a4ec1630-fa8e-4556-93cc-143522af072b"), Name: "amy", Age: 20, Other: "yes"}, {ID: uuid.MustParse("0e0ef81a-4afc-42ad-8387-bc005ecad5c5"), Name: "John", Age: 25, Other: "maybe"}},
		},
		"more complex multi-value, all like queries with existing calls": {
			filter: map[string]any{
				"other": []string{"aaa%", "%ooo"},
			},
			query: func(db *gorm.DB) *gorm.DB {
				return db.Where("name LIKE ?", "%a%").Or("name LIKE ?", "%o%")
			},
			existing: []ObjectA{{ID: uuid.MustParse("2b0108fe-7639-4791-b209-cd605040455a"), Name: "jessica", Age: 53, Other: "aaaooo"}, {ID: uuid.MustParse("7343f441-afa2-4248-be7f-8a2d18227a52"), Name: "amy", Age: 20, Other: "aaaooo"}, {ID: uuid.MustParse("340d161c-b3c6-4209-8782-f0a9963d19eb"), Name: "John", Age: 25, Other: "aaaooo"}},
			expected: []ObjectA{{ID: uuid.MustParse("2b0108fe-7639-4791-b209-cd605040455a"), Name: "jessica", Age: 53, Other: "aaaooo"}, {ID: uuid.MustParse("7343f441-afa2-4248-be7f-8a2d18227a52"), Name: "amy", Age: 20, Other: "aaaooo"}, {ID: uuid.MustParse("340d161c-b3c6-4209-8782-f0a9963d19eb"), Name: "John", Age: 25, Other: "aaaooo"}},
		},
		"multi-value, some like queries with existing calls": {
			filter: map[string]any{
				"name": []string{"%essica"},
			},
			query: func(db *gorm.DB) *gorm.DB {
				return db.Where("name = ?", "jessica")
			},
			existing: []ObjectA{{ID: uuid.MustParse("1ba784aa-0b7a-4098-9a89-946a18e33c3b"), Name: "jessica", Age: 53}, {ID: uuid.MustParse("41531e1f-79d3-4100-b89f-83eba11a03fa"), Name: "amy", Age: 20}, {ID: uuid.MustParse("7b575b5e-10f6-4662-ae25-c66c6cdfc921"), Name: "John", Age: 25}},
			expected: []ObjectA{{ID: uuid.MustParse("1ba784aa-0b7a-4098-9a89-946a18e33c3b"), Name: "jessica", Age: 53}},
		},
		"more complex multi-value, some like queries with existing calls": {
			filter: map[string]any{
				"other": []string{"aa%", "bb"},
			},
			query: func(db *gorm.DB) *gorm.DB {
				return db.Where("name = ?", "jessica").Or("name LIKE ?", "%o%")
			},
			existing: []ObjectA{{ID: uuid.MustParse("3a8bd18b-ba4e-430c-95f5-0a517fe2dc1a"), Name: "jessica", Age: 53, Other: "aab"}, {ID: uuid.MustParse("9bd65c44-0a99-4721-8b7b-bd510fdaa3d7"), Name: "amy", Age: 20}, {ID: uuid.MustParse("8f6cf7f5-4a8e-4e6c-bcc3-1f3a81bf8cdb"), Name: "John", Age: 25, Other: "bb"}},
			expected: []ObjectA{{ID: uuid.MustParse("3a8bd18b-ba4e-430c-95f5-0a517fe2dc1a"), Name: "jessica", Age: 53, Other: "aab"}, {ID: uuid.MustParse("8f6cf7f5-4a8e-4e6c-bcc3-1f3a81bf8cdb"), Name: "John", Age: 25, Other: "bb"}},
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
			dbQuery := testData.query(db).Where(testData.filter).Find(&actual)
			err = dbQuery.Error
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
