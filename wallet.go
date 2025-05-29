package pi42

import (
	"encoding/json"
	"fmt"
)

// WalletAPI provides access to wallet information endpoints
type WalletAPI struct {
	client *Client
}

// NewWalletAPI creates a new Wallet API instance
func NewWalletAPI(client *Client) *WalletAPI {
	return &WalletAPI{client: client}
}

// FuturesWalletDetails gets all details of Futures wallet with structured response
// marginAsset: Asset to query wallet details for (e.g., "INR", "USDT")
func (api *WalletAPI) FuturesWalletDetails(marginAsset string) (*FuturesWalletResponse, error) {
	endpoint := "/v1/wallet/futures-wallet/details"

	params := make(map[string]string)
	if marginAsset != "" {
		params["marginAsset"] = marginAsset
	} else {
		params["marginAsset"] = "INR"
	}

	data, err := api.client.Get(endpoint, params, false)
	if err != nil {
		return nil, err
	}

	var result FuturesWalletResponse
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("error parsing response: %v", err)
	}

	return &result, nil
}

// FundingWalletDetails gets details of funding wallet with structured response
// marginAsset: Asset to query wallet details for (e.g., "INR", "USDT")
func (api *WalletAPI) FundingWalletDetails(marginAsset string) (*FundingWalletResponse, error) {
	endpoint := "/v1/wallet/funding-wallet/details"

	params := make(map[string]string)
	if marginAsset != "" {
		params["marginAsset"] = marginAsset
	} else {
		params["marginAsset"] = "INR"
	}

	data, err := api.client.Get(endpoint, params, false)
	if err != nil {
		return nil, err
	}

	var result FundingWalletResponse
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("error parsing response: %v", err)
	}

	return &result, nil
}
