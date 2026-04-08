package services

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

// ChainInfo holds metadata about a supported blockchain.
type ChainInfo struct {
	Symbol   string
	Decimals int
	Name     string
}

// Supported chains mapped by chain identifier.
var SupportedChains = map[string]ChainInfo{
	"solana":     {Symbol: "SOL", Decimals: 9, Name: "Solana"},
	"ethereum":   {Symbol: "ETH", Decimals: 18, Name: "Ethereum"},
	"bitcoin":    {Symbol: "BTC", Decimals: 8, Name: "Bitcoin"},
	"usdt_erc20": {Symbol: "USDT", Decimals: 6, Name: "USDT (ERC-20)"},
	"usdc_erc20": {Symbol: "USDC", Decimals: 6, Name: "USDC (ERC-20)"},
	"usdt_trc20": {Symbol: "USDT", Decimals: 6, Name: "USDT (TRC-20)"},
	"bsc":        {Symbol: "BNB", Decimals: 18, Name: "BNB Chain"},
	"polygon":    {Symbol: "MATIC", Decimals: 18, Name: "Polygon"},
	"arbitrum":   {Symbol: "ETH", Decimals: 18, Name: "Arbitrum"},
}

// ChainService provides chain-agnostic balance queries.
type ChainService interface {
	GetBalance(ctx context.Context, address string, chain string) (float64, error)
	GetChainInfo(chain string) ChainInfo
	IsChainSupported(chain string) bool
}

// MoralisProvider implements ChainService using the Moralis API.
type MoralisProvider struct {
	apiKey    string
	proxyURL  string
	client    *http.Client
	solanaRPC *SolanaService // fallback for Solana
}

// Moralis balance response structures
type moralisNativeBalance struct {
	Balance string `json:"balance"`
}

type moralisTokenBalance struct {
	Balance  string `json:"balance"`
	Decimals int    `json:"decimals"`
	Symbol   string `json:"symbol"`
}

// Moralis chain ID mapping
var moralisChainIDs = map[string]string{
	"ethereum": "0x1",
	"bsc":      "0x38",
	"polygon":  "0x89",
	"arbitrum": "0xa4b1",
	"bitcoin":  "0x0", // Moralis uses special handling for BTC
}

func NewMoralisProvider(apiKey string, proxyURL string, solanaRPC *SolanaService) *MoralisProvider {
	var client *http.Client
	if proxyURL != "" {
		if parsed, err := url.Parse(proxyURL); err == nil {
			client = &http.Client{
				Timeout:   30 * time.Second,
				Transport: &http.Transport{Proxy: http.ProxyURL(parsed)},
			}
		}
	}
	if client == nil {
		client = &http.Client{Timeout: 30 * time.Second}
	}

	return &MoralisProvider{
		apiKey:    apiKey,
		proxyURL:  proxyURL,
		client:    client,
		solanaRPC: solanaRPC,
	}
}

func (m *MoralisProvider) GetBalance(ctx context.Context, address string, chain string) (float64, error) {
	info, ok := SupportedChains[chain]
	if !ok {
		return 0, fmt.Errorf("unsupported chain: %s", chain)
	}

	// Use direct Solana RPC for Solana chain (already working)
	if chain == "solana" && m.solanaRPC != nil {
		return m.solanaRPC.GetBalance(ctx, address)
	}

	// Use Moralis for other chains
	return m.getMoralisBalance(ctx, address, chain, info)
}

func (m *MoralisProvider) getMoralisBalance(ctx context.Context, address string, chain string, info ChainInfo) (float64, error) {
	// For token chains (USDT/USDC), use token balance endpoint
	if chain == "usdt_erc20" || chain == "usdc_erc20" || chain == "usdt_trc20" {
		return m.getMoralisTokenBalance(ctx, address, chain)
	}

	// For native chains, use native balance endpoint
	return m.getMoralisNativeBalance(ctx, address, chain, info)
}

func (m *MoralisProvider) getMoralisNativeBalance(ctx context.Context, address string, chain string, info ChainInfo) (float64, error) {
	chainID, ok := moralisChainIDs[chain]
	if !ok {
		return 0, fmt.Errorf("no Moralis chain ID for: %s", chain)
	}

	apiURL := fmt.Sprintf("https://deep-index.moralis.io/api/v2.2/%s/balance?chain=%s", address, chainID)

	req, err := http.NewRequestWithContext(ctx, "GET", apiURL, nil)
	if err != nil {
		return 0, err
	}
	req.Header.Set("X-API-Key", m.apiKey)
	req.Header.Set("Accept", "application/json")

	resp, err := m.client.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return 0, fmt.Errorf("Moralis API error: HTTP %d", resp.StatusCode)
	}

	var result moralisNativeBalance
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return 0, err
	}

	// Parse balance string (wei/satoshi) and convert to main unit
	return parseBalanceString(result.Balance, info.Decimals)
}

func (m *MoralisProvider) getMoralisTokenBalance(ctx context.Context, address string, chain string) (float64, error) {
	// Map token chains to their parent chain and token addresses
	type tokenInfo struct {
		parentChain string
		tokenAddr   string
	}

	tokens := map[string]tokenInfo{
		"usdt_erc20": {parentChain: "0x1", tokenAddr: "0xdAC17F958D2ee523a2206206994597C13D831ec7"},
		"usdc_erc20": {parentChain: "0x1", tokenAddr: "0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48"},
		"usdt_trc20": {parentChain: "0x1", tokenAddr: ""}, // TRC-20 not directly supported, will return error
	}

	tok, ok := tokens[chain]
	if !ok || tok.tokenAddr == "" {
		return 0, fmt.Errorf("token balance not supported for chain: %s", chain)
	}

	apiURL := fmt.Sprintf("https://deep-index.moralis.io/api/v2.2/%s/erc20?chain=%s&token_addresses=%s", address, tok.parentChain, tok.tokenAddr)

	req, err := http.NewRequestWithContext(ctx, "GET", apiURL, nil)
	if err != nil {
		return 0, err
	}
	req.Header.Set("X-API-Key", m.apiKey)
	req.Header.Set("Accept", "application/json")

	resp, err := m.client.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return 0, fmt.Errorf("Moralis token API error: HTTP %d", resp.StatusCode)
	}

	var result []moralisTokenBalance
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return 0, err
	}

	if len(result) == 0 {
		return 0, nil
	}

	return parseBalanceString(result[0].Balance, result[0].Decimals)
}

func (m *MoralisProvider) GetChainInfo(chain string) ChainInfo {
	if info, ok := SupportedChains[chain]; ok {
		return info
	}
	return ChainInfo{Symbol: "???", Decimals: 0, Name: "Unknown"}
}

func (m *MoralisProvider) IsChainSupported(chain string) bool {
	_, ok := SupportedChains[chain]
	return ok
}

// parseBalanceString converts a wei/satoshi string to float
func parseBalanceString(balanceStr string, decimals int) (float64, error) {
	if balanceStr == "" || balanceStr == "0" {
		return 0, nil
	}

	var val float64
	if _, err := fmt.Sscanf(balanceStr, "%f", &val); err != nil {
		return 0, err
	}

	divisor := 1.0
	for i := 0; i < decimals; i++ {
		divisor *= 10
	}

	return val / divisor, nil
}
