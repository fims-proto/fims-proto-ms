package domain

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVoucher_Audit(t *testing.T) {
	type fields struct {
		creator    string
		reviewer   string
		auditor    string
		isAudited  bool
		isReviewed bool
	}
	type args struct {
		auditor string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "normal_success",
			fields: fields{
				creator:    "user-a",
				reviewer:   "",
				auditor:    "",
				isAudited:  false,
				isReviewed: false,
			},
			args: args{auditor: "user-b"},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				assert.NoError(t, err)
				return false
			},
		},
		{
			name: "noAuditor_error",
			fields: fields{
				creator:    "user-a",
				reviewer:   "",
				auditor:    "",
				isAudited:  false,
				isReviewed: false,
			},
			args: args{auditor: ""},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				assert.Equal(t, errAuditEmptyAuditor, err.(domainErr).Slug())
				return true
			},
		},
		{
			name: "repeatAudit_error",
			fields: fields{
				creator:    "user-a",
				reviewer:   "",
				auditor:    "user-b",
				isAudited:  true,
				isReviewed: false,
			},
			args: args{auditor: "user-b"},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				assert.Equal(t, errAuditRepeatAudit, err.(domainErr).Slug())
				return true
			},
		},
		{
			name: "auditorSameAsCreator_error",
			fields: fields{
				creator:    "user-a",
				reviewer:   "",
				auditor:    "",
				isAudited:  false,
				isReviewed: false,
			},
			args: args{auditor: "user-a"},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				assert.Equal(t, errAuditAuditorSameAsCreator, err.(domainErr).Slug())
				return true
			},
		},
		{
			name: "auditorSameAsReviewer_error",
			fields: fields{
				creator:    "user-a",
				reviewer:   "user-b",
				auditor:    "",
				isAudited:  false,
				isReviewed: true,
			},
			args: args{auditor: "user-b"},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				assert.Equal(t, errAuditAuditorSameAsReviewer, err.(domainErr).Slug())
				return true
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := &Voucher{
				creator:    tt.fields.creator,
				reviewer:   tt.fields.reviewer,
				auditor:    tt.fields.auditor,
				isAudited:  tt.fields.isAudited,
				isReviewed: tt.fields.isReviewed,
			}
			tt.wantErr(t, v.Audit(tt.args.auditor), fmt.Sprintf("Audit(%v)", tt.args.auditor))
		})
	}
}

func TestVoucher_CancelAudit(t *testing.T) {
	type fields struct {
		creator    string
		auditor    string
		reviewer   string
		isAudited  bool
		isReviewed bool
		isPosted   bool
	}
	type args struct {
		auditor string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "normal_success",
			fields: fields{
				creator:    "user-a",
				auditor:    "user-b",
				reviewer:   "",
				isAudited:  true,
				isReviewed: false,
				isPosted:   false,
			},
			args: args{auditor: "user-b"},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				assert.NoError(t, err)
				return false
			},
		},
		{
			name: "notAudited_error",
			fields: fields{
				creator:    "user-a",
				auditor:    "",
				reviewer:   "",
				isAudited:  false,
				isReviewed: false,
				isPosted:   false,
			},
			args: args{auditor: "user-b"},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				assert.Equal(t, errCancelAuditNotAudited, err.(domainErr).Slug())
				return true
			},
		},
		{
			name: "differentAuditor_error",
			fields: fields{
				creator:    "user-a",
				auditor:    "user-b",
				reviewer:   "",
				isAudited:  true,
				isReviewed: false,
				isPosted:   false,
			},
			args: args{auditor: "user-c"},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				assert.Equal(t, errCancelAuditDifferentAuditor, err.(domainErr).Slug())
				return true
			},
		},
		{
			name: "alreadyPosted_error",
			fields: fields{
				creator:    "user-a",
				auditor:    "user-b",
				reviewer:   "user-c",
				isAudited:  true,
				isReviewed: true,
				isPosted:   true,
			},
			args: args{auditor: "user-b"},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				assert.Equal(t, errCancelAuditPosted, err.(domainErr).Slug())
				return true
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := &Voucher{
				creator:    tt.fields.creator,
				auditor:    tt.fields.auditor,
				reviewer:   tt.fields.reviewer,
				isAudited:  tt.fields.isAudited,
				isReviewed: tt.fields.isReviewed,
				isPosted:   tt.fields.isPosted,
			}
			tt.wantErr(t, v.CancelAudit(tt.args.auditor), fmt.Sprintf("CancelAudit(%v)", tt.args.auditor))
		})
	}
}
