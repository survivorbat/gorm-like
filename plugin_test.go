package gormlike

import (
	"testing"

	"github.com/ing-bank/gormtestutil"
	"github.com/stretchr/testify/assert"
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
	assert.NoError(t, err)
	assert.NotNil(t, db.Callback().Query().Get("gormlike:query"))
}
