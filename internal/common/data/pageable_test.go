package data

import (
	"reflect"
	"testing"
)

func Test_newPageRequest(t *testing.T) {
	type args struct {
		page         int
		size         int
		sortFields   map[string]string
		chooseFields []string
	}
	tests := []struct {
		name    string
		args    args
		want    Pageable
		wantErr bool
	}{
		{
			name: "first_page_success",
			args: args{
				page:         1,
				size:         20,
				sortFields:   nil,
				chooseFields: nil,
			},
			want: pageRequest{
				page:    1,
				size:    20,
				offset:  0,
				sorts:   nil,
				chooses: nil,
			},
			wantErr: false,
		},
		{
			name: "normal_page_success",
			args: args{
				page:         3,
				size:         20,
				sortFields:   nil,
				chooseFields: nil,
			},
			want: pageRequest{
				page:    3,
				size:    20,
				offset:  40,
				sorts:   nil,
				chooses: nil,
			},
			wantErr: false,
		},
		{
			name: "invalid_page_error",
			args: args{
				page:         0,
				size:         20,
				sortFields:   nil,
				chooseFields: nil,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "invalid_size_error",
			args: args{
				page:         1,
				size:         0,
				sortFields:   nil,
				chooseFields: nil,
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewPageRequest(tt.args.page, tt.args.size, tt.args.sortFields, tt.args.chooseFields)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewPageRequest() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewPageRequest() got = %v, want %v", got, tt.want)
			}
		})
	}
}
