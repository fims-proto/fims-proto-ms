package domain

import (
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"

	"github.com/stretchr/testify/assert"
)

func TestVoucher_UpdateLineItems(t *testing.T) {
	type fields struct {
		isReviewed bool
		isAudited  bool
	}
	type args struct {
		items []*LineItem
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
				isReviewed: false,
				isAudited:  false,
			},
			args: args{items: prepareBalancedItems()},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				assert.NoError(t, err)
				return false
			},
		},
		{
			name: "emptyLineItem_error",
			fields: fields{
				isReviewed: false,
				isAudited:  false,
			},
			args: args{items: nil},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				assert.Equal(t, errVoucherEmptyLineItem, err.(domainErr).Slug())
				return true
			},
		},
		{
			name: "notBalanced_error",
			fields: fields{
				isReviewed: false,
				isAudited:  false,
			},
			args: args{items: prepareImbalancedItems()},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				assert.Equal(t, errVoucherNotBalanced, err.(domainErr).Slug())
				return true
			},
		},
		{
			name: "voucherReviewed_error",
			fields: fields{
				isReviewed: true,
				isAudited:  false,
			},
			args: args{items: prepareBalancedItems()},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				assert.Equal(t, errUpdateReviewed, err.(domainErr).Slug())
				return true
			},
		},
		{
			name: "voucherAudited_error",
			fields: fields{
				isReviewed: false,
				isAudited:  true,
			},
			args: args{items: prepareBalancedItems()},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				assert.Equal(t, errUpdateAudited, err.(domainErr).Slug())
				return true
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := &Voucher{
				isReviewed: tt.fields.isReviewed,
				isAudited:  tt.fields.isAudited,
			}
			tt.wantErr(t, v.UpdateLineItems(tt.args.items), fmt.Sprintf("UpdateLineItems(%v)", tt.args.items))
		})
	}
}

func TestVoucher_UpdateTransactionTime(t *testing.T) {
	type fields struct {
		isReviewed bool
		isAudited  bool
	}
	type args struct {
		transactionTime time.Time
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
				isReviewed: false,
				isAudited:  false,
			},
			args: args{transactionTime: time.Now()},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				assert.NoError(t, err)
				return false
			},
		},
		{
			name: "emptyTransactionTime_success",
			fields: fields{
				isReviewed: false,
				isAudited:  false,
			},
			args: args{transactionTime: time.Time{}},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				assert.Equal(t, errUpdateZeroTransactionTime, err.(domainErr).Slug())
				return true
			},
		},
		{
			name: "futureTransactionTime_success",
			fields: fields{
				isReviewed: false,
				isAudited:  false,
			},
			args: args{transactionTime: time.Now().Add(time.Hour)},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				assert.Equal(t, errVoucherFutureTransactionTime, err.(domainErr).Slug())
				return true
			},
		},
		{
			name: "voucherReviewed_success",
			fields: fields{
				isReviewed: true,
				isAudited:  false,
			},
			args: args{transactionTime: time.Now()},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				assert.Equal(t, errUpdateReviewed, err.(domainErr).Slug())
				return true
			},
		},
		{
			name: "voucherAudited_success",
			fields: fields{
				isReviewed: false,
				isAudited:  true,
			},
			args: args{transactionTime: time.Now()},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				assert.Equal(t, errUpdateAudited, err.(domainErr).Slug())
				return true
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := &Voucher{
				isReviewed: tt.fields.isReviewed,
				isAudited:  tt.fields.isAudited,
			}
			tt.wantErr(t, v.UpdateTransactionTime(tt.args.transactionTime, uuid.New()), fmt.Sprintf("UpdateTransactionTime(%v)", tt.args.transactionTime))
		})
	}
}
