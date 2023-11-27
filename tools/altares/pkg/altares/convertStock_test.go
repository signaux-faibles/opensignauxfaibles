package altares

import (
	"reflect"
	"testing"
)

func Test_removeColumns(t *testing.T) {
	record := []string{"1", "2", "3"}

	type args struct {
		record []string
		remove []int
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{"on retire la colonne 0 sur 3 colonnes", args{record, []int{0}}, []string{"2", "3"}},
		{"on retire les colonnes 0 et 2 sur 3 colonnes", args{record, []int{0, 2}}, []string{"2"}},
		{"on retire une colonne qui n'existe pas", args{record, []int{-1}}, record},
		{"on retire nil sur 3 colonnes", args{record, nil}, record},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := removeColumns(tt.args.record, tt.args.remove...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("removeColumns() = %v, want %v", got, tt.want)
			}
		})
	}
}
