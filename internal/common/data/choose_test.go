package data

import (
	"reflect"
	"testing"
)

func Test_newChoose(t *testing.T) {
	type args struct {
		field string
	}
	tests := []struct {
		name    string
		args    args
		want    Choose
		wantErr bool
	}{
		{
			name:    "normal_success",
			args:    args{field: "fieldInCamelCaseSnake"},
			want:    chooseRequest{field: "field_in_camel_case_snake"},
			wantErr: false,
		},
		{
			name:    "invalid_field_error",
			args:    args{field: "fieldWith Space"},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "invalid_field_error_2",
			args:    args{field: "fieldWith%sInvalidChar"},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := newChoose(tt.args.field)
			if (err != nil) != tt.wantErr {
				t.Errorf("newChoose() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newChoose() got = %v, want %v", got, tt.want)
			}
		})
	}
}
