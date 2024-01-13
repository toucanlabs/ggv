package parser

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParse(t *testing.T) {

	result := NewParser().Parse("github.com/gin-gonic/gin")

	assert.Equal(t, "", result)

}
