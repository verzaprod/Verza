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

  // Save deployment info
  const deploymentInfo = {
    network: hre.network.name,
    chainId: (await ethers.provider.getNetwork()).chainId,
    contractAddress: vcRegistry.address,
    implementationAddress: await upgrades.erc1967.getImplementationAddress(vcRegistry.address),
    adminAddress: await upgrades.erc1967.getAdminAddress(vcRegistry.address),
    deployer: deployer.address,
    blockNumber: (await ethers.provider.getBlockNumber()),
    timestamp: new Date().toISOString(),
  };

  console.log("\nDeployment Summary:");
  
  // Save contract addresses to config
  contracts.vcRegistry = vcRegistry.address;
  
  // Update config file with deployed contract addresses
  await updateConfig(networkName, contracts);
  console.log("Contract addresses saved to config file");
  console.log(JSON.stringify(deploymentInfo, null, 2));

  // Save to file
  const fs = require("fs");
  const path = require("path");
  const deploymentsDir = path.join(__dirname, "..", "deployments");
  
  if (!fs.existsSync(deploymentsDir)) {
    fs.mkdirSync(deploymentsDir, { recursive: true });
  }
  
  const deploymentFile = path.join(deploymentsDir, `${hre.network.name}.json`);
  fs.writeFileSync(deploymentFile, JSON.stringify(deploymentInfo, null, 2));
  
  console.log(`\nDeployment info saved to: ${deploymentFile}`);
}

main()
  .then(() => process.exit(0))
  .catch((error) => {
    console.error(error);
    process.exit(1);
  });