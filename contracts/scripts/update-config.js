const fs = require('fs');
const path = require('path');
require('dotenv').config();

/**
 * Updates the contract-config.json file with deployed contract addresses
 * @param {string} network - Network name (e.g., 'hederaTestnet')
 * @param {object} contracts - Object containing contract addresses
 */
async function updateConfig(network, contracts) {
  const configPath = path.join(__dirname, '..', 'contract-config.json');
  
  // Read existing config
  let config = {};
  if (fs.existsSync(configPath)) {
    const configData = fs.readFileSync(configPath, 'utf8');
    config = JSON.parse(configData);
  } else {
    config = { networks: {} };
  }

  // Update config with new contract addresses
  config.networks[network] = {
    ...config.networks[network],
    ...contracts,
    deployedAt: new Date().toISOString(),
    deployer: process.env.EVM_ADDRESS,
    accountId: process.env.ACCOUNT_ID
  };

  // Write updated config back to file
  fs.writeFileSync(configPath, JSON.stringify(config, null, 2));
  console.log(`Updated contract-config.json with ${network} contract addresses`);
}

// Example usage:
// updateConfig('hederaTestnet', {
//   vcRegistry: '0x123...',
//   escrowContract: '0x456...',
//   verifierMarketplace: '0x789...'
// });

module.exports = { updateConfig };