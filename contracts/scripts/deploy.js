const { ethers, upgrades } = require("hardhat");

async function main() {
  const [deployer] = await ethers.getSigners();

  console.log("Deploying contracts with the account:", deployer.address);
  console.log("Account balance:", (await deployer.getBalance()).toString());

  // Deploy VCRegistry as upgradeable proxy
  const VCRegistry = await ethers.getContractFactory("VCRegistry");
  
  console.log("Deploying VCRegistry...");
  const vcRegistry = await upgrades.deployProxy(
    VCRegistry,
    [deployer.address], // admin address
    {
      initializer: "initialize",
      kind: "uups",
    }
  );

  await vcRegistry.deployed();

  console.log("VCRegistry deployed to:", vcRegistry.address);
  console.log("Implementation address:", await upgrades.erc1967.getImplementationAddress(vcRegistry.address));
  console.log("Admin address:", await upgrades.erc1967.getAdminAddress(vcRegistry.address));

  // Register the deployer as an issuer for testing
  console.log("Registering deployer as issuer...");
  const registerTx = await vcRegistry.registerIssuer(deployer.address);
  await registerTx.wait();
  console.log("Deployer registered as issuer");

  // Verify deployment
  const isAuthorized = await vcRegistry.isAuthorizedIssuer(deployer.address);
  console.log("Deployer is authorized issuer:", isAuthorized);

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