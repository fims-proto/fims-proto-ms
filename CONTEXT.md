# CONTEXT.md

This file captures the bounded context glossary for FIMS. Terms here are meaningful to domain experts and should be used consistently across code, tests, and documentation.

For a full business-facing overview of the system, see [docs/system-overview.md](docs/system-overview.md).

## System-Level Terms

### FIMS (Financial Information Management System / 财务信息管理系统)

A standalone accounting SaaS product for Chinese SMEs. Each customer receives an isolated deployment (not shared-tenant). Implements **小企业会计准则 (Small Business Accounting Standards)**; multi-standard support is on the roadmap.

### Set of Books (SoB / 账套)

The top-level accounting entity within a FIMS deployment. One deployment hosts multiple SoBs — enabling management of multiple subsidiary companies, legal entities, or (for accounting firms) multiple client companies. Each SoB is fully independent: its own chart of accounts, accounting periods, journals, and ledger balances.

### Small Business Accounting Standards (小企业会计准则)

The Chinese accounting standard that FIMS currently implements. Pre-loaded as seed data in `dataload/xqykjzz/`. Support for additional standards (e.g., 企业会计准则) is planned.

---

## General Ledger Glossary

### Period Closing (结账)

The act of sealing an accounting period so that no further journal entries can be posted to it. Closing a period:
1. Validates all journals in the period are posted
2. Validates P&L accounts have zero ending balance (cleared by Monthly Closing Journal)
3. Validates trial balance (sum of all signed amounts equals zero)
4. Validates the Current Year Profit account is zero for period 12 (cleared by Year-End Closing Journal)
5. Marks the period as closed and opens the next period

### Monthly Closing Journal (月末结账凭证)

A system-generated journal that reverses all leaf P&L account balances to zero and transfers the net result to the Current Year Profit account (003103). Generated automatically at month-end before closing the period. Skipped if there are no P&L balances.

### Year-End Closing Journal (年末结账凭证)

A system-generated journal that transfers the Current Year Profit account (003103) balance to Retained Earnings (003104000002). Only applicable in period 12. Skipped if the Current Year Profit account has zero balance.

### Continuous Period Closing (连续结账)

A batch operation that closes a sequence of accounting periods from the current period to a user-specified target period in a single atomic transaction. For each period, the system automatically creates the Monthly Closing Journal (and Year-End Closing Journal if period 12) before closing the period. The entire batch rolls back if any period fails validation. Maximum 12 periods per batch.

### Trial Balance (试算平衡)

A validation that confirms the sum of all signed amounts across level-1 accounts equals zero (opening, period, and ending balances). Passes when the books are in balance. Used as a precondition for period closing.

### Current Year Profit (本年利润)

Account number 003103. Accumulates the net P&L result across all months of the fiscal year via Monthly Closing Journals. Must be transferred to Retained Earnings (003104000002) via the Year-End Closing Journal before period 12 can be closed.

---

## Reports Glossary

### Report Template (报表模板)

Per-SoB structural configuration for a mandatory report type. Holds sections, items, and formula rules. No period association. Two templates exist per SoB after initialization: Balance Sheet (资产负债表) and Income Statement (利润表). Updated via `PATCH /sob/{sobId}/report/{reportId}`. Changes to a template do not cascade to existing instances — instances must be manually regenerated.

### Report Instance (报表实例)

A generated financial report for a specific accounting period. Exactly one instance exists per (SoB, class, period) — enforced by a DB partial unique index. Generation is idempotent: calling generate for a period that already has an instance regenerates it (recalculates amounts from current ledgers) rather than creating a duplicate. Automatically generated for all classes when a period is closed via `ClosePeriodHandler`. Accessible by natural key via `GET /sob/{sobId}/report/{class}/{period}` (YYYY-MM format); returns 404 if no instance has been generated for that period yet.

---

## Cash Flow Glossary

### Cash Equivalent Account (现金等价物科目)

A GL account flagged as `is_cash_equivalent = true`, representing holdings of cash or near-cash instruments. By default after SoB initialization: 库存现金 (1001), 银行存款 (1002), and 其他货币资金 (1012, including all child accounts). User-editable per SoB. The flag is stored per-account with no tree inheritance — child accounts of a cash-equivalent account do not inherit the flag automatically.

### Cash Flow Item (现金流量项目)

A predefined classification category for cash movements, defined by the accounting standard. Each item belongs to one of three categories: OPERATING (经营活动), INVESTING (投资活动), or FINANCING (筹资活动). Each item has a direction: INFLOW (现金流入) or OUTFLOW (现金流出). Items are per-SoB (seeded at SoB creation from the accounting standard's catalog — currently always xqykjzz). 小企业会计准则 defines 16 items: OP_01..OP_06, INV_01..INV_05, FIN_01..FIN_05.

### Cash Flow Classification (现金流量归类)

The act of tagging a non-cash-equivalent journal line with a Cash Flow Item ID. Required whenever a journal entry contains at least one cash-equivalent line and at least one non-cash-equivalent line. Omitted when all lines are cash-equivalent (internal transfer) or when no cash-equivalent line is present (accrual entry). Validated at journal creation and update time.

### Internal Cash Transfer (内部现金划转)

A journal entry where every line touches a cash-equivalent account (e.g., moving money from 银行存款 to 库存现金). No Cash Flow Item tagging is required — the entry does not represent a net cash inflow or outflow to the entity.

### Cash Flow Statement (现金流量表)

The third mandatory financial report (alongside Balance Sheet and Income Statement) per 小企业会计准则. Structured into three activity sections plus an opening/closing cash balance reconciliation. Amounts are aggregated from posted journal lines tagged with Cash Flow Items, not from ledger balances. Two amount columns: 本月金额 (current period) and 本年累计金额 (year-to-date from period 1 through current period). Auto-generated at period close alongside BS and IS.
