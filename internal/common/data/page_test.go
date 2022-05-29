package data

import (
	"reflect"
	"testing"
)

func TestNewPage(t *testing.T) {
	type args struct {
		content          []int
		page             int
		size             int
		numberOfElements int
	}
	tests := []struct {
		name    string
		args    args
		want    Page[int]
		wantErr bool
	}{
		{
			name: "first_page_success",
			args: args{
				content:          []int{1},
				page:             1,
				size:             5,
				numberOfElements: 16,
			},
			want: Page[int]{
				Content:          []int{1},
				Page:             1,
				Size:             5,
				Total:            4,
				NumberOfElements: 16,
				IsFirst:          true,
				IsLast:           false,
			},
			wantErr: false,
		},
		{
			name: "normal_success",
			args: args{
				content:          []int{1},
				page:             2,
				size:             5,
				numberOfElements: 15,
			},
			want: Page[int]{
				Content:          []int{1},
				Page:             2,
				Size:             5,
				Total:            3,
				NumberOfElements: 15,
				IsFirst:          false,
				IsLast:           false,
			},
			wantErr: false,
		},
		{
			name: "last_page_success",
			args: args{
				content:          []int{1},
				page:             4,
				size:             5,
				numberOfElements: 16,
			},
			want: Page[int]{
				Content:          []int{1},
				Page:             4,
				Size:             5,
				Total:            4,
				NumberOfElements: 16,
				IsFirst:          false,
				IsLast:           true,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewPage(tt.args.content, tt.args.page, tt.args.size, tt.args.numberOfElements)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewPage() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewPage() got = %v, want %v", got, tt.want)
			}
		})
	}
}
