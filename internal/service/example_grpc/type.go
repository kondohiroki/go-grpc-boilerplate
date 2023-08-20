package example_grpc

import "time"

type SponsorDataConditions struct {
	NSLInterestRate    float64   `json:"nsl_interest_rate"`
	NSDInterestRate    float64   `json:"nsd_interest_rate"`
	SDLInterestRate    float64   `json:"sdl_interest_rate"`
	SDTInterestRate    float64   `json:"sdt_interest_rate"`
	LoanToValueRate    float64   `json:"loan_to_value_rate"`
	FundTransferMethod string    `json:"fund_transfer_method"`
	Status             string    `json:"status"`
	CreatedTime        time.Time `json:"created_time"`
	UpdatedTime        time.Time `json:"updated_time"`
}

type SupplierBankAccount struct {
	SupplierAccountID  string    `json:"supplier_account_id"`
	SupplierBankType   string    `json:"supplier_bank_type"`
	AccountDisplayName string    `json:"account_display_name"`
	BankCode           string    `json:"bank_code"`
	BankName           string    `json:"bank_name"`
	BankAccountNo      string    `json:"bank_account_no"`
	Status             string    `json:"status"`
	CreatedTime        time.Time `json:"created_time"`
	UpdatedTime        time.Time `json:"updated_time"`
}

type SupplierBranch struct {
	SupplierBranchID             string    `json:"supplier_branch_id"`
	BranchType                   string    `json:"branch_type"`
	BranchName                   string    `json:"branch_name"`
	SupplierAccountRefNo         string    `json:"supplier_account_ref_no"`
	SupplierAccountBankName      string    `json:"supplier_account_bank_name"`
	SupplierAccountBankAccountNo string    `json:"supplier_account_bank_account_no"`
	Status                       string    `json:"status"`
	CreatedTime                  time.Time `json:"created_time"`
	UpdatedTime                  time.Time `json:"updated_time"`
}
