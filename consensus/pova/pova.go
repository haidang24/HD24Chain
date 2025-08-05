// POVA là consensus engine cho chuỗi POVA
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
	"github.com/ethereum/go-ethereum/crypto"
    "log"
	)

type POVA struct {
	config      *params.ChainConfig // Cấu hình chuỗi
	validators  []common.Address // Danh sách validator 
	currentStep int // Bước hiện tại
}

// Khởi tạo POVA
func New(config *params.ChainConfig, validators []common.Address) *POVA { 
	return &POVA{
		config:      config,
		validators:  validators,
		currentStep: 0,
	}
}

// Lấy địa chỉ validator tạo block
func (p *POVA) Author(header *types.Header) (common.Address, error) {
	return header.Coinbase, nil  //header.Coinbase là địa chỉ validator tạo block
}

// Kiểm tra header của block
func (p *POVA) VerifyHeader(chain consensus.ChainHeaderReader, header *types.Header) error {
	if header.Number.Uint64() == 0 {
		return nil // Nếu block là block đầu tiên, trả về nil
	}
	// Tính toán validator tạo block
	blockNumber := header.Number.Uint64() // Lấy số thứ tự block
	expected := p.validators[(blockNumber-1)%uint64(len(p.validators))] // Lấy địa chỉ validator tạo block
	if header.Coinbase != expected { // Nếu địa chỉ validator tạo block không khớp với địa chỉ validator tính toán, trả về lỗi
		return errors.New("unauthorized block proposer") // Lỗi không được phép tạo block
	}
	// // Kiểm tra signature của block
	// if err := p.verifySignature(header, expected); err != nil {
	// 	return err // Lỗi không được phép tạo block
	// }
	return nil // Nếu khớp, trả về nil
}

// // Kiểm tra signature của block
// func (p *POVA) verifySignature(header *types.Header, expected common.Address) error {
// 	hash := header.Hash().Bytes() // Lấy hash của block
// 	sig := header.Extra // Lấy signature của block
// 	pubkey, err := crypto.SigToPub(hash, sig) // Lấy public key từ signature
// 	if err != nil {
// 		return err // Lỗi không được phép tạo block
// 	}
// 	recovered := crypto.PubkeyToAddress(*pubkey) // Lấy địa chỉ từ public key
// 	if recovered != expected {
// 		return errors.New("invalid block signature") // Lỗi không được phép tạo block
// 	}
// 	return nil // Nếu khớp, trả về nil
// }

// Kiểm tra nhiều header của block
func (p *POVA) VerifyHeaders(chain consensus.ChainHeaderReader, headers []*types.Header) (chan<- struct{}, <-chan error) {
	abort := make(chan struct{}) // Kênh để dừng việc kiểm tra
	results := make(chan error, len(headers)) // Kênh để lấy kết quả kiểm tra
	
	go func() {
		for _, header := range headers {
			select {
			case <-abort:
				return // Nếu dừng việc kiểm tra, trả về
			case results <- p.VerifyHeader(chain, header): // Kiểm tra header của block
			}
		}
	}()
	
	return abort, results // Trả về kênh để dừng việc kiểm tra và kênh để lấy kết quả kiểm tra
}

//uncles là các block cha của block hiện tại
// Kiểm tra uncles của block
func (p *POVA) VerifyUncles(chain consensus.ChainReader, block *types.Block) error {
	// POVA không sử dụng uncles, nên trả về nil vì không có uncles
	return nil
}

