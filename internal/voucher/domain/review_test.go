package domain

import (
	"fmt"
	"testing"

	"github.com/google/uuid"

	"github.com/stretchr/testify/assert"
)

func TestVoucher_Review(t *testing.T) {
	type fields struct {
		creator    uuid.UUID
		auditor    uuid.UUID
		reviewer   uuid.UUID
		isAudited  bool
		isReviewed bool
		isPosted   bool
	}
	type args struct {
		reviewer uuid.UUID
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
				creator:    userA,
				auditor:    uuid.Nil,
				reviewer:   uuid.Nil,
				isAudited:  false,
				isReviewed: false,
				isPosted:   false,
			},
			args: args{reviewer: userB},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				assert.NoError(t, err)
				return false
			},
		},
		{
			name: "emptyReviewer_success",
			fields: fields{
				creator:    userA,
				auditor:    uuid.Nil,
				reviewer:   uuid.Nil,
				isAudited:  false,
				isReviewed: false,
				isPosted:   false,
			},
			args: args{reviewer: uuid.Nil},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				assert.Equal(t, errReviewEmptyReviewer, err.(domainErr).Slug())
				return true
			},
		},
		{
			name: "repeatReview_success",
			fields: fields{
				creator:    userA,
				auditor:    uuid.Nil,
				reviewer:   userB,
				isAudited:  false,
				isReviewed: true,
				isPosted:   false,
			},
			args: args{reviewer: userB},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				assert.Equal(t, errReviewRepeatReview, err.(domainErr).Slug())
				return true
			},
		},
		{
			name: "reviewerSameAsCreator_success",
			fields: fields{
				creator:    userA,
				auditor:    uuid.Nil,
				reviewer:   uuid.Nil,
				isAudited:  false,
				isReviewed: false,
				isPosted:   false,
			},
			args: args{reviewer: userA},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				assert.Equal(t, errReviewReviewerSameAsCreator, err.(domainErr).Slug())
				return true
			},
		},
		{
			name: "reviewerSameAsCreator_success",
			fields: fields{
				creator:    userA,
				auditor:    userB,
				reviewer:   uuid.Nil,
				isAudited:  true,
				isReviewed: false,
				isPosted:   false,
			},
			args: args{reviewer: userB},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				assert.Equal(t, errReviewReviewerSameAsAuditor, err.(domainErr).Slug())
				return true
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := &Voucher{
				creator:    tt.fields.creator,
				reviewer:   tt.fields.reviewer,
				isReviewed: tt.fields.isReviewed,
				auditor:    tt.fields.auditor,
				isAudited:  tt.fields.isAudited,
				isPosted:   tt.fields.isPosted,
			}
			tt.wantErr(t, v.Review(tt.args.reviewer), fmt.Sprintf("Review(%v)", tt.args.reviewer))
		})
	}
}

func TestVoucher_CancelReview(t *testing.T) {
	type fields struct {
		creator    uuid.UUID
		auditor    uuid.UUID
		reviewer   uuid.UUID
		isAudited  bool
		isReviewed bool
		isPosted   bool
	}
	type args struct {
		reviewer uuid.UUID
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
				creator:    userA,
				auditor:    uuid.Nil,
				reviewer:   userB,
				isAudited:  false,
				isReviewed: true,
				isPosted:   false,
			},
			args: args{reviewer: userB},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				assert.NoError(t, err)
				return false
			},
		},
		{
			name: "notReviewed_error",
			fields: fields{
				creator:    userA,
				auditor:    uuid.Nil,
				reviewer:   uuid.Nil,
				isAudited:  false,
				isReviewed: false,
				isPosted:   false,
			},
			args: args{reviewer: userB},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				assert.Equal(t, errCancelReviewNotReviewed, err.(domainErr).Slug())
				return true
			},
		},
		{
			name: "differentReviewer_error",
			fields: fields{
				creator:    userA,
				auditor:    uuid.Nil,
				reviewer:   userB,
				isAudited:  false,
				isReviewed: true,
				isPosted:   false,
			},
			args: args{reviewer: userC},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				assert.Equal(t, errCancelReviewDifferentReviewer, err.(domainErr).Slug())
				return true
			},
		},
		{
			name: "alreadyPosted_error",
			fields: fields{
				creator:    userA,
				auditor:    userB,
				reviewer:   userC,
				isAudited:  true,
				isReviewed: true,
				isPosted:   true,
			},
			args: args{reviewer: userC},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				assert.Equal(t, errCancelReviewPosted, err.(domainErr).Slug())
				return true
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := &Voucher{
				creator:    tt.fields.creator,
				reviewer:   tt.fields.reviewer,
				isReviewed: tt.fields.isReviewed,
				auditor:    tt.fields.auditor,
				isAudited:  tt.fields.isAudited,
				isPosted:   tt.fields.isPosted,
			}
			tt.wantErr(t, v.CancelReview(tt.args.reviewer), fmt.Sprintf("CancelReview(%v)", tt.args.reviewer))
		})
	}
}
