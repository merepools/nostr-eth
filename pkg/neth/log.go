package neth

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/big"
	"sort"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

const (
	ERC20TransferTopic = "0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef"
	DataKeyFrom        = "from"
	DataKeyTo          = "to"
	DataKeyTopic       = "topic"
	DataKeyValue       = "value"
)

type Log struct {
	Hash      string           `json:"hash"`
	TxHash    string           `json:"tx_hash"`
	ChainID   string           `json:"chain_id"`
	Topic     string           `json:"topic"`
	CreatedAt time.Time        `json:"created_at"`
	UpdatedAt time.Time        `json:"updated_at"`
	Nonce     int64            `json:"nonce"`
	Sender    string           `json:"sender"`
	To        string           `json:"to"`
	Value     *big.Int         `json:"value"`
	Data      *json.RawMessage `json:"data"`
}

type LogTransferData struct {
	To    string `json:"to"`
	From  string `json:"from"`
	Topic string `json:"topic"`
	Value string `json:"value"`
}

// generate hash for transfer using a provided index, from, to and the tx hash
func (t *Log) GenerateUniqueHash() string {
	buf := new(bytes.Buffer)

	// Write each value to the buffer as bytes
	// Convert t.Value to a fixed-length byte representation
	valueBytes := t.Value.Bytes()
	buf.Write(common.LeftPadBytes(valueBytes, 32))
	if t.Data != nil {
		buf.Write(sortedJSONBytes(t.Data))
	}

	buf.Write(common.FromHex(t.TxHash))
	buf.Write(common.FromHex(t.ChainID))

	hash := crypto.Keccak256Hash(buf.Bytes())
	return hash.Hex()
}

func (t *Log) ToRounded(decimals int64) float64 {
	v, _ := t.Value.Float64()

	if decimals == 0 {
		return v
	}

	// Calculate value * 10^x
	multiplier, _ := new(big.Int).Exp(big.NewInt(10), big.NewInt(decimals), nil).Float64()

	result, _ := new(big.Float).Quo(big.NewFloat(v), big.NewFloat(multiplier)).Float64()

	return result
}

// Update updates the transfer using the given transfer
func (t *Log) Update(tx *Log) {
	// update all fields
	t.Hash = tx.Hash
	t.TxHash = tx.TxHash
	t.ChainID = tx.ChainID
	t.Topic = tx.Topic
	t.CreatedAt = tx.CreatedAt
	t.UpdatedAt = time.Now()
	t.Nonce = tx.Nonce
	t.Sender = tx.Sender
	t.To = tx.To
	t.Value = tx.Value
	t.Data = tx.Data
}

func (t *Log) GetPoolTopic() *string {
	if t.Data == nil {
		return nil
	}

	var data map[string]any

	json.Unmarshal(*t.Data, &data)

	v, ok := data["topic"].(string)
	if !ok {
		return nil
	}

	topic := strings.ToLower(fmt.Sprintf("%s/%s", t.To, v))

	return &topic
}

func (t *Log) GetEventData() (map[string]interface{}, error) {
	if t.Data == nil {
		return nil, nil
	}

	var data map[string]interface{}

	err := json.Unmarshal(*t.Data, &data)
	if err != nil {
		return nil, err
	}

	return data, nil
}

// Convert a log to json bytes
func (t *Log) MarshalJSON() ([]byte, error) {
	b, err := json.Marshal(t)
	if err != nil {
		return nil, err
	}

	return b, nil
}

func sortedJSONBytes(data *json.RawMessage) []byte {
	if data == nil {
		return nil
	}

	var m map[string]interface{}
	if err := json.Unmarshal(*data, &m); err != nil {
		// If it's not a JSON object, return the raw bytes
		return *data
	}

	// Get sorted keys
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	// Build sorted buffer
	var buf bytes.Buffer
	for _, k := range keys {
		v := m[k]
		keyBytes, _ := json.Marshal(k)
		valueBytes, _ := json.Marshal(v)
		buf.Write(keyBytes)
		buf.Write(valueBytes)
	}

	return buf.Bytes()
}
