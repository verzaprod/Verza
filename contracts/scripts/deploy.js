const { ethers, upgrades } = require("hardhat");
const hre = require("hardhat");
const { updateConfig } = require("./update-config");

async function main() {
  const [deployer] = await ethers.getSigners();
  const networkName = hre.network.name;
  const contracts = {};

  console.log("Deploying contracts with the account:", deployer.address);
  console.log("Deploying to network:", networkName);

  // Deploy VCRegistry as upgradeable proxy
  const VCRegistry = await ethers.getContractFactory("VCRegistry");
  
  console.log("Deploying VCRegistry...");
  const vcRegistry = await upgrades.deployProxy(
    VCRegistry,
    [deployer.address, "VerzaVC", "VVC"], // admin address, name, symbol
    {
      initializer: "initialize",
      kind: "uups",
    }
  );

  // Wait for deployment transaction to be mined
  await vcRegistry.waitForDeployment();

  console.log("VCRegistry deployed to:", await vcRegistry.getAddress());
  console.log("Implementation address:", await upgrades.erc1967.getImplementationAddress(await vcRegistry.getAddress()));
  console.log("Admin address:", await upgrades.erc1967.getAdminAddress(await vcRegistry.getAddress()));

  // Register the deployer as an issuer for testing
  console.log("Registering deployer as issuer...");
  const registerTx = await vcRegistry.registerIssuer(deployer.address);
  await registerTx.wait();
  console.log("Deployer registered as issuer");

  // Verify deployment
  const isAuthorized = await vcRegistry.isAuthorizedIssuer(deployer.address);
  console.log("Deployer is authorized issuer:", isAuthorized);
  
  // Save contract address to config
  contracts.vcRegistry = await vcRegistry.getAddress();

  // Deploy VerifierMarketplace as upgradeable proxy
  const VerifierMarketplace = await ethers.getContractFactory("VerifierMarketplace");
  
  console.log("Deploying VerifierMarketplace...");
  const verifierMarketplace = await upgrades.deployProxy(
    VerifierMarketplace,
    [
      deployer.address, // admin
      ethers.parseEther("100"), // minimum stake (100 tokens)
      ethers.parseEther("10"), // base verification fee (10 tokens)
      ethers.ZeroAddress // staking token (use native token)
    ],
    {
      initializer: "initialize",
      kind: "uups",
    }
  );

  await verifierMarketplace.waitForDeployment();
  console.log("VerifierMarketplace deployed to:", await verifierMarketplace.getAddress());
  contracts.verifierMarketplace = await verifierMarketplace.getAddress();

  // Deploy EscrowContract as upgradeable proxy
  const EscrowContract = await ethers.getContractFactory("EscrowContract");
  
  console.log("Deploying EscrowContract...");
  const escrowContract = await upgrades.deployProxy(
    EscrowContract,
    [
      deployer.address, // admin
      await verifierMarketplace.getAddress(), // verifier marketplace
      ethers.ZeroAddress, // fraud detection (placeholder)
      ethers.ZeroAddress, // payment token (use native token)
      deployer.address // fee recipient
    ],
    {
      initializer: "initialize",
      kind: "uups",
    }
  );

  await escrowContract.waitForDeployment();
  console.log("EscrowContract deployed to:", await escrowContract.getAddress());
  contracts.escrowContract = await escrowContract.getAddress();

  // Save deployment info
  const deploymentInfo = {
    network: hre.network.name,
    chainId: (await ethers.provider.getNetwork()).chainId,
    deployer: deployer.address,
    blockNumber: (await ethers.provider.getBlockNumber()),
    timestamp: new Date().toISOString(),
    contracts: {
      vcRegistry: await vcRegistry.getAddress(),
      verifierMarketplace: await verifierMarketplace.getAddress(),
      escrowContract: await escrowContract.getAddress()
    }
  };

  // Try to get implementation addresses (may fail on Hedera testnet)
  try {
    deploymentInfo.vcRegistryImplementation = await upgrades.erc1967.getImplementationAddress(await vcRegistry.getAddress());
    deploymentInfo.verifierMarketplaceImplementation = await upgrades.erc1967.getImplementationAddress(await verifierMarketplace.getAddress());
    deploymentInfo.escrowContractImplementation = await upgrades.erc1967.getImplementationAddress(await escrowContract.getAddress());
  } catch (error) {
    console.log("Note: Could not retrieve implementation addresses (Hedera testnet limitation)");
  }

  console.log("\nDeployment Summary:");
  console.log("VCRegistry:", await vcRegistry.getAddress());
  console.log("VerifierMarketplace:", await verifierMarketplace.getAddress());
  console.log("EscrowContract:", await escrowContract.getAddress());
  
  // Update config file with deployed contract addresses
  await updateConfig(networkName, contracts);
  console.log("Contract addresses saved to config file");
  console.log(JSON.stringify(deploymentInfo, (key, value) => 
    typeof value === 'bigint' ? value.toString() : value, 2));

  // Save to file
  const fs = require("fs");
  const path = require("path");
  const deploymentsDir = path.join(__dirname, "..", "deployments");
  
  if (!fs.existsSync(deploymentsDir)) {
    fs.mkdirSync(deploymentsDir, { recursive: true });
  }
  
  const deploymentFile = path.join(deploymentsDir, `${hre.network.name}.json`);
  fs.writeFileSync(deploymentFile, JSON.stringify(deploymentInfo, (key, value) => 
    typeof value === 'bigint' ? value.toString() : value, 2));
  
  console.log(`\nDeployment info saved to: ${deploymentFile}`);
}

main()
  .then(() => process.exit(0))
  .catch((error) => {
    console.error(error);
    process.exit(1);
  });