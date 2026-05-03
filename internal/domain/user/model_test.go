package user

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStatus_IsValid(t *testing.T) {
	testCases := []struct {
		name   string
		status Status
		want   bool
	}{
		{name: "active", status: StatusActive, want: true},
		{name: "inactive", status: StatusInactive, want: true},
		{name: "suspended", status: StatusSuspended, want: true},
		{name: "empty", status: "", want: false},
		{name: "unsupported", status: Status("pending"), want: false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.want, tc.status.IsValid())
		})
	}
}
