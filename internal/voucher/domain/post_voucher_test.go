package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVoucher_Post(t *testing.T) {
	type fields struct {
		isReviewed bool
		isAudited  bool
		isPosted   bool
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "normal_success",
			fields: fields{
				isReviewed: true,
				isAudited:  true,
				isPosted:   false,
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				assert.NoError(t, err)
				return false
			},
		},
		{
			name: "notReviewed_error",
			fields: fields{
				isReviewed: false,
				isAudited:  true,
				isPosted:   false,
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				assert.Equal(t, errPostNotReviewed, err.(domainErr).Slug())
				return true
			},
		},
		{
			name: "notAudited_error",
			fields: fields{
				isReviewed: true,
				isAudited:  false,
				isPosted:   false,
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				assert.Equal(t, errPostNotAudited, err.(domainErr).Slug())
				return true
			},
		},
		{
			name: "repeatPost_error",
			fields: fields{
				isReviewed: true,
				isAudited:  true,
				isPosted:   true,
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				assert.Equal(t, errPostRepeatPost, err.(domainErr).Slug())
				return true
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := &Voucher{
				isReviewed: tt.fields.isReviewed,
				isAudited:  tt.fields.isAudited,
				isPosted:   tt.fields.isPosted,
			}
			tt.wantErr(t, v.Post(), "Post()")
		})
	}
}
