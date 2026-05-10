package errors

// This file is the single source of truth for all slug string constants used in
// the error framework.  Every NewSlugError / NewNotFoundError / NewConflictError
// call site must reference a constant from this file rather than an inline
// string literal.  Add new constants here whenever a new error slug is
// introduced, and remember to add the corresponding i18n entry in
// i18n/zh-CN.json.
//
// HTTP status intent is noted on each constant for quick reference.
// The actual status is controlled by which constructor is used (see slug_err.go).

// Common

const (
	SlugRecordNotFound           = "record-not-found"
	SlugDuplicateEntry           = "duplicate-entry"
	SlugForeignKeyViolation      = "foreign-key-violation"
	SlugCheckConstraintViolation = "check-constraint-violation"
	SlugNotNullViolation         = "not-null-violation"
)

// SoB

const (
	SlugSobEmptyId              = "sob-emptyId"
	SlugSobEmptyName            = "sob-emptyName"
	SlugSobNameTooLong          = "sob-nameTooLong"
	SlugSobDescriptionTooLong   = "sob-descriptionTooLong"
	SlugSobEmptyBaseCurrency    = "sob-emptyBaseCurrency"
	SlugSobInvalidStartingYear  = "sob-invalidStartingYear"
	SlugSobInvalidStartingMonth = "sob-invalidStartingMonth"
	SlugSobInvalidAccountLevel  = "sob-invalidAccountLevel"
	SlugSobCannotShortenLevel   = "sob-cannotShortenLevel"
	SlugSobDuplicateName        = "sob-duplicate-name"
)

// SoB / General

const (
	SlugEmptySobId = "emptySobId"
)

// Period

const (
	SlugPeriodEmptyId                  = "period-emptyId"
	SlugPeriodInvalidFiscalYear        = "period-invalidFiscalYear"
	SlugPeriodInvalidPeriodNumber      = "period-invalidPeriodNumber"
	SlugPeriodNotFound                 = "period-notFound"
	SlugPeriodDuplicateNumber          = "period-duplicate-number"
	SlugPeriodClosed                   = "period-closed"
	SlugPeriodCloseClosed              = "period-close-closed"
	SlugPeriodCloseIsNotCurrent        = "period-close-isNotCurrent"
	SlugPeriodCloseNotAllPosted        = "period-close-notAllJournalsPosted"
	SlugPeriodCloseUnclearedPnL        = "period-close-unclearedProfitAndLoss"
	SlugPeriodCloseUnclearedProfit     = "period-close-unclearedCurrentYearProfitAccount"
	SlugPeriodCloseOpeningUnequal      = "period-close-openingBalanceUnequal"
	SlugPeriodClosePeriodUnequal       = "period-close-periodBalanceUnequal"
	SlugPeriodCloseEndingUnequal       = "period-close-endingBalanceUnequal"
	SlugInvalidPeriodFormat            = "invalid-period-format"
	SlugPeriodRangeInvalid             = "period-range-invalid"
	SlugPeriodRangeNotContinuous       = "period-range-not-continuous"
	SlugPeriodBatchCloseTargetInPast   = "period-batchClose-targetInPast"
	SlugPeriodBatchCloseTooManyPeriods = "period-batchClose-tooManyPeriods"
)

// Account

const (
	SlugAccountNilId              = "account-nilId"
	SlugAccountNilSob             = "account-nilSob"
	SlugAccountNilSuperiorId      = "account-nilSuperiorId"
	SlugAccountEmptyTitle         = "account-emptyTitle"
	SlugAccountTitleTooLong       = "account-titleTooLong"
	SlugAccountEmptyRawNumber     = "account-emptyRawNumber"
	SlugAccountNotFound           = "account-notFound"
	SlugAccountDuplicateNumber    = "account-duplicate-number"
	SlugInvalidAccountClass       = "invalid-account-class"
	SlugInvalidAccountGroup       = "invalid-account-group"
	SlugAccountDeleteHasChildren  = "account-delete-hasChildren"
	SlugAccountDeleteUsedByJLine  = "account-delete-usedByJournalLine"
	SlugAccountDeleteHasOpBalance = "account-delete-hasOpeningBalance"

	SlugAccountClassMismatch      = "account-classMismatch"
	SlugAccountGroupMismatch      = "account-groupMismatch"
	SlugAccountLevelExceedsLimit  = "account-levelExceedsLimit"
	SlugAccountCodeLengthExceeded = "account-codeLengthExceeded"
)

// Ledger

const (
	SlugLedgerNilId        = "ledger-nilId"
	SlugLedgerNilSobId     = "ledger-nilSobId"
	SlugLedgerNilPeriodId  = "ledger-nilPeriodId"
	SlugLedgerNilAccountId = "ledger-nilAccountId"
	SlugLedgerNilAccount   = "ledger-nilAccount"

	SlugLedgerTransactionsMissingFilter = "ledger-transactions-missingFilter"
)

// Journal

