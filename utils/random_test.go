package utils

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRandomString(t *testing.T) {
	length := 10
	str := RandomString(length)

	require.NotEmpty(t, str)
	require.Len(t, str, length)
}
