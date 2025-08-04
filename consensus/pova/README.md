# POVA (Proof of Validator Authority) Consensus Engine

## Overview

POVA is a custom consensus mechanism for Go Ethereum that implements a deterministic validator rotation system. Validators take turns creating blocks in a fixed order based on the formula: `(blockNumber-1) % len(validators)`.

## Features

- **Deterministic Validator Rotation**: Validators rotate in a predictable order
- **Configurable Block Time**: Adjustable time intervals between blocks (default: 15 seconds)
- **Pre-merge Mode**: Operates as proof-of-work consensus, no beacon client required
- **State Management**: Proper state root calculation and commitment
- **Full Ethereum API Compatibility**: Supports all standard Ethereum RPC endpoints

## Architecture

### Core Components

- **Validator Selection**: Uses block number to determine which validator should create the next block
- **Block Time Control**: Implements delays to ensure consistent block intervals
- **State Management**: Automatically commits state and sets state root in block headers
- **Header Verification**: Validates that blocks are created by authorized validators

### Consensus Algorithm

1. **Block Proposal**: Validator is selected based on `(blockNumber-1) % len(validators)`
2. **Header Preparation**: Sets coinbase, difficulty, and timestamp
3. **Block Sealing**: Implements time delays to maintain block intervals
4. **State Finalization**: Commits state and calculates state root
5. **Verification**: Other nodes verify block was created by correct validator

## Implementation

### Key Methods

#### `Prepare(chain, header)`
- Sets the correct validator as coinbase address
- Calculates block timestamp based on configured period
- Sets difficulty to 1 (constant for POVA)

#### `Seal(chain, block, results, stop)`
- Implements block time control with configurable delays
- Ensures minimum time between consecutive blocks
- Handles graceful shutdown via stop channel

#### `FinalizeAndAssemble(chain, header, state, txs, uncles, receipts, withdrawals)`
- Commits state to database
- Calculates and sets state root in block header
- Creates block with proper state information

#### `VerifyHeader(chain, header)`
- Validates block was created by authorized validator
- Uses deterministic formula for validator selection
- Returns error for unauthorized block proposers

## Configuration

### POVAConfig Structure

```go
type POVAConfig struct {
    Validators []common.Address `json:"validators"` // List of validator addresses
    Period     uint64           `json:"period"`     // Block time in seconds
}
```

### Genesis Configuration

```json
{
    "config": {
        "chainId": 12345,
        "homesteadBlock": 0,
        "eip150Block": 0,
        "eip155Block": 0,
        "eip158Block": 0,
        "byzantiumBlock": 0,
        "constantinopleBlock": 0,
        "petersburgBlock": 0,
        "istanbulBlock": 0,
        "berlinBlock": 0,
        "londonBlock": 0,
        "terminalTotalDifficulty": 1000000000000000000000000,
        "pova": {
            "validators": [
                "0x772cA2FDAe794772abD1fE776E073125e0C60360",
                "0x75239202dF072612fF56471437f7C0de47e3fF0C"
            ],
            "period": 15
        }
    },
    "difficulty": "0x1",
    "gasLimit": "0x8000000",
    "alloc": {
        "0x772cA2FDAe794772abD1fE776E073125e0C60360": {
            "balance": "0x200000000000000000000000000000000000000000000000000000000000000"
        },
        "0x75239202dF072612fF56471437f7C0de47e3fF0C": {
            "balance": "0x200000000000000000000000000000000000000000000000000000000000000"
        }
    },
    "coinbase": "0x772cA2FDAe794772abD1fE776E073125e0C60360",
    "extraData": "0x0000000000000000000000000000000000000000000000000000000000000000",
    "nonce": "0x0000000000000042",
    "mixhash": "0x0000000000000000000000000000000000000000000000000000000000000000",
    "parentHash": "0x0000000000000000000000000000000000000000000000000000000000000000",
    "timestamp": "0x00"
}
```

## Usage

### Building

```bash
# Build Go Ethereum with POVA
go build -buildvcs=false -o mychain.exe ./cmd/geth
```

### Initialization

```bash
# Initialize blockchain with genesis
./mychain.exe --datadir pova_test init test/pova_genesis.json
```

### Running

