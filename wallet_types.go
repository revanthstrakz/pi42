package pi42

// FuturesWalletResponse represents the futures wallet information
type FuturesWalletResponse struct {
	InrBalance             string `json:"inrBalance"`
	WalletBalance          string `json:"walletBalance"`
	WithdrawableBalance    string `json:"withdrawableBalance"`
	MaintenanceMargin      string `json:"maintenanceMargin"`
	UnrealisedPnlCross     string `json:"unrealisedPnlCross"`
	UnrealisedPnlIsolated  string `json:"unrealisedPnlIsolated"`
	MaxWithdrawableBalance string `json:"maxWithdrawableBalance"`
	LockedBalance          string `json:"lockedBalance"`
	MarginBalance          string `json:"marginBalance"`
	PnlPercentCross        string `json:"pnlPercentCross"`
	PnlPercentIsolated     string `json:"pnlPercentIsolated"`
	LockedBalanceCross     string `json:"lockedBalanceCross"`
	LockedBalanceIsolated  string `json:"lockedBalanceIsolated"`
	MarginAsset            string `json:"marginAsset"`
}

// FundingWalletResponse represents the funding wallet information
type FundingWalletResponse struct {
	InrBalance          string `json:"inrBalance"`
	WalletBalance       string `json:"walletBalance"`
	WithdrawableBalance string `json:"withdrawableBalance"`
	LockedBalance       string `json:"lockedBalance"`
	MarginAsset         string `json:"marginAsset"`
}
