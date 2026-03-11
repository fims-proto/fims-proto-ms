# FIMS Business Design Document

**Financial Information Management System - Solution Description**

_Version: 1.1_
_Date: 2026-03-10_

**Changelog v1.1**:

- Added LedgerEntry entity description
- Updated TransactionDate type description
- Corrected periodNumber range (supports custom fiscal periods)
- Added DataSource.None value
- Updated user traits storage description (JSONB format)
- Added SectionType and ItemType enums
- Corrected minimum journal lines validation (1 line minimum)

---

## Table of Contents

1. [System Overview](#1-system-overview)
2. [Core Business Domains](#2-core-business-domains)
3. [Multi-Entity Model](#3-multi-entity-model)
4. [General Ledger Management](#4-general-ledger-management)
5. [Journal Lifecycle & Workflow](#5-journal-lifecycle--workflow)
6. [Period Management & Closing](#6-period-management--closing)
7. [Financial Reporting](#7-financial-reporting)
8. [Numbering & Identifier Management](#8-numbering--identifier-management)
9. [User Management](#9-user-management)
10. [Key Business Rules & Constraints](#10-key-business-rules--constraints)

---

## 1. System Overview

### 1.1 Purpose

FIMS (Financial Information Management System) is an **accounting system** designed to provide comprehensive financial management capabilities for multiple independent accounting entities. The system follows double-entry bookkeeping principles and supports the complete accounting cycle from transaction recording to financial reporting.

### 1.2 Deployment Model

- **Single-Tenant Deployment**: One system instance per customer (tenant)
- **Multiple Accounting Entities**: Within one deployment, supports multiple independent Set of Books (SoB/账套)
- **Customer = Tenant**: Each customer gets their own isolated deployment instance
- **SoB = Accounting Entity**: Within the customer's deployment, they can create multiple SoBs for different legal entities, subsidiaries, or divisions

### 1.3 Architecture Philosophy

- **Hexagonal Architecture (Ports & Adapters)**: Clear separation between business logic and infrastructure
- **CQRS Pattern**: Commands (writes) and Queries (reads) are separated for scalability
- **Domain-Driven Design**: Rich domain models with business rules enforced at the domain level
- **Multi-Entity Support**: Complete data isolation between Set of Books (SoB) within same deployment

### 1.4 Core Principles

- **Data Integrity**: All financial transactions are validated and balanced
- **Audit Trail**: Complete tracking of who created, reviewed, audited, and posted transactions
- **Segregation of Duties**: Different users must perform creation, review, audit, and posting
- **Period Control**: Transactions can only be posted to open periods

---

## 2. Core Business Domains

FIMS consists of five primary business domains:

### 2.1 Set of Books (SoB) - 账套

**Business Purpose**: Represents a complete, independent accounting entity.

**Key Attributes**:

- **Name & Description**: Identifies the accounting entity
- **Base Currency**: The functional currency for all transactions
- **Starting Period**: The initial accounting period (year and month)
- **Account Code Structure**: Configurable hierarchical account numbering (2-10 levels, each level 1-6 digits)

**Business Significance**:

- Each SoB is completely isolated from other SoBs within the same deployment
- A customer may have multiple SoBs (e.g., separate legal entities, subsidiaries, divisions)
- All accounting data (accounts, periods, journals, reports) belongs to a specific SoB

### 2.2 General Ledger - 总账

**Business Purpose**: Core accounting engine that manages the chart of accounts, accounting periods, ledger balances, and journal entries.

**Sub-Components**:

#### 2.2.1 Chart of Accounts

- **Based on Chinese Accounting Standards**: First-level accounts (一级科目) follow standard CoA structure
  - Currently supports: **小企业会计准则** (Accounting Standards for Small Enterprises)
  - Pre-configured first-level accounts (e.g., 1001-库存现金, 1002-银行存款, 1122-应收账款)
- **Hierarchical Structure**: Accounts organized in a tree (parent-child relationships)
  - Users add custom detail accounts (明细科目) under standard first-level accounts
- **Account Classes**: Assets (资产), Liabilities (负债), Equities (权益), Costs (成本), Profits & Losses (损益), Commons (共同)
- **Account Groups**: Sub-classification within each class (e.g., Current Assets, Fixed Assets)
- **Balance Direction**: Debit or Credit (determines natural balance)
- **Parent vs. Detail Accounts**:
  - **Parent Accounts (一级科目/上级科目)**: Summary accounts that aggregate child balances
  - **Detail/Leaf Accounts (明细科目)**: Lowest-level accounts that accept transactions
  - Only detail accounts can have transactions; parent accounts show rolled-up totals
- **Account Numbering**: Configurable digit length per level, zero-padded as prefix (e.g., 1001001 = 1001 + 001)
#### 2.2.2 Accounting Periods

- **Monthly Periods**: Accounting periods aligned to calendar months
- **Fiscal Year Support**: Periods span from starting month with flexible period numbering
- **Period Sequence**:
  - **periodNumber**: Accounting period sequence number (supports custom fiscal periods, not limited to 1-12)
  - **fiscalYear**: Accounting year (1970-9999)
  - Automatic calculation of next/previous periods with year rollover support
- **Period States**:
  - **Open**: Current period accepting new transactions
  - **Closed**: Historical period, no longer accepting transactions
  - **Current Flag**: Only one period can be current at a time

#### 2.2.3 Ledger (General Ledger)

- **One Ledger per Account per Period**: Tracks account balance movement
- **Scope**: Covers both parent accounts (一级科目) and detail accounts (明细科目)
- **Balance Components**:
  - **Opening Amount**: Signed decimal value (positive for debit accounts, negative for credit accounts)
  - **Period Amount**: Signed net movement during the period (positive = net debit, negative = net credit)
  - **Period Debit/Credit**: Positive values representing total debit and credit movements (kept for query performance)
  - **Ending Amount**: Calculated as openingAmount + periodAmount
- **Hierarchical Posting**: When journals post, both detail accounts and all parent accounts update

#### 2.2.4 Ledger Entry (账簿分录)


- **Purpose**: Stores detailed transaction history when journals are posted
- **Creation**: One ledger entry is created for each journal line during posting
- **Components**:
  - **Journal Reference**: Links to both journal and journal line IDs for traceability
  - **Account Reference**: The affected account ID
  - **Transaction Date**: Business date (year, month, day) of the transaction
  - **Amount**: Signed amount of the transaction
- **Usage**: Used by ledger explorer to display detailed transaction history for each account
- **Distinction**: Unlike Ledger (which stores aggregated balances), LedgerEntry stores individual transaction details

#### 2.2.5 Journal (分录)

See detailed section below (Section 5)

### 2.3 Report - 报表

**Business Purpose**: Generate standard financial statements and custom reports.

**Report Types**:

- **Balance Sheet (资产负债表)**: Statement of financial position
- **Income Statement (利润表)**: Profit and loss statement

**Report Components**:

- **Template vs. Instance**: Templates define structure; instances contain actual data for specific periods
- **Sections**: Logical groupings (Assets, Liabilities, Revenue, Expenses, etc.)
- **Items**: Individual journal lines with display order
- **Data Sources**:
  - **Sum**: Aggregate ledger balances by account filters (class, group, specific accounts)
  - **Formula**: Calculate from other items using rules (Net, Debit, Credit, Transaction)
- **Amount Types**:
  - Balance Sheet: Year Opening Balance, Period Ending Balance
  - Income Statement: Last Year Amount, Year-to-Date Amount, Period Amount

**Generation Process**:

1. Copy template to create report instance for specific period
2. Aggregate ledger data based on data source definitions
3. Calculate formulas referencing other items
4. Validate results (e.g., Balance Sheet must balance: Assets = Liabilities + Equity)

### 2.4 Numbering - 编号

**Business Purpose**: Automatically generate unique, sequential identifiers for business documents.

**Key Features**:

- **Configurable Patterns**: Prefix + Auto-incrementing Counter + Suffix
- **Context-Aware**: Different numbering sequences based on business object properties
- **Property Matchers**: Configure different sequences for different journal types, periods, etc.
- **Example**: Journal numbering might be "JV-202401-0001", "JV-202401-0002" for January 2024

**Use Cases**:

- Journal document numbers
- Account codes (automated sequential assignment within each level)
- Report identifiers

### 2.5 User - 用户

**Business Purpose**: Manage system users and their identities.

**Integration**:

- **Authentication**: Delegated to Ory Kratos (external identity provider)
- **User Records**: FIMS maintains minimal user information (ID, name, email)
- **Authorization**: Users referenced in journals (creator, reviewer, auditor, poster)

---

## 3. Multi-Entity Model

### 3.1 Deployment Model

- **Single Customer per Deployment**: Each customer (tenant) receives their own isolated FIMS deployment instance
- **No Multi-Tenancy**: There is NO sharing of infrastructure or data between different customers
- **Dedicated Resources**: Each deployment is completely independent

### 3.2 Set of Books (Multiple Accounting Entities)

- **Multiple SoBs per Deployment**: Within one customer's deployment, they can create multiple Set of Books
- **SoB Isolation**: Each Set of Books represents a completely isolated accounting entity
- **Data Segregation**: All domain entities (accounts, periods, journals, reports) belong to exactly one SoB
- **Shared Nothing Between SoBs**: No data sharing between different SoBs within the same deployment

### 3.3 Use Cases for Multiple SoBs

A single customer may create multiple SoBs for:

- **Multiple Legal Entities**: Parent company and subsidiaries
- **Different Divisions**: Separate accounting for different business units
- **Different Currencies**: Separate books for operations in different countries
- **Testing/Staging**: Production SoB vs. test SoB

### 3.4 Technical Implementation (Current)

- **Single Database**: One database per deployment with SoB ID filtering
- **DataSource Abstraction**: System designed to support future architectural changes
- **No Subdomain Routing**: All SoBs accessed through same application URL

### 3.5 User Access

- Users within a deployment can access multiple SoBs
- User permissions and roles managed per SoB (though current implementation is minimal)

---

## 4. General Ledger Management

### 4.1 Chart of Accounts Structure

**Standard CoA Foundation**:

- **Based on Chinese Accounting Standards (会计准则)**: System currently supports **小企业会计准则** (Accounting Standards for Small Enterprises)
- **Pre-configured First Level**: First-level accounts (一级科目) follow the standard CoA structure
- **User Customization**: Users typically:
  1. Start with standard first-level accounts (e.g., 1001-库存现金, 1002-银行存款)
  2. Add their own detail accounts (明细科目) underneath as needed
  3. Customize based on business requirements

**Hierarchical Design Example**:

```
1001 - 库存现金 (Cash on Hand) - 一级科目 from standard CoA
  ├─ 1001001 - 人民币 (RMB) - 明细科目 (Detail Account - user-defined)
  ├─ 1001002 - 美元 (USD) - 明细科目 (Detail Account - user-defined)
  └─ 1001003 - 欧元 (EUR) - 明细科目 (Detail Account - user-defined)

1002 - 银行存款 (Bank Deposits) - 一级科目 from standard CoA
  ├─ 1002001 - 中国银行-账户A - 明细科目 (Detail Account - user-defined)
  └─ 1002002 - 工商银行-账户B - 明细科目 (Detail Account - user-defined)

1122 - 应收账款 (Accounts Receivable) - 一级科目 from standard CoA
```

**Account Numbering Rules**:

- **Configurable Length per Level**: SoB defines code length at each level (e.g., [4, 3, 2] = level 1: 4 digits, level 2: 3 digits, level 3: 2 digits)
- **Zero-Padding as Prefix**: Numbers are left-padded with zeros
  - Example: Level 2 defined as 3 digits
  - First detail account under 1001 (一级科目): **1001** + **001** = **1001001**
  - Second detail account: **1001** + **002** = **1001002**
- **Number Hierarchy**: Internal representation [1001, 1] for account "1001001"
- **Concatenation**: Account number is concatenation of all level numbers with zero-padding
- **Unique within SoB**: Account numbers must be unique per SoB

**Account Classes & Groups**:

```
| Class ID | Class Name       | Possible Groups                  |
| -------- | ---------------- | -------------------------------- |
| 1        | Assets           | Current, Fixed, Intangible, etc. |
| 2        | Liabilities      | Current, Long-term               |
| 3        | Equities         | Share Capital, Retained Earnings |
| 4        | Costs            | Production, Manufacturing        |
| 5        | Profits & Losses | Revenue, Expenses                |
| 7        | Commons          | Special clearing accounts        |
```

**Balance Direction**:

- **Debit Accounts**: Assets, Costs, Expenses (normal debit balance)
- **Credit Accounts**: Liabilities, Equities, Revenue (normal credit balance)
- **Validation**: Used in report validation and trial balance checks

### 4.2 Ledger Posting Mechanics

**Posting Process** (when journal is posted):

1. **Journal Line Processing**: For each journal journal line
   - Update detail account ledger (明细科目) - add debit/credit to period activity
   - Update all parent account ledgers (上级科目) - hierarchical rollup

2. **Hierarchical Rollup**: Parent account balances aggregate all child activity
   - Example: Posting to "1001001 - Some Cash Account" (detail account) also updates "1001 - 库存现金" (top-level parent)

3. **Balance Calculation**:

   ```
   Ending Amount = Opening Amount + Period Amount
   ```

   Where:
   - Positive amounts represent debits
   - Negative amounts represent credits
   - `PeriodDebit` and `PeriodCredit` are maintained separately for query performance (always positive values)

4. **Batch Updates**: Multiple journal lines to same account are merged before posting (performance optimization)

---

## 5. Journal Lifecycle & Workflow

### 5.1 Journal Structure

**Header Information**:

- **Document Number**: Unique identifier (auto-generated by numbering service)
- **Journal Type**: Classification (currently only "General Journal" supported)
- **Period**: The accounting period this journal belongs to
- **Transaction Date**: Business date of the transaction (year, month, day without timezone)
  - Stored as custom type with separate Year, Month, Day fields
  - Automatically validates date validity (e.g., rejects Feb 30)
  - Serialized as ISO 8601 format "YYYY-MM-DD"
- **Header Text**: Description of the transaction
- **Attachment Quantity**: Number of supporting documents
- **Amount**: Transaction amount (sum of all positive/debit journal line amounts)

**Journal Lines**:

- **Account**: The general ledger account being debited or credited
- **Amount**: Signed decimal value (positive = debit, negative = credit)
- **Line Text**: Description of this specific line

**Business Rules**:

- **Balanced Entry**: Sum of all signed journal line amounts must equal zero (trial balance)
- **Minimum Lines**: At least 1 journal line (system validates non-empty list, double-entry convention typically requires 2+)
- **Non-zero Amounts**: Each line must have a non-zero signed amount
- **Detail Accounts Only**: Can only post to detail accounts/明细科目 (not parent accounts/上级科目)

### 5.2 Workflow States

**State Progression**:

```
Draft → Reviewed → Audited → Posted
          ↓          ↓
        Cancel     Cancel
        Review     Audit
```

**State Definitions**:

1. **Draft** (Created):
   - Initial state after journal creation
   - Can be edited by creator
   - Not yet validated by others

2. **Reviewed** (复核):
   - First level of approval
   - Performed by a different user than creator (segregation of duties)
   - Confirms transaction is valid and properly documented
   - Can be canceled back to draft

3. **Audited** (审核):
   - Second level of approval
   - Performed by a different user than creator and reviewer
   - Final validation before posting
   - Can be canceled back to reviewed state

4. **Posted** (登账/过账):
   - Final state - journal affects ledger balances
   - Cannot be modified or deleted
   - Updates general ledger balances
   - Cannot be reversed (new journal needed for corrections)

### 5.3 Segregation of Duties

**Mandatory Separation**:

- **Creator ≠ Reviewer ≠ Auditor ≠ Poster**: All four roles must be performed by different users
- **Purpose**: Prevent fraud and errors through independent verification
- **System Enforcement**: Domain rules validate user IDs at each state transition

**Role Descriptions**:

- **Creator**: Person who records the transaction
- **Reviewer (复核人)**: Person who verifies accuracy and completeness
- **Auditor (审核人)**: Person who provides final approval before posting
- **Poster (过账人)**: Person who commits the transaction to ledgers

### 5.4 Posting Requirements

**Pre-Posting Validations**:

- ✓ Journal must be both reviewed AND audited
- ✓ Period must be open (not closed)
- ✓ Period must be the current period
- ✓ Debit and credit totals must balance
- ✓ Poster must be different from creator, reviewer, and auditor

**Posting Effects**:

- Updates ledger balances for all affected accounts (including parent accounts)
- Sets journal state to "posted"
- Records poster user ID and posting timestamp

**Post-Posting Restrictions**:

- Journal cannot be edited
- Journal cannot be deleted
- State cannot be reversed
- Corrections require creating a reversing journal

---

## 6. Period Management & Closing

### 6.1 Period Lifecycle

**Period States**:

- **Open & Not Current**: Historical periods already closed
- **Open & Current**: Active period accepting new journals
- **Closed & Not Current**: Period is locked, no further changes allowed

**Period Transitions**:

```
Period Creation → Start (become current) → Close → Next Period Start
```

### 6.2 Period Initialization

**When a SoB is created**:

1. Create starting period (based on SoB's starting year and month)
2. Mark it as current period
3. Create initial ledgers for all accounts (zero balances)

**When a period closes**:

1. Close current period
2. Create next period (automatically increment month/year)
3. Mark next period as current
4. Initialize ledgers for next period (carry forward balances)

### 6.3 Period Closing Process

**Pre-Closing Validations**:

1. **All Journals Posted**: No unposted journals remain in the period
2. **Trial Balance**: Sum of all signed ledger amounts across all accounts equals zero
3. **Profit & Loss Cleared**: All P&L accounts (revenue/expense) have zero ending amount
   - Requires creating closing entries to transfer P&L to Retained Earnings

**Closing Steps**:

1. Validate all closing requirements
2. Set current period as closed and not current
3. Create next period (if doesn't exist)
4. Set next period as current
5. Initialize ledgers for next period:
   - **Balance Sheet Accounts**: Carry forward ending balance as opening balance
   - **P&L Accounts**: Start with zero balance (already cleared to equity)

**Post-Closing State**:

- Closed period cannot accept new journals
- Closed period cannot be modified
- Next period is now active

---

## 7. Financial Reporting

### 7.1 Report Design Philosophy

**Template-Instance Pattern**:

- **Templates**: Define report structure, sections, items, and calculation logic
- **Instances**: Generated from templates for specific periods with actual data

**Report Classes**:

- **Balance Sheet (资产负债表)**: Point-in-time financial position
- **Income Statement (利润表)**: Period performance (revenue and expenses)

### 7.2 Report Structure

**Three-Level Hierarchy**:

```
Report
  └─ Sections (e.g., Assets, Liabilities, Revenue, Expenses)
      └─ Items (individual journal lines with amounts)
```

**Section Types**:

| SectionType   | Chinese | Description         |
| ------------- | ------- | ------------------- |
| `assets`      | 资产    | Assets section      |
| `liabilities` | 负债    | Liabilities section |
| `equity`      | 权益    | Equity section      |
| `revenue`     | 收入    | Revenue section     |
| `expenses`    | 费用    | Expenses section    |

**Item Types**:

| ItemType           | Chinese  | Description                |
| ------------------ | -------- | -------------------------- |
| `gross_profit`     | 毛利     | Gross profit line item     |
| `operating_profit` | 营业利润 | Operating profit line item |
| `total_profit`     | 利润总额 | Total profit line item     |
| `net_profit`       | 净利润   | Net profit line item       |

**Item Components**:

- **Type**: Header, Detail, or Subtotal (or specific ItemType for calculated profit items)
- **Display Order**: Controls visual sequence
- **Bold/Indent**: Formatting hints
- **Data Source**: Where the data comes from (Sum, Formulas, or None)
- **Amounts**: Multiple columns based on amount types

### 7.3 Data Sources

Report items can have three types of data sources:

**1. Sum Data Source** - Sum of previous items:

- **Purpose**: Calculate subtotals by summing amounts from all previous items in the same section
- **Calculation**: Simply adds up amounts from all items processed so far
- **Use Case**: Subtotal rows, section totals
- **Example**:
  ```
  Item 1: Current Assets (Formulas) = $10,000
  Item 2: Fixed Assets (Formulas) = $5,000
  Item 3: Total Assets (Sum) = $10,000 + $5,000 = $15,000
  ```

**2. Formulas Data Source** - Aggregate ledger balances:

- **Purpose**: Calculate amounts from specified general ledger accounts
- **Components**:
  - **Account Filter**: Which account(s) to aggregate (by account ID)
  - **Formula Rule**: Which balance/value to use (see below)
  - **Sum Factor**: +1 (add) or -1 (subtract) for combining multiple formulas
- **Multiple Formulas**: One item can have multiple formulas (e.g., Account A + Account B - Account C)

**Formula Rules** (determines which value to extract from ledger):

| Rule            | Chinese  | Description                                                         | Use Case                                              |
| --------------- | -------- | ------------------------------------------------------------------- | ----------------------------------------------------- |
| **Net**         | 净值     | Net balance (debit - credit, adjusted by account balance direction) | Most common - gets the "normal" balance of an account |
| **Debit**       | 借方余额 | Debit balance only                                                  | When you need specific debit balance                  |
| **Credit**      | 贷方余额 | Credit balance only                                                 | When you need specific credit balance                 |
| **Transaction** | 发生额   | Period activity (debit movements - credit movements)                | Income statement items showing period activity        |

**3. None Data Source** - No automatic calculation:

- **Purpose**: Items with manually entered or placeholder values
- **Use Case**: Header rows, spacing items, or items requiring manual input
- **Calculation**: No automatic value calculation; amounts are manually entered or left blank

**Item Structure Fields**:

- **Sum Factor** (on Item level): 0, +1, or -1
  - Controls whether this item's amount is added to section total
  - 0 = not included in parent sum (e.g., header rows)
  - +1 = add to parent sum (normal items)
  - -1 = subtract from parent sum (e.g., contra-accounts)
- **Display Sum Factor**: Whether to show +/- sign in the report

### 7.4 Amount Types

**Balance Sheet Amount Types**:

- **Year Opening Balance**: Balance at start of fiscal year
- **Period Ending Balance**: Balance at end of reporting period

**Income Statement Amount Types**:

- **Last Year Amount**: Same period last fiscal year
- **Year-to-Date Amount**: Cumulative from start of fiscal year to reporting period
- **Period Amount**: Activity within reporting period only

### 7.5 Report Generation Process

**Template Instantiation**:

1. User selects report template
2. Specifies target period and desired amount types
3. System copies template structure to new report instance
4. Report instance is linked to specific period

**Data Population** (bottom-up calculation):

1. **Load Periods**: Retrieve current period and related historical periods for comparison
2. **Process Each Item**:
   - **For "Formulas" items**:
     - For each formula, aggregate ledger balances for the specified account
     - Apply formula rule (Net/Debit/Credit/Transaction) to extract the value
     - Multiply by sum factor (+1 or -1)
     - Sum all formulas in the item
   - **For "Sum" items**:
     - Use the accumulated sum of all previous items in the same section
3. **Calculate Section Totals**:
   - Sum all items in section (respecting each item's sum factor)
   - Roll up subsection amounts
4. **Store Results**: Persist calculated amounts in report instance

**Example Calculation Flow**:

```
Section: Assets
  Item 1: Cash (Formulas: Account 1001, Rule=Net, SumFactor=1)
    → Query ledger for 1001 → Net balance = $10,000

  Item 2: Bank Deposits (Formulas: Account 1002, Rule=Net, SumFactor=1)
    → Query ledger for 1002 → Net balance = $5,000

  Item 3: Current Assets Subtotal (Sum, SumFactor=1)
    → Sum of previous items = $10,000 + $5,000 = $15,000

  Section Total: $15,000 (sum of all items with SumFactor != 0)
```

**Validation**:

- **Balance Sheet**: Assets = Liabilities + Equity
- **Income Statement**: Custom validation rules (if configured)

**Re-generation**:

- Reports can be regenerated to reflect latest ledger data
- Useful when journals are posted after initial report generation

---

## 8. Numbering & Identifier Management

### 8.1 Purpose

Provide consistent, unique, sequential identifiers for business documents without manual intervention.

### 8.2 Configuration Structure

**Identifier Configuration**:

- **Target Business Object**: What type of document (e.g., "Journal")
- **Property Matchers**: Conditions that activate this configuration
  - Example: Period ID = "xyz" AND Journal Type = "General"
- **Counter**: Auto-incrementing number (starts at 0 or 1)
- **Prefix**: String before the counter (e.g., "JV-202401-")
- **Suffix**: String after the counter (e.g., empty or "-DRAFT")

### 8.3 Generation Logic

**When a journal is created**:

1. System calls numbering service with context (period ID, journal type)
2. Service finds matching configuration by property matchers
3. Service increments configuration counter
4. Service formats identifier: `{prefix}{counter}{suffix}`
5. Returns identifier (e.g., "JV-202401-0042")

**Example Configurations**:

```
Configuration:
  Target: Journal
  Matchers: Period = "Jan 2024" AND Type = "General"
  Prefix: "记"
  Counter: 42
  Result: "号"
```

**Benefits**:

- Sequential numbering within logical groups
- No number conflicts
- Audit trail (gaps indicate deleted/voided documents)

---

## 9. User Management

### 9.1 Authentication & Identity

**External Identity Provider**:

- **Ory Kratos**: Handles user authentication, password management, recovery
- **Ory Oathkeeper**: API gateway providing session management and routing

**User Registration Flow**:

1. Admin creates user via Kratos API (provides email)
2. System generates recovery link
3. User follows link to set password
4. User gains access to FIMS

### 9.2 User Data in FIMS

**Minimal User Records**:

- **User ID**: UUID (same as Kratos identity ID)
- **Traits**: JSONB field containing flexible user attributes (matches Ory Kratos identity structure)
  - Email for identification
  - Name for display
  - Extensible for additional attributes as needed

**Purpose**:

- Link journal activities to users (creator, reviewer, auditor, poster)
- Display user names in audit trails
- No sensitive authentication data (passwords, recovery tokens) stored in FIMS
- JSONB format allows seamless integration with Ory Kratos identity data

### 9.3 Authorization (Current State)

**Minimal Authorization**:

- Current implementation does not enforce role-based access control
- All authenticated users can perform all actions (within segregation of duties rules)
- Future enhancement: Implement roles and permissions per SoB

---

## 10. Key Business Rules & Constraints

### 10.1 Double-Entry Bookkeeping

- ✓ Every journal must have a balanced signed amount sum (sum of all journal line amounts equals zero)
- ✓ Minimum one journal line per journal (system validates non-empty list; double-entry convention typically requires 2+)
- ✓ Each journal line has a non-zero signed amount (positive = debit, negative = credit)

### 10.2 Segregation of Duties

- ✓ Creator ≠ Reviewer ≠ Auditor ≠ Poster (four different users)
- ✓ Enforced at domain level during state transitions
- ✓ Prevents single-user fraud

### 10.3 Period Control

- ✓ Journals can only be posted to the current period
- ✓ Current period must be open (not closed)
- ✓ Only one period can be current at a time
- ✓ Historical periods are read-only after closing

### 10.4 Period Closing Requirements

- ✓ All journals must be posted (none in draft/reviewed/audited state)
- ✓ All profit & loss accounts must have zero ending amount
- ✓ Trial balance must be satisfied (sum of all signed amounts equals zero)

### 10.5 Account Structure

- ✓ Only detail accounts (明细科目) can have transactions
- ✓ Parent accounts (上级科目/一级科目) aggregate child account balances
- ✓ Account numbers must be unique within a SoB
- ✓ Account code length defined by SoB (2-10 levels, 1-6 digits each)

### 10.6 Ledger Balance Integrity

- ✓ Opening amount is a signed decimal value (positive or negative)
- ✓ Period amount is signed net movement (endingAmount = openingAmount + periodAmount)
- ✓ Account balance direction guides expected amount sign
- ✓ Hierarchical posting updates all parent accounts

### 10.7 Journal Immutability

- ✓ Posted journals cannot be edited or deleted
- ✓ Corrections require creating new reversing/correcting journals
- ✓ Maintains complete audit trail

### 10.8 Report Validation

- ✓ Balance Sheet: Assets = Liabilities + Equity
- ✓ Income Statement: All items must calculate correctly
- ✓ Amount types must match report class

### 10.9 Transaction Boundaries

- ✓ All write operations use database transactions
- ✓ Journal posting is atomic (all ledgers update or none)
- ✓ Period closing is atomic (close period + create next + initialize ledgers)

### 10.10 SoB Isolation (Within Deployment)

- ✓ No cross-SoB data access within same deployment
- ✓ All domain entities scoped to single SoB
- ✓ Data isolation enforced at application level

---

## Appendix A: Common Business Scenarios

### Scenario 1: Simple Sales Transaction

```
Transaction: Sold goods for $10,000 cash

Journal:
  Line 1: Account 1101 - Cash           Amount: +$10,000  (debit)
  Line 2: Account 4010 - Sales Revenue  Amount: -$10,000 (credit)
```

### Scenario 2: Purchase on Credit

```
Transaction: Purchased inventory from Vendor ABC for $5,000 on credit

Journal:
  Line 1: Account 1500 - Inventory           Amount: +$5,000  (debit)
  Line 2: Account 2100 - Accounts Payable    Amount: -$5,000 (credit)

Result:
  - Inventory account increases by $5,000
  - A/P account increases by $5,000 (credit = negative amount)
```

### Scenario 3: Month-End Closing

```
Period: January 2024
Accounts:
  - 4010 - Sales Revenue (P&L): -$100,000 (credit balance, normal for revenue)
  - 5010 - Expenses (P&L): +$60,000 (debit balance, normal for expenses)
  - 3200 - Retained Earnings (Equity): existing balance

Closing Entry:
  Line 1: Account 4010 - Sales Revenue      Amount: +$100,000 (to clear credit)
  Line 2: Account 5010 - Expenses           Amount: -$60,000 (to clear debit)
  Line 3: Account 3200 - Retained Earnings  Amount: +$60,000
  Line 4: Account 3200 - Retained Earnings  Amount: -$100,000

(Above shown as two journals; actual may vary)

After Closing:
  - All P&L accounts have zero balance
  - Retained Earnings increased by $40,000 (net income)
  - Period can now be closed
  - February period initialized with B/S balances carried forward
```

---

## Appendix B: Glossary

| Term                       | Chinese           | Definition                                                                             |
| -------------------------- | ----------------- | -------------------------------------------------------------------------------------- |
| SoB                        | 账套              | Set of Books - Independent accounting entity within a deployment                       |
| Accounting Standards       | 会计准则          | Chinese accounting standards framework (currently supports 小企业会计准则)             |
| Small Enterprise Standards | 小企业会计准则    | Accounting Standards for Small Enterprises - defines standard CoA structure            |
| Account                    | 科目/会计科目     | Individual account in chart of accounts                                                |
| First-level Account        | 一级科目          | Top-level account from standard CoA (e.g., 1001, 1002, 1122)                           |
| Parent Account             | 上级科目          | Summary account that aggregates child balances (cannot post transactions)              |
| Detail Account             | 明细科目          | Leaf-level account that can accept transactions (user-defined under standard accounts) |
| Chart of Accounts          | 科目表/会计科目表 | Organized hierarchical list of all accounts                                            |
| Journal                    | 分录              | Journal entry recording financial transactions (interchangeable term)                  |
| Journal Entry              | 分录              | Same as Journal - both terms are correct in accounting                                 |
| Period                     | 会计期间          | Accounting period (monthly)                                                            |
| Ledger                     | 总账              | General ledger - balance records for all accounts (both parent and detail)             |
| Ledger Entry               | 账簿分录          | Detailed transaction record created when journal is posted (one per journal line)      |
| Review                     | 复核              | First-level approval of journal by reviewer                                            |
| Audit                      | 审核              | Second-level approval of journal by auditor                                            |
| Post                       | 登账/过账         | Final commit of journal to ledgers (makes it affect balances)                          |
| Balance Sheet              | 资产负债表        | Statement of financial position (assets = liabilities + equity)                        |
| Income Statement           | 利润表            | Profit and loss statement (revenue - expenses)                                         |
| Trial Balance              | 试算平衡          | Verification that sum of all signed amounts equals zero                                |
| Formula Rule               | 公式规则          | Rule for extracting value from ledger (Net/Debit/Credit/Transaction)                   |
| Sum Factor                 | 汇总因子          | +1 (add), -1 (subtract), or 0 (exclude) - controls how item contributes to parent sum  |

---

**End of Business Design Document**

_This document serves as the baseline reference for discussing FIMS business functionality with AI assistants and team members._
