package data

import (
	"reflect"
	"testing"
)

func Test_newPageRequest(t *testing.T) {
	type args struct {
		page int
		size int
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
				page: 1,
				size: 20,
			},
			want: pageRequest{
				page:   1,
				size:   20,
				offset: 0,
			},
			wantErr: false,
		},
		{
			name: "normal_page_success",
			args: args{
				page: 3,
				size: 20,
			},
			want: pageRequest{
				page:   3,
				size:   20,
				offset: 40,
			},
			wantErr: false,
		},
		{
			name: "invalid_page_error",
			args: args{
				page: 0,
				size: 20,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "invalid_size_error",
			args: args{
				page: 1,
				size: 0,
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := newPageRequest(tt.args.page, tt.args.size, nil, nil, nil)
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
