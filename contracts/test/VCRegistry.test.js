const { expect } = require("chai");
const { ethers, upgrades } = require("hardhat");

describe("VCRegistry", function () {
  let VCRegistry;
  let vcRegistry;
  let admin;
  let issuer1;
  let issuer2;
  let user;
  let addrs;

  const vcHash1 = ethers.utils.keccak256(ethers.utils.toUtf8Bytes("test-vc-1"));
  const vcHash2 = ethers.utils.keccak256(ethers.utils.toUtf8Bytes("test-vc-2"));
  const statusUri = "https://example.com/status/123";

  beforeEach(async function () {
    [admin, issuer1, issuer2, user, ...addrs] = await ethers.getSigners();

    VCRegistry = await ethers.getContractFactory("VCRegistry");
    vcRegistry = await upgrades.deployProxy(
      VCRegistry,
      [admin.address],
      {
        initializer: "initialize",
        kind: "uups",
      }
    );
    await vcRegistry.deployed();
  });

  describe("Deployment", function () {
    it("Should set the right admin", async function () {
      expect(await vcRegistry.hasRole(await vcRegistry.DEFAULT_ADMIN_ROLE(), admin.address)).to.be.true;
    });

    it("Should not be paused initially", async function () {
      expect(await vcRegistry.paused()).to.be.false;
    });
  });

  describe("Issuer Management", function () {
    it("Should register issuer by admin", async function () {
      await expect(vcRegistry.connect(admin).registerIssuer(issuer1.address))
        .to.emit(vcRegistry, "IssuerAdded")
        .withArgs(issuer1.address);

      expect(await vcRegistry.isAuthorizedIssuer(issuer1.address)).to.be.true;
      expect(await vcRegistry.authorizedIssuers(issuer1.address)).to.be.true;
    });

    it("Should revoke issuer by admin", async function () {
      // First register issuer
      await vcRegistry.connect(admin).registerIssuer(issuer1.address);
      expect(await vcRegistry.isAuthorizedIssuer(issuer1.address)).to.be.true;

      // Then revoke
      await expect(vcRegistry.connect(admin).revokeIssuer(issuer1.address))
        .to.emit(vcRegistry, "IssuerRemoved")
        .withArgs(issuer1.address);

      expect(await vcRegistry.isAuthorizedIssuer(issuer1.address)).to.be.false;
      expect(await vcRegistry.authorizedIssuers(issuer1.address)).to.be.false;
    });

    it("Should not allow non-admin to register issuer", async function () {
      await expect(
        vcRegistry.connect(user).registerIssuer(issuer1.address)
      ).to.be.revertedWith("AccessControl: account");
    });

    it("Should not register zero address as issuer", async function () {
      await expect(
        vcRegistry.connect(admin).registerIssuer(ethers.constants.AddressZero)
      ).to.be.revertedWith("VCRegistry: Invalid issuer address");
    });
  });

  describe("VC Anchoring", function () {
    beforeEach(async function () {
      // Register issuer1
      await vcRegistry.connect(admin).registerIssuer(issuer1.address);
    });

    it("Should anchor VC by authorized issuer", async function () {
      await expect(
        vcRegistry.connect(issuer1).anchor(vcHash1, statusUri)
      )
        .to.emit(vcRegistry, "Anchored")
        .withArgs(vcHash1, issuer1.address, await getBlockTimestamp(), statusUri);

      const credential = await vcRegistry.getCredential(vcHash1);
      expect(credential.issuer).to.equal(issuer1.address);
      expect(credential.revoked).to.be.false;
      expect(credential.uri).to.equal(statusUri);
    });

    it("Should not anchor VC by unauthorized issuer", async function () {
      await expect(
        vcRegistry.connect(issuer2).anchor(vcHash1, statusUri)
      ).to.be.revertedWith("AccessControl: account");
    });

    it("Should not anchor VC with zero hash", async function () {
      await expect(
        vcRegistry.connect(issuer1).anchor(ethers.constants.HashZero, statusUri)
      ).to.be.revertedWith("VCRegistry: Invalid VC hash");
    });

    it("Should not anchor VC with empty URI", async function () {
      await expect(
        vcRegistry.connect(issuer1).anchor(vcHash1, "")
      ).to.be.revertedWith("VCRegistry: URI cannot be empty");
    });

    it("Should not anchor same VC twice", async function () {
      await vcRegistry.connect(issuer1).anchor(vcHash1, statusUri);
      
      await expect(
        vcRegistry.connect(issuer1).anchor(vcHash1, statusUri)
      ).to.be.revertedWith("VCRegistry: VC already anchored");
    });

    it("Should not anchor when paused", async function () {
      await vcRegistry.connect(admin).pause();
      
      await expect(
        vcRegistry.connect(issuer1).anchor(vcHash1, statusUri)
      ).to.be.revertedWith("Pausable: paused");
    });
  });

  describe("VC Revocation", function () {
    beforeEach(async function () {
      // Register issuer1 and anchor a VC
      await vcRegistry.connect(admin).registerIssuer(issuer1.address);
      await vcRegistry.connect(issuer1).anchor(vcHash1, statusUri);
    });

    it("Should revoke VC by issuer", async function () {
      await expect(
        vcRegistry.connect(issuer1).revoke(vcHash1)
      )
        .to.emit(vcRegistry, "Revoked")
        .withArgs(vcHash1, issuer1.address, await getBlockTimestamp());

      const credential = await vcRegistry.getCredential(vcHash1);
      expect(credential.revoked).to.be.true;
      expect(credential.revokedAt).to.be.gt(0);
    });

    it("Should revoke VC by admin", async function () {
      await expect(
        vcRegistry.connect(admin).revoke(vcHash1)
      )
        .to.emit(vcRegistry, "Revoked")
        .withArgs(vcHash1, issuer1.address, await getBlockTimestamp());

      const credential = await vcRegistry.getCredential(vcHash1);
      expect(credential.revoked).to.be.true;
    });

    it("Should not revoke VC by unauthorized user", async function () {
      await expect(
        vcRegistry.connect(user).revoke(vcHash1)
      ).to.be.revertedWith("VCRegistry: Only issuer or admin can revoke");
    });

    it("Should not revoke non-existent VC", async function () {
      await expect(
        vcRegistry.connect(issuer1).revoke(vcHash2)
      ).to.be.revertedWith("VCRegistry: VC not found");
    });

    it("Should not revoke already revoked VC", async function () {
      await vcRegistry.connect(issuer1).revoke(vcHash1);
      
      await expect(
        vcRegistry.connect(issuer1).revoke(vcHash1)
      ).to.be.revertedWith("VCRegistry: VC already revoked");
    });

    it("Should not revoke when paused", async function () {
      await vcRegistry.connect(admin).pause();
      
      await expect(
        vcRegistry.connect(issuer1).revoke(vcHash1)
      ).to.be.revertedWith("Pausable: paused");
    });
  });

  describe("VC Validation", function () {
    beforeEach(async function () {
      await vcRegistry.connect(admin).registerIssuer(issuer1.address);
    });

    it("Should return valid status for anchored VC", async function () {
      await vcRegistry.connect(issuer1).anchor(vcHash1, statusUri);
      
      const [valid, revoked, issuedAt, revokedAt, uri] = await vcRegistry.isValid(vcHash1);
      
      expect(valid).to.be.true;
      expect(revoked).to.be.false;
      expect(issuedAt).to.be.gt(0);
      expect(revokedAt).to.equal(0);
      expect(uri).to.equal(statusUri);
    });

    it("Should return invalid status for revoked VC", async function () {
      await vcRegistry.connect(issuer1).anchor(vcHash1, statusUri);
      await vcRegistry.connect(issuer1).revoke(vcHash1);
      
      const [valid, revoked, issuedAt, revokedAt, uri] = await vcRegistry.isValid(vcHash1);
      
      expect(valid).to.be.false;
      expect(revoked).to.be.true;
      expect(issuedAt).to.be.gt(0);
      expect(revokedAt).to.be.gt(0);
      expect(uri).to.equal(statusUri);
    });

    it("Should return invalid status for non-existent VC", async function () {
      const [valid, revoked, issuedAt, revokedAt, uri] = await vcRegistry.isValid(vcHash2);
      
      expect(valid).to.be.false;
      expect(revoked).to.be.false;
      expect(issuedAt).to.equal(0);
      expect(revokedAt).to.equal(0);
      expect(uri).to.equal("");
    });
  });

  describe("Pausable", function () {
    it("Should pause and unpause by admin", async function () {
      await vcRegistry.connect(admin).pause();
      expect(await vcRegistry.paused()).to.be.true;
      
      await vcRegistry.connect(admin).unpause();
      expect(await vcRegistry.paused()).to.be.false;
    });

    it("Should not pause by non-admin", async function () {
      await expect(
        vcRegistry.connect(user).pause()
      ).to.be.revertedWith("AccessControl: account");
    });
  });

  describe("Upgradeability", function () {
    it("Should upgrade by upgrader role", async function () {
      const VCRegistryV2 = await ethers.getContractFactory("VCRegistry");
      
      await expect(
        upgrades.upgradeProxy(vcRegistry.address, VCRegistryV2)
      ).to.not.be.reverted;
    });
  });

  // Helper function to get current block timestamp
  async function getBlockTimestamp() {
    const blockNumber = await ethers.provider.getBlockNumber();
    const block = await ethers.provider.getBlock(blockNumber + 1);
    return block.timestamp;
  }
});