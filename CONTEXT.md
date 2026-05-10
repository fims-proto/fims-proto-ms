# CONTEXT.md

This file captures the bounded context glossary for the FIMS General Ledger domain. Terms here are meaningful to domain experts and should be used consistently across code, tests, and documentation.

## Glossary

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
