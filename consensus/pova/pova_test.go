package pova

import (
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/params"
)

func TestPOVA_Author(t *testing.T) {
	validators := []common.Address{
		common.HexToAddress("0x1234567890123456789012345678901234567890"),
		common.HexToAddress("0x2345678901234567890123456789012345678901"),
	}
	
	config := &params.ChainConfig{
		ChainID: big.NewInt(1),
	}
	
	pova := New(config, validators)
	
	header := &types.Header{
		Coinbase: validators[0],
	}
	
	author, err := pova.Author(header)
	if err != nil {
		t.Fatalf("Author failed: %v", err)
	}
	
	if author != validators[0] {
		t.Errorf("Expected author %v, got %v", validators[0], author)
	}
}

func TestPOVA_VerifyHeader(t *testing.T) {
	validators := []common.Address{
		common.HexToAddress("0x1234567890123456789012345678901234567890"),
		common.HexToAddress("0x2345678901234567890123456789012345678901"),
	}
	
	config := &params.ChainConfig{
		ChainID: big.NewInt(1),
	}
	
	pova := New(config, validators)
	
	// Test genesis block (should pass)
	genesisHeader := &types.Header{
		Number: big.NewInt(0),
	}
	
	err := pova.VerifyHeader(nil, genesisHeader)
	if err != nil {
		t.Errorf("Genesis block verification failed: %v", err)
	}
	
	// Test valid proposer
	validHeader := &types.Header{
		Number:   big.NewInt(1),
		Coinbase: validators[0], // First validator should be the proposer for block 1
	}
	
	err = pova.VerifyHeader(nil, validHeader)
	if err != nil {
		t.Errorf("Valid header verification failed: %v", err)
	}
	
	// Test invalid proposer
	invalidHeader := &types.Header{
		Number:   big.NewInt(1),
		Coinbase: common.HexToAddress("0x9999999999999999999999999999999999999999"),
	}
	
	err = pova.VerifyHeader(nil, invalidHeader)
	if err == nil {
		t.Error("Invalid header verification should have failed")
	}
}

func TestPOVA_Prepare(t *testing.T) {
	validators := []common.Address{
		common.HexToAddress("0x1234567890123456789012345678901234567890"),
		common.HexToAddress("0x2345678901234567890123456789012345678901"),
	}
	
	config := &params.ChainConfig{
		ChainID: big.NewInt(1),
	}
	
	pova := New(config, validators)
	
	header := &types.Header{
		Number: big.NewInt(1),
	}
	
	err := pova.Prepare(nil, header)
	if err != nil {
		t.Fatalf("Prepare failed: %v", err)
	}
	
	if header.Coinbase != validators[0] {
		t.Errorf("Expected coinbase %v, got %v", validators[0], header.Coinbase)
	}
	
	if header.Difficulty.Cmp(big.NewInt(1)) != 0 {
		t.Errorf("Expected difficulty 1, got %v", header.Difficulty)
	}
}

func TestPOVA_Seal(t *testing.T) {
	validators := []common.Address{
		common.HexToAddress("0x1234567890123456789012345678901234567890"),
		common.HexToAddress("0x2345678901234567890123456789012345678901"),
	}
	
	config := &params.ChainConfig{
		ChainID: big.NewInt(1),
	}
	
	pova := New(config, validators)
	
	block := types.NewBlock(&types.Header{
		Number: big.NewInt(1),
	}, nil, nil, nil, nil)
	
	results := make(chan *types.Block, 1)
	stop := make(chan struct{})
	
	err := pova.Seal(nil, block, results, stop)
	if err != nil {
		t.Fatalf("Seal failed: %v", err)
	}
	
	select {
	case sealedBlock := <-results:
		if sealedBlock != block {
			t.Error("Sealed block should be the same as input block")
		}
	case <-time.After(time.Second):
		t.Error("Seal timeout")
	}
}

func TestPOVA_ValidatorRotation(t *testing.T) {
	validators := []common.Address{
		common.HexToAddress("0x1234567890123456789012345678901234567890"),
		common.HexToAddress("0x2345678901234567890123456789012345678901"),
		common.HexToAddress("0x3456789012345678901234567890123456789012"),
	}
	
	config := &params.ChainConfig{
		ChainID: big.NewInt(1),
	}
	
	pova := New(config, validators)
	
	// Test first block (block 1 should be created by validator 0)
	header1 := &types.Header{Number: big.NewInt(1)}
	pova.Prepare(nil, header1)
	if header1.Coinbase != validators[0] {
		t.Errorf("Block 1: Expected %v, got %v", validators[0], header1.Coinbase)
	}
	
	// Test second block (block 2 should be created by validator 1)
	header2 := &types.Header{Number: big.NewInt(2)}
	pova.Prepare(nil, header2)
	if header2.Coinbase != validators[1] {
		t.Errorf("Block 2: Expected %v, got %v", validators[1], header2.Coinbase)
	}
	
	// Test third block (block 3 should be created by validator 2)
	header3 := &types.Header{Number: big.NewInt(3)}
	pova.Prepare(nil, header3)
	if header3.Coinbase != validators[2] {
		t.Errorf("Block 3: Expected %v, got %v", validators[2], header3.Coinbase)
	}
	
	// Test fourth block (block 4 should be created by validator 0 again)
	header4 := &types.Header{Number: big.NewInt(4)}
	pova.Prepare(nil, header4)
	if header4.Coinbase != validators[0] {
		t.Errorf("Block 4: Expected %v, got %v", validators[0], header4.Coinbase)
	}
} 