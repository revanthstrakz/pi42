package pi42

// ExchangeInfoResponse represents the full response from the Exchange Info endpoint
type ExchangeInfoResponse struct {
	Markets         []string           `json:"markets"`
	Contracts       []ContractData     `json:"contracts"`
	Tags            []string           `json:"tags"`
	AssetPrecisions map[string]int     `json:"assetPrecisions"`
	ConversionRates map[string]float64 `json:"conversionRates"`
}

// ContractData represents a trading contract from the exchange info response
type ContractData struct {
	Name                            string        `json:"name"`
	ContractName                    string        `json:"contractName"`
	Slug                            string        `json:"slug"`
	Tags                            []string      `json:"tags"`
	Filters                         []Filter      `json:"filters"`
	MakerFee                        float64       `json:"makerFee"`
	TakerFee                        float64       `json:"takerFee"`
	BaseAsset                       string        `json:"baseAsset"`
	OrderTypes                      []OrderType   `json:"orderTypes"`
	QuoteAsset                      string        `json:"quoteAsset"`
	MaxLeverage                     string        `json:"maxLeverage"`
	ContractType                    string        `json:"contractType"`
	MarginBuffer                    string        `json:"marginBuffer,omitempty"`
	DepthGrouping                   []string      `json:"depthGrouping,omitempty"`
	LiquidationFee                  string        `json:"liquidationFee"`
	PricePrecision                  string        `json:"pricePrecision"`
	IsDefaultContract               bool          `json:"isDefaultContract"`
	QuantityPrecision               string        `json:"quantityPrecision"`
	MaintMarginPercent              string        `json:"maintMarginPercent,omitempty"`
	LimitPriceVarAllowed            string        `json:"limitPriceVarAllowed"`
	MarginBufferPercentage          string        `json:"marginBufferPercentage"`
	MaintenanceMarginPercentage     string        `json:"maintenanceMarginPercentage"`
	IconUrl                         string        `json:"iconUrl"`
	ReduceMarginAllowedRatioPercent int           `json:"reduceMarginAllowedRatioPercent"`
	Market                          string        `json:"market"`
	MarginAssetsSupported           []string      `json:"marginAssetsSupported"`
	FundingFeeInterval              int           `json:"fundingFeeInterval"`
	MaintenanceMarginConfig         []interface{} `json:"maintenanceMarginConfig"`
}

// Filter represents a trading filter applied to a contract
type Filter struct {
	FilterType string `json:"filterType"`
	MinQty     string `json:"minQty,omitempty"`
	MaxQty     string `json:"maxQty,omitempty"`
	Limit      string `json:"limit,omitempty"`
	Notional   string `json:"notional,omitempty"`
}

// PreferenceUpdateResponse represents the response from updating trading preferences
type PreferenceUpdateResponse struct {
	ContractName    string `json:"contractName"`
	MarginMode      string `json:"marginMode"`
	UpdatedLeverage int    `json:"updatedLeverage"`
}

// LeverageUpdateResponse represents the response from updating leverage
type LeverageUpdateResponse struct {
	UpdatedLeverage int    `json:"updatedLeverage"`
	ContractName    string `json:"contractName"`
}
