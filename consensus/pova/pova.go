package pova

import (
	"errors"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/params"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/consensus"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/ethereum/go-ethereum/trie"
	)

type POVA struct {
	config      *params.ChainConfig
	validators  []common.Address
	currentStep int
}

func New(config *params.ChainConfig, validators []common.Address) *POVA {
	return &POVA{
		config:      config,
		validators:  validators,
		currentStep: 0,
	}
}

func (p *POVA) Author(header *types.Header) (common.Address, error) {
	return header.Coinbase, nil
}

func (p *POVA) VerifyHeader(chain consensus.ChainHeaderReader, header *types.Header) error {
	if header.Number.Uint64() == 0 {
		return nil
	}
	// Calculate which validator should have created this block
	blockNumber := header.Number.Uint64()
	expected := p.validators[(blockNumber-1)%uint64(len(p.validators))]
	if header.Coinbase != expected {
		return errors.New("unauthorized block proposer")
	}
	return nil
}

func (p *POVA) VerifyHeaders(chain consensus.ChainHeaderReader, headers []*types.Header) (chan<- struct{}, <-chan error) {
	abort := make(chan struct{})
	results := make(chan error, len(headers))
	
	go func() {
		for _, header := range headers {
			select {
			case <-abort:
				return
			case results <- p.VerifyHeader(chain, header):
			}
		}
	}()
	
	return abort, results
}

func (p *POVA) VerifyUncles(chain consensus.ChainReader, block *types.Block) error {
	// POVA doesn't use uncles, so we just return nil
	return nil
}

func (p *POVA) Prepare(chain consensus.ChainHeaderReader, header *types.Header) error {
	// Calculate which validator should create this block
	blockNumber := header.Number.Uint64()
	header.Coinbase = p.validators[(blockNumber-1)%uint64(len(p.validators))]
	header.Difficulty = big.NewInt(1)
	
	// Get the period from config (default 15 seconds)
	period := uint64(15)
	if p.config.POVA != nil {
		period = p.config.POVA.Period
	}
	
	// Calculate expected block time based on period
	parent := chain.GetHeader(header.ParentHash, header.Number.Uint64()-1)
	if parent != nil {
		// Ensure minimum time between blocks
		expectedTime := parent.Time + period
		currentTime := uint64(time.Now().Unix())
		if currentTime < expectedTime {
			header.Time = expectedTime
		} else {
			header.Time = currentTime
		}
	} else {
		header.Time = uint64(time.Now().Unix())
	}
	
	return nil
}

func (p *POVA) Finalize(chain consensus.ChainHeaderReader, header *types.Header, state *state.StateDB, txs []*types.Transaction, uncles []*types.Header, withdrawals []*types.Withdrawal) {
	// POVA doesn't have any special finalization logic
}

func (p *POVA) FinalizeAndAssemble(chain consensus.ChainHeaderReader, header *types.Header, state *state.StateDB, txs []*types.Transaction, uncles []*types.Header, receipts []*types.Receipt, withdrawals []*types.Withdrawal) (*types.Block, error) {
	// Commit the state to get the state root
	state.Commit(header.Number.Uint64(), true)
	
	// Set the state root in the header
	header.Root = state.IntermediateRoot(true)
	
	// Create the block with the committed state root
	block := types.NewBlock(header, txs, uncles, receipts, trie.NewStackTrie(nil))
	return block, nil
}

func (p *POVA) Seal(chain consensus.ChainHeaderReader, block *types.Block, results chan<- *types.Block, stop <-chan struct{}) error {
	// Get the period from config (default 15 seconds)
	period := uint64(15)
	if p.config.POVA != nil {
		period = p.config.POVA.Period
	}
	
	// Calculate delay until next block time
	parent := chain.GetHeader(block.ParentHash(), block.NumberU64()-1)
	if parent != nil {
		expectedTime := parent.Time + period
		currentTime := uint64(time.Now().Unix())
		if currentTime < expectedTime {
			delay := time.Duration(expectedTime-currentTime) * time.Second
			select {
			case <-time.After(delay):
			case <-stop:
				return errors.New("sealing stopped")
			}
		}
	}
	
	// Send the block
	select {
	case results <- block:
		return nil
	case <-stop:
		return errors.New("sealing stopped")
	}
}

func (p *POVA) SealHash(header *types.Header) common.Hash {
	return header.Hash()
}

func (p *POVA) CalcDifficulty(chain consensus.ChainHeaderReader, time uint64, parent *types.Header) *big.Int {
	return big.NewInt(1)
}

func (p *POVA) APIs(chain consensus.ChainHeaderReader) []rpc.API {
	return []rpc.API{}
}

func (p *POVA) Close() error {
	return nil
}