// Chuẩn bị block
func (p *POVA) Prepare(chain consensus.ChainHeaderReader, header *types.Header) error {
	// Tính toán validator tạo block
	blockNumber := header.Number.Uint64() // Lấy số thứ tự block
	//công thức tính toán validator tạo block: (blockNumber-1)%uint64(len(p.validators))
	header.Coinbase = p.validators[(blockNumber-1)%uint64(len(p.validators))] 
	header.Difficulty = big.NewInt(1) // Đặt độ khó của block là 1
	
	// Set default period là 15 giây
	period := uint64(15) 

	// Nếu có POVA trong config thì lấy chu kỳ từ config (genesis.json)
	if p.config.POVA != nil {
		period = p.config.POVA.Period 
	}
	
	// Tính toán thời gian block dựa trên period
	parent := chain.GetHeader(header.ParentHash, header.Number.Uint64()-1) // Lấy header cha của block
	
	//Nếu có header cha thì tính toán thời gian block dựa trên period
	if parent != nil { 
		// Đảm bảo thời gian giữa các block là tối thiểu
		expectedTime := parent.Time + period // Tính toán thời gian dựa trên period
		currentTime := uint64(time.Now().Unix()) // Lấy thời gian hiện tại
		if currentTime < expectedTime { // Nếu thời gian hiện tại nhỏ hơn thời gian dựa trên period
			header.Time = expectedTime // Đặt thời gian của block là thời gian dựa trên period
		} else {
			header.Time = currentTime // Đặt thời gian của block là thời gian hiện tại
		}
	} else {
		header.Time = uint64(time.Now().Unix()) // Đặt thời gian của block là thời gian hiện tại
	}
	
	return nil // Trả về nil
}

// Kết thúc block
func (p *POVA) Finalize(chain consensus.ChainHeaderReader, header *types.Header, state *state.StateDB, txs []*types.Transaction, uncles []*types.Header, withdrawals []*types.Withdrawal) {
	// POVA không có logic kết thúc đặc biệt
	// Không cần thiết
}

// Kết thúc và tạo block
func (p *POVA) FinalizeAndAssemble(chain consensus.ChainHeaderReader, header *types.Header, state *state.StateDB, txs []*types.Transaction, uncles []*types.Header, receipts []*types.Receipt, withdrawals []*types.Withdrawal) (*types.Block, error) {
	state.Commit(header.Number.Uint64(), true) // Commit state để lấy state root
	header.Root = state.IntermediateRoot(true) // Đặt state root trong header
	
	block := types.NewBlock(header, txs, uncles, receipts, trie.NewStackTrie(nil)) // Tạo block với state root đã commit
	return block, nil // Trả về block
}

// Tạo block
func (p *POVA) Seal(chain consensus.ChainHeaderReader, block *types.Block, results chan<- *types.Block, stop <-chan struct{}) error {
	// Set default period là 15 giây
	period := uint64(15) 

	// Nếu có POVA trong config thì lấy chu kỳ từ config (genesis.json)
	if p.config.POVA != nil {
		period = p.config.POVA.Period 
	}
	
	// Tính toán thời gian chờ cho block tiếp theo
	parent := chain.GetHeader(block.ParentHash(), block.NumberU64()-1) // Lấy header cha của block
	if parent != nil { // Nếu có header cha
		expectedTime := parent.Time + period // Tính toán thời gian dựa trên period
		currentTime := uint64(time.Now().Unix()) // Lấy thời gian hiện tại
		if currentTime < expectedTime { // Nếu thời gian hiện tại nhỏ hơn thời gian dựa trên period
			delay := time.Duration(expectedTime-currentTime) * time.Second // Tính toán thời gian chờ
			
			// Tối đa 30 giây
			maxWait := 30 * time.Second 
			
			// bỏ qua validator nếu thời gian chờ lớn hơn 30 giây
			select {
               case <-time.After(delay):
               case <-time.After(maxWait):
                  	log.Println("Validator timeout, skipping block", "validator", block.Coinbase.Hex()) // skip validator
               case <-stop:
				return errors.New("sealing stopped") // Lỗi khi dừng việc tạo block
			}
		}
	}
	
	// Gửi block
	select {
	case results <- block:
		return nil // Trả về nil
	case <-stop:
		return errors.New("sealing stopped") // Lỗi khi dừng việc tạo block	
	}
}

// Lấy hash của block
func (p *POVA) SealHash(header *types.Header) common.Hash {
	return header.Hash() // Trả về hash của header
}

// Tính độ khó của block
func (p *POVA) CalcDifficulty(chain consensus.ChainHeaderReader, time uint64, parent *types.Header) *big.Int {
	return big.NewInt(1) // Đặt độ khó của block là 1
}

// Lấy API của block
func (p *POVA) APIs(chain consensus.ChainHeaderReader) []rpc.API {
	return []rpc.API{} // Trả về API của block
}

// Đóng POVA
func (p *POVA) Close() error {
	return nil // Trả về nil
}