```bash
# Start node with mining
./mychain.exe --datadir pova_test \
    --mine \
    --miner.etherbase=0x772cA2FDAe794772abD1fE776E073125e0C60360 \
    --http \
    --http.addr=0.0.0.0 \
    --http.port=8545 \
    --http.api=eth,net,web3,personal,miner,admin,debug \
    --allow-insecure-unlock \
    --networkid=99999 \
    --verbosity=3
```

### Testing

```bash
# Run unit tests
go test ./consensus/pova

# Run integration tests
./test/test_pova.sh  # Linux/Mac
.\test\test_pova.ps1 # Windows
```

## API Endpoints

The POVA node provides standard Ethereum RPC endpoints:

- **HTTP**: `http://localhost:8545`
- **WebSocket**: `ws://localhost:8546`
- **IPC**: `\\.\pipe\geth.ipc`

### Example API Calls

```bash
# Get block number
curl -X POST -H "Content-Type: application/json" \
  --data '{"jsonrpc":"2.0","method":"eth_blockNumber","params":[],"id":1}' \
  http://localhost:8545

# Get latest block
curl -X POST -H "Content-Type: application/json" \
  --data '{"jsonrpc":"2.0","method":"eth_getBlockByNumber","params":["latest", true],"id":1}' \
  http://localhost:8545
```

## Validator Rotation

The validator selection follows a deterministic pattern:

- **Block 1**: `(1-1) % 2 = 0` â†’ Validator 1
- **Block 2**: `(2-1) % 2 = 1` â†’ Validator 2
- **Block 3**: `(3-1) % 2 = 0` â†’ Validator 1
- **Block 4**: `(4-1) % 2 = 1` â†’ Validator 2
- And so on...

## Performance Characteristics

- **Block Time**: Configurable (default: 15 seconds)
- **Validator Rotation**: Deterministic and predictable
- **State Management**: Efficient with proper state root calculation
- **API Response**: Fast with standard Ethereum RPC

## Security Considerations

### Advantages
- **Deterministic**: Predictable validator rotation
- **Fast**: No complex proof-of-work calculations
- **Simple**: Easy to understand and implement
- **Configurable**: Flexible block time and validator set

### Limitations
- **Centralized**: Requires trusted validators
- **No Slashing**: No mechanism to punish malicious validators
- **Fixed Set**: Validator set is static (can be extended)
- **Not Suitable for Mainnet**: Lacks economic security of PoW/PoS

## Troubleshooting

### Common Issues

#### "Zero state root hash"
- **Cause**: State not properly committed
- **Solution**: Fixed in `FinalizeAndAssemble` method

#### "Post-merge network, but no beacon client"
- **Cause**: Incorrect terminal difficulty configuration
- **Solution**: Set high `terminalTotalDifficulty` and remove `terminalTotalDifficultyPassed`

#### "Unauthorized block proposer"
- **Cause**: Block created by wrong validator
- **Solution**: Check validator addresses and rotation logic

#### Node not mining blocks
- **Cause**: Missing etherbase or incorrect configuration
- **Solution**: Verify etherbase address and genesis configuration

### Debug Mode

```bash
# Run with verbose logging
./mychain.exe --datadir pova_test --verbosity=4 --log.vmodule="consensus/*=5" --mine
```

## Development

### Adding New Features

1. **Dynamic Validator Set**: Implement validator set updates
2. **Slashing Mechanism**: Add penalties for malicious behavior
3. **Validator Weighting**: Support weighted validator selection
4. **Fork Choice**: Implement fork choice rules

### Testing

```bash
# Run all tests
go test ./consensus/pova -v

# Run specific test
go test ./consensus/pova -run TestPOVA_ValidatorRotation
```

## License

This implementation follows the same license as Go Ethereum.

## Contributing

1. Fork the repository
2. Create a feature branch
3. Implement changes with tests
4. Submit a pull request

## References

- [Go Ethereum Documentation](https://geth.ethereum.org/docs/)
- [Ethereum Consensus Documentation](https://ethereum.org/en/developers/docs/consensus-mechanisms/)
- [Consensus Engine Interface](https://github.com/ethereum/go-ethereum/blob/master/consensus/consensus.go) 