const (
	SlugJournalEmptyId                  = "journal-emptyId"
	SlugJournalEmptyPeriodId            = "journal-emptyPeriodId"
	SlugJournalEmptyPeriod              = "journal-emptyPeriod"
	SlugJournalEmptyHeaderText          = "journal-emptyHeaderText"
	SlugJournalEmptyNumber              = "journal-emptyNumber"
	SlugJournalInvalidAttachmentQty     = "journal-invalidAttachmentQuantity"
	SlugJournalEmptyCreator             = "journal-emptyCreator"
	SlugJournalEmptyReviewer            = "journal-emptyReviewer"
	SlugJournalEmptyAuditor             = "journal-emptyAuditor"
	SlugJournalEmptyPoster              = "journal-emptyPoster"
	SlugJournalInvalidPostStatus        = "journal-invalidPostStatus"
	SlugJournalZeroTransactionDate      = "journal-zeroTransactionDate"
	SlugJournalEmptyJournalLines        = "journal-emptyJournalLines"
	SlugJournalNilJournalLine           = "journal-nilJournalLine"
	SlugJournalNotBalanced              = "journal-notBalanced"
	SlugJournalInvalidJournalType       = "journal-invalidJournalType"
	SlugJournalMissingReferenceId       = "journal-missingReferenceJournalId"
	SlugJournalUnexpectedReferenceId    = "journal-unexpectedReferenceJournalId"
	SlugJournalNotFound                 = "journal-notFound"
	SlugJournalReferenceNotFound        = "journal-referenceJournalNotFound"
	SlugJournalDuplicateDocumentNumber  = "journal-duplicate-document-number"
	SlugJournalUpdateAudited            = "journal-update-audited"
	SlugJournalUpdateReviewed           = "journal-update-reviewed"
	SlugJournalUpdateNotCreator         = "journal-update-notCreator"
	SlugJournalPostPeriodClosed         = "journal-post-periodClosed"
	SlugJournalPostPeriodNotCurrent     = "journal-post-periodNotCurrent"
	SlugJournalPostRepeatPost           = "journal-post-repeatPost"
	SlugJournalPostNotAudited           = "journal-post-notAudited"
	SlugJournalPostNotReviewed          = "journal-post-notReviewed"
	SlugJournalPostEmptyPoster          = "journal-post-emptyPoster"
	SlugJournalAuditRepeatAudit         = "journal-audit-repeatAudit"
	SlugJournalAuditEmptyAuditor        = "journal-audit-emptyAuditor"
	SlugJournalAuditSameAsCreator       = "journal-audit-auditorSameAsCreator"
	SlugJournalAuditSameAsReviewer      = "journal-audit-auditorSameAsReviewer"
	SlugJournalCancelAuditNotAudited    = "journal-cancelAudit-notAudited"
	SlugJournalCancelAuditDiffAuditor   = "journal-cancelAudit-differentAuditor"
	SlugJournalCancelAuditPosted        = "journal-cancelAudit-posted"
	SlugJournalReviewRepeat             = "journal-review-repeatReview"
	SlugJournalReviewEmptyReviewer      = "journal-review-emptyReviewer"
	SlugJournalReviewSameAsCreator      = "journal-review-reviewerSameAsCreator"
	SlugJournalReviewSameAsAuditor      = "journal-review-reviewerSameAsAuditor"
	SlugJournalCancelReviewNotReviewed  = "journal-cancelReview-notReviewed"
	SlugJournalCancelReviewDiffReviewer = "journal-cancelReview-differentReviewer"
	SlugJournalCancelReviewPosted       = "journal-cancelReview-posted"
	SlugJournalClosingAlreadyExists     = "journal-closing-alreadyExists"
	SlugJournalClosingUnpostedExist     = "journal-closing-unpostedJournalsExist"
	SlugJournalClosingNoBalance         = "journal-closing-noBalanceToClear"
	SlugJournalYearEndNotYearEnd        = "journal-yearEndClosing-notYearEndPeriod"
	SlugJournalYearEndAlreadyExists     = "journal-yearEndClosing-alreadyExists"
	SlugJournalYearEndPnLNotCleared     = "journal-yearEndClosing-pnlNotCleared"
	SlugJournalYearEndNoBalance         = "journal-yearEndClosing-noBalanceToClear"
	SlugJournalDeleteNotSystemJournal   = "journal-delete-notSystemJournal"
)

// Journal Line

const (
	SlugJournalLineEmptyId               = "journalLine-emptyId"
	SlugJournalLineNilAccount            = "journalLine-nilAccount"
	SlugJournalLineEmptyAccountId        = "journalLine-emptyAccountId"
	SlugJournalLineEmptyText             = "journalLine-emptyText"
	SlugJournalLineEmptyAmount           = "journalLine-emptyAmount"
	SlugJournalLineInvalidDimension      = "journalLine-invalidDimensionOption"
	SlugJournalLineDuplicateDimCategory  = "journalLine-duplicateDimensionCategory"
	SlugJournalLineDisallowedDimCategory = "journalLine-disallowedDimensionCategory"
	SlugJournalLineMissingDimCategory    = "journalLine-missingRequiredDimensionCategory"
)

