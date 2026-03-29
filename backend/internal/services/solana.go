package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"net/url"
	"time"
)

type SolanaService struct {
	rpcURL  string
	proxyURL string
	client  *http.Client
}

type RPCRequest struct {
	JsonRPC string        `json:"jsonrpc"`
	Method  string        `json:"method"`
	Params  []interface{} `json:"params"`
	ID      int           `json:"id"`
}

type RPCResponse struct {
	JsonRPC string          `json:"jsonrpc"`
	Result  json.RawMessage `json:"result"`
	Error   *RPCError       `json:"error,omitempty"`
	ID      int             `json:"id"`
}

type RPCError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type GetBalanceResponse struct {
	Context struct {
		Slot int64 `json:"slot"`
	} `json:"context"`
	Value float64 `json:"value"`
}

type GetMultipleAccountsResponse struct {
	Context struct {
		Slot int64 `json:"slot"`
	} `json:"context"`
	Value []AccountData `json:"value"`
}

type AccountData struct {
	Executable bool   `json:"executable"`
	Lamports   int64  `json:"lamports"`
	Owner      string `json:"owner"`
	Data       []string `json:"data"`
}

func NewSolanaService(rpcURL string, proxyURL string) *SolanaService {
	var client *http.Client
	if proxyURL != "" {
		proxyURLParsed, err := url.Parse(proxyURL)
		if err == nil {
			client = &http.Client{
				Timeout: 30 * time.Second,
				Transport: &http.Transport{Proxy: http.ProxyURL(proxyURLParsed)},
			}
		} else {
			client = &http.Client{Timeout: 30 * time.Second}
		}
	} else {
		client = &http.Client{Timeout: 30 * time.Second}
	}

	return &SolanaService{
		rpcURL:  rpcURL,
		proxyURL: proxyURL,
		client:  client,
	}
}

func (s *SolanaService) GetBalance(ctx context.Context, address string) (float64, error) {
	params := []interface{}{address, "processed"}
	req := RPCRequest{
		JsonRPC: "2.0",
		Method:  "getBalance",
		Params:  params,
		ID:      1,
	}

	body, err := json.Marshal(req)
	if err != nil {
		return 0, err
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", s.rpcURL, bytes.NewReader(body))
	if err != nil {
		return 0, err
	}
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := s.client.Do(httpReq)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return 0, fmt.Errorf("HTTP error: %d", resp.StatusCode)
	}

	var result RPCResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return 0, err
	}

	if result.Error != nil {
		return 0, fmt.Errorf("RPC error: %s", result.Error.Message)
	}

	var balanceResp GetBalanceResponse
	if err := json.Unmarshal(result.Result, &balanceResp); err != nil {
		return 0, err
	}

	return float64(balanceResp.Value) / math.Pow(10, 9), nil
}

func (s *SolanaService) GetMultipleBalances(ctx context.Context, addresses []string) (map[string]float64, error) {
	if len(addresses) == 0 {
		return make(map[string]float64), nil
	}

	params := []interface{}{addresses, map[string]string{"dataSlice": "base64"}}
	req := RPCRequest{
		JsonRPC: "2.0",
		Method:  "getMultipleAccounts",
		Params:  params,
		ID:      1,
	}

	body, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", s.rpcURL, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := s.client.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("HTTP error: %d", resp.StatusCode)
	}

	var result RPCResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	if result.Error != nil {
		return nil, fmt.Errorf("RPC error: %s", result.Error.Message)
	}

	var accountsResp GetMultipleAccountsResponse
	if err := json.Unmarshal(result.Result, &accountsResp); err != nil {
		return nil, err
	}

	balances := make(map[string]float64)
	for i, addr := range addresses {
		if i < len(accountsResp.Value) && accountsResp.Value[i].Lamports > 0 {
			balances[addr] = float64(accountsResp.Value[i].Lamports) / math.Pow(10, 9)
		} else {
			balances[addr] = 0
		}
	}

	return balances, nil
}

func (s *SolanaService) TestConnection(ctx context.Context) error {
	params := []interface{}{}
	req := RPCRequest{
		JsonRPC: "2.0",
		Method:  "getVersion",
		Params:  params,
		ID:      1,
	}

	body, _ := json.Marshal(req)
	httpReq, _ := http.NewRequestWithContext(ctx, "POST", s.rpcURL, bytes.NewReader(body))
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := s.client.Do(httpReq)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("HTTP error: %d", resp.StatusCode)
	}

	var result RPCResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return err
	}

	if result.Error != nil {
		return fmt.Errorf("RPC error: %s", result.Error.Message)
	}

	return nil
}

func (s *SolanaService) GetTokenBalance(ctx context.Context, address, mint string) (float64, error) {
	return 0, fmt.Errorf("not implemented")
}

func readBody(resp *http.Response) string {
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	return string(body)
}