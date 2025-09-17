const fs = require('fs');
const path = require('path');
const { execSync } = require('child_process');

async function generateBindings() {
  try {
    console.log('Generating Go bindings for VCRegistry contract...');
    
    // Paths
    const artifactsDir = path.join(__dirname, '..', 'artifacts', 'contracts', 'VCRegistry.sol');
    const abiPath = path.join(artifactsDir, 'VCRegistry.json');
    const binPath = path.join(artifactsDir, 'VCRegistry.bin');
    const outputDir = path.join(__dirname, '..', '..', 'pkg', 'contracts');
    const outputFile = path.join(outputDir, 'vcregistry.go');
    
    // Ensure output directory exists
    if (!fs.existsSync(outputDir)) {
      fs.mkdirSync(outputDir, { recursive: true });
    }
    
    // Check if artifacts exist
    if (!fs.existsSync(abiPath)) {
      throw new Error(`ABI file not found: ${abiPath}. Please compile contracts first.`);
    }
    
    // Read and extract ABI
    const artifact = JSON.parse(fs.readFileSync(abiPath, 'utf8'));
    const abi = JSON.stringify(artifact.abi);
    const bytecode = artifact.bytecode;
    
    // Write ABI to temporary file
    const tempAbiPath = path.join(__dirname, 'temp_abi.json');
    fs.writeFileSync(tempAbiPath, abi);
    
    // Write bytecode to bin file
    fs.writeFileSync(binPath, bytecode.replace('0x', ''));
    
    // Generate Go bindings using abigen
    const abigenCmd = `abigen --abi ${tempAbiPath} --bin ${binPath} --pkg contracts --type VCRegistry --out ${outputFile}`;
    
    console.log('Running abigen command:', abigenCmd);
    execSync(abigenCmd, { stdio: 'inherit' });
    
    // Clean up temporary files
    fs.unlinkSync(tempAbiPath);
    
    console.log(`Go bindings generated successfully: ${outputFile}`);
    
    // Generate additional helper file
    const helperContent = `package contracts

import (
	"context"
	"crypto/ecdsa"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

// VCRegistryClient wraps the generated VCRegistry contract with helper methods
type VCRegistryClient struct {
	*VCRegistry
	client   *ethclient.Client
	auth     *bind.TransactOpts
	address  common.Address
}

// NewVCRegistryClient creates a new VCRegistry client
func NewVCRegistryClient(client *ethclient.Client, contractAddress common.Address, privateKey *ecdsa.PrivateKey) (*VCRegistryClient, error) {
	contract, err := NewVCRegistry(contractAddress, client)
	if err != nil {
		return nil, err
	}

	chainID, err := client.NetworkID(context.Background())
	if err != nil {
		return nil, err
	}

	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, chainID)
	if err != nil {
		return nil, err
	}

	return &VCRegistryClient{
		VCRegistry: contract,
		client:     client,
		auth:       auth,
		address:    contractAddress,
	}, nil
}

// AnchorVC anchors a verifiable credential on-chain
func (c *VCRegistryClient) AnchorVC(ctx context.Context, vcHash [32]byte, statusURI string) (*types.Transaction, error) {
	return c.Anchor(c.auth, vcHash, statusURI)
}

// RevokeVC revokes a verifiable credential on-chain
func (c *VCRegistryClient) RevokeVC(ctx context.Context, vcHash [32]byte) (*types.Transaction, error) {
	return c.Revoke(c.auth, vcHash)
}

// IsVCValid checks if a verifiable credential is valid
func (c *VCRegistryClient) IsVCValid(ctx context.Context, vcHash [32]byte) (bool, bool, *big.Int, *big.Int, string, error) {
	return c.IsValid(&bind.CallOpts{Context: ctx}, vcHash)
}

// GetVCCredential retrieves credential information
func (c *VCRegistryClient) GetVCCredential(ctx context.Context, vcHash [32]byte) (VCRegistryCredential, error) {
	return c.GetCredential(&bind.CallOpts{Context: ctx}, vcHash)
}

// IsAuthorizedIssuerCheck checks if an address is an authorized issuer
func (c *VCRegistryClient) IsAuthorizedIssuerCheck(ctx context.Context, issuer common.Address) (bool, error) {
	return c.IsAuthorizedIssuer(&bind.CallOpts{Context: ctx}, issuer)
}

// SetGasPrice sets the gas price for transactions
func (c *VCRegistryClient) SetGasPrice(gasPrice *big.Int) {
	c.auth.GasPrice = gasPrice
}

// SetGasLimit sets the gas limit for transactions
func (c *VCRegistryClient) SetGasLimit(gasLimit uint64) {
	c.auth.GasLimit = gasLimit
}

// GetContractAddress returns the contract address
func (c *VCRegistryClient) GetContractAddress() common.Address {
	return c.address
}

// GetAuth returns the transaction auth
func (c *VCRegistryClient) GetAuth() *bind.TransactOpts {
	return c.auth
}
`;
    
    const helperPath = path.join(outputDir, 'client.go');
    fs.writeFileSync(helperPath, helperContent);
    
    console.log(`Helper client generated: ${helperPath}`);
    
  } catch (error) {
    console.error('Error generating bindings:', error.message);
    process.exit(1);
  }
}

if (require.main === module) {
  generateBindings();
}

module.exports = { generateBindings };