// Dimension

const (
	SlugDimCategoryEmptyId       = "dimension-category-emptyId"
	SlugDimCategoryEmptySobId    = "dimension-category-emptySobId"
	SlugDimCategoryEmptyName     = "dimension-category-emptyName"
	SlugDimCategoryNameTooLong   = "dimension-category-nameTooLong"
	SlugDimCategoryDuplicateName = "dimension-category-duplicate-name"
	SlugDimCategoryDeleteHasUsed = "dimension-deleteCategory-hasUsedOptions"

	SlugDimOptionEmptyId         = "dimension-option-emptyId"
	SlugDimOptionEmptyCategoryId = "dimension-option-emptyCategoryId"
	SlugDimOptionEmptyName       = "dimension-option-emptyName"
	SlugDimOptionNameTooLong     = "dimension-option-nameTooLong"
	SlugDimOptionDuplicateName   = "dimension-option-duplicate-name"
	SlugDimOptionDeleteIsUsed    = "dimension-deleteOption-isUsed"
)

// Report

const (
	SlugReportEmptyId                  = "report-emptyId"
	SlugReportEmptySobId               = "report-emptySobId"
	SlugReportTemplateHasPeriod        = "report-templateHasPeriod"
	SlugReportEmptyPeriodId            = "report-emptyPeriodId"
	SlugReportEmptyTitle               = "report-emptyTitle"
	SlugReportEmptySections            = "report-emptySections"
	SlugReportEmptyAmountTypes         = "report-emptyAmountTypes"
	SlugReportInvalidAmountType        = "report-invalidAmountType"
	SlugReportCopyTemplate             = "report-copyTemplate"
	SlugReportDuplicateTitle           = "report-duplicate-title"
	SlugReportGenerateEmptyPeriod      = "report-generate-emptyPeriod"
	SlugReportValidationMissingSection = "report-validation-missingSectionType"
	SlugReportValidationMissingItem    = "report-validation-missingItemType"
	SlugReportBalanceSheetImbalance    = "report-balanceSheet-imbalance"
	SlugReportIncomeProfitMismatch     = "report-incomeStatement-profitMismatch"
)

const (
	SlugReportSectionNilId    = "report-section-nilId"
	SlugReportSectionZeroSeq  = "report-section-zeroSequence"
	SlugReportSectionNotFound = "report-section-notFound"
)

const (
	SlugReportItemNil                        = "report-item-nil"
	SlugReportItemNilId                      = "report-item-nilId"
	SlugReportItemEmptyText                  = "report-item-emptyText"
	SlugReportItemInvalidLevel               = "report-item-invalidLevel"
	SlugReportItemZeroSeq                    = "report-item-zeroSequence"
	SlugReportItemInvalidSumFactor           = "report-item-invalidSumFactor"
	SlugReportItemInvalidDataSourceWithForms = "report-item-invalidDataSourceWithFormulas"
	SlugReportItemRootIsBreakdown            = "report-item-rootLevelIsBreakdownItem"
	SlugReportItemBreakdownNoChild           = "report-item-breakdownItemCannotAddChild"
	SlugReportItemNotEditable                = "report-item-notEditable"
	SlugReportItemNotFound                   = "report-item-notFound"
	SlugReportItemTextRequired               = "report-item-textRequired"
	SlugReportItemLevelRequired              = "report-item-levelRequired"
	SlugReportItemSumFactorRequired          = "report-item-sumFactorRequired"
	SlugReportItemDataSourceRequired         = "report-item-dataSourceRequired"
)

const (
	SlugReportFormulaNilId            = "report-formula-nilId"
	SlugReportFormulaZeroSeq          = "report-formula-zeroSequence"
	SlugReportFormulaEmptyAccountId   = "report-formula-emptyAccountId"
	SlugReportFormulaInvalidSumFactor = "report-formula-invalidSumFactor"
)

// Numbering

const (
	SlugNumberingIdEmpty               = "numbering-id-empty"
	SlugNumberingConfigIdEmpty         = "numbering-configId-empty"
	SlugNumberingIdentifierEmpty       = "numbering-identifier-empty"
	SlugNumberingTargetObjectEmpty     = "numbering-targetObject-empty"
	SlugNumberingPropertyMatchersEmpty = "numbering-propertyMatchers-empty"
	SlugNumberingPropertyNameEmpty     = "numbering-propertyName-empty"
	SlugNumberingPropertyValueEmpty    = "numbering-propertyValue-empty"
	SlugNumberingConfigDuplicate       = "numbering-config-duplicate"
	SlugNumberingIdentifierDuplicate   = "numbering-identifier-duplicate"
)

// User

const (
	SlugUserEmptyId     = "user-emptyId"
	SlugUserEmptyTraits = "user-emptyTraits"
)
