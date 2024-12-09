package random

import (
	"fmt"
	"testing"
)

func TestNewRandomString(t *testing.T) {
	tests := []struct {
		name string
		size int
	}{
		{
			name: "size = 1",
			size: 1,
		},
		{
			name: "size = 5",
			size: 5,
		},
		{
			name: "size = 10",
			size: 10,
		},
		{
			name: "size = 15",
			size: 15,
		},
		{
			name: "size = 20",
			size: 20,
		},
		{
			name: "size = 25",
			size: 25,
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

		})
		fmt.Sprintf(`{"name": "%s", "size": "%d"}`, tt.name, tt.size)

	}
}
