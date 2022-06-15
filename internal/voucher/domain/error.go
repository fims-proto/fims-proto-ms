package domain

import "fmt"

type domainErr struct {
	slug string
	args []any
}

func (e domainErr) Error() string {
	return fmt.Sprintf("Slug Error: %s.", e.slug)
}

func (e domainErr) Slug() string {
	return e.slug
}

func newDomainErr(slug string, args ...any) domainErr {
	return domainErr{
		slug: slug,
		args: args,
	}
}

const (
	errVoucherEmptyId             = "voucher-emptyId"
	errVoucherEmptySobId          = "voucher-emptySobId"
	errVoucherEmptyNumber         = "voucher-emptyNumber"
	errVoucherEmptyCreator        = "voucher-emptyCreator"
	errVoucherEmptyReviewer       = "voucher-emptyReviewer"
	errVoucherEmptyAuditor        = "voucher-emptyAuditor"
	errVoucherInvalidPostStatus   = "voucher-invalidPostStatus"
	errVoucherEmptyLineItem       = "voucher-emptyLineItem"
	errVoucherNotBalanced         = "voucher-notBalanced"
	errVoucherZeroTransactionTime = "voucher-zeroTransactionTime"
	errVoucherTypeNotSupported    = "voucher-type-notSupported"

	errLineItemEmptyId            = "voucher-lineItem-emptyId"
	errLineItemEmptyAccountId     = "voucher-lineItem-emptyAccountId"
	errLineItemEmptySummary       = "voucher-lineItem-emptySummary"
	errLineItemEmptyDebitCredit   = "voucher-lineItem-emptyDebitCredit"
	errLineItemDebitCreditCoExist = "voucher-lineItem-debitCreditCoExist"

	errUpdateAudited             = "voucher-update-audited"
	errUpdateReviewed            = "voucher-update-reviewed"
	errUpdateZeroTransactionTime = "voucher-update-zeroTransactionTime"

	errAuditEmptyAuditor           = "voucher-audit-emptyAuditor"
	errAuditRepeatAudit            = "voucher-audit-repeatAudit"
	errAuditAuditorSameAsCreator   = "voucher-audit-auditorSameAsCreator"
	errAuditAuditorSameAsReviewer  = "voucher-audit-auditorSameAsReviewer"
	errCancelAuditNotAudited       = "voucher-cancelAudit-notAudited"
	errCancelAuditDifferentAuditor = "voucher-cancelAudit-differentCancelAuditor"
	errCancelAuditPosted           = "voucher-cancelAudit-posted"

	errReviewEmptyReviewer           = "voucher-review-emptyReviewer"
	errReviewRepeatReview            = "voucher-review-repeatReview"
	errReviewReviewerSameAsCreator   = "voucher-review-reviewerSameAsCreator"
	errReviewReviewerSameAsAuditor   = "voucher-review-reviewerSameAsAuditor"
	errCancelReviewNotReviewed       = "voucher-cancelReview-notReviewed"
	errCancelReviewDifferentReviewer = "voucher-cancelReview-differentCancelReviewer"
	errCancelReviewPosted            = "voucher-cancelReview-posted"

	errPostRepeatPost  = "voucher-post-repeatPost"
	errPostNotAudited  = "voucher-post-notAudited"
	errPostNotReviewed = "voucher-post-notReviewed"
)
