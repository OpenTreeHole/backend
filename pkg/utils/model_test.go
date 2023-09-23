package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestModelsToIDs(t *testing.T) {
	type args struct {
		models any
	}

	type User struct {
		ID int
	}

	tests := []struct {
		name    string
		args    args
		wantIds []int
	}{
		{
			name: "struct slice",
			args: args{
				models: []User{
					{ID: 1},
					{ID: 2},
					{ID: 3},
				},
			},
			wantIds: []int{1, 2, 3},
		},
		{
			name: "struct pointer slice",
			args: args{
				models: []*User{
					{ID: 1},
					{ID: 2},
					{ID: 3},
				},
			},
			wantIds: []int{1, 2, 3},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.wantIds, ModelsToIDs(tt.args.models), "ModelsToIDs(%v)", tt.args.models)
		})
	}
}
