package dummyjson

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_productClient_GetNewProducts(t *testing.T) {
	client := NewProductClient("https://dummyjson.com")
	dto, err := client.GetNewProducts(context.TODO(), 10)
	require.NoError(t, err)
	require.NotEqual(t, len(dto), 0)
}
