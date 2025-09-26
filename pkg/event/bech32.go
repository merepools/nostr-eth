package event

import (
	"encoding/binary"
	"fmt"
	"strconv"

	"github.com/btcsuite/btcutil/bech32"
)

// EncodeEventIDToNevent encodes an event ID to NIP-19 nevent format
// This function creates a bech32-encoded string with the nevent prefix
// that includes the event ID and optional metadata (relay URL, author pubkey, kind)
func EncodeEventIDToNevent(eventID string, relayURL string, authorPubkey string, kind int) (string, error) {
	// Convert hex event ID to bytes
	eventIDBytes, err := hexToBytes(eventID)
	if err != nil {
		return "", fmt.Errorf("invalid event ID: %v", err)
	}

	// Build TLV structure using proper NIP-19 format
	var tlvData []byte

	// Type 0 (special): 32-byte event ID
	tlvData = append(tlvData, 0)                       // type
	tlvData = append(tlvData, byte(len(eventIDBytes))) // length
	tlvData = append(tlvData, eventIDBytes...)         // value

	// Type 1 (relay): relay URL (optional)
	if relayURL != "" {
		relayBytes := []byte(relayURL)
		tlvData = append(tlvData, 1)                     // type
		tlvData = append(tlvData, byte(len(relayBytes))) // length
		tlvData = append(tlvData, relayBytes...)         // value
	}

	// Type 2 (author): author public key (optional)
	if authorPubkey != "" {
		authorBytes, err := hexToBytes(authorPubkey)
		if err != nil {
			return "", fmt.Errorf("invalid author pubkey: %v", err)
		}
		tlvData = append(tlvData, 2)                      // type
		tlvData = append(tlvData, byte(len(authorBytes))) // length
		tlvData = append(tlvData, authorBytes...)         // value
	}

	// Type 3 (kind): event kind (optional)
	if kind > 0 {
		kindBytes := make([]byte, 4)
		binary.BigEndian.PutUint32(kindBytes, uint32(kind))
		tlvData = append(tlvData, 3)                    // type
		tlvData = append(tlvData, byte(len(kindBytes))) // length
		tlvData = append(tlvData, kindBytes...)         // value
	}

	// Convert to 5-bit groups for bech32 encoding
	converted, err := bech32.ConvertBits(tlvData, 8, 5, true)
	if err != nil {
		return "", fmt.Errorf("failed to convert bits: %v", err)
	}

	// Encode using bech32 with nevent prefix
	encoded, err := bech32.Encode("nevent", converted)
	if err != nil {
		return "", fmt.Errorf("bech32 encoding failed: %v", err)
	}

	return encoded, nil
}

// hexToBytes converts a hex string to bytes
func hexToBytes(hex string) ([]byte, error) {
	if len(hex)%2 != 0 {
		return nil, fmt.Errorf("hex string must have even length")
	}

	bytes := make([]byte, len(hex)/2)
	for i := 0; i < len(hex); i += 2 {
		val, err := strconv.ParseUint(hex[i:i+2], 16, 8)
		if err != nil {
			return nil, fmt.Errorf("invalid hex character at position %d: %v", i, err)
		}
		bytes[i/2] = byte(val)
	}
	return bytes, nil
}
