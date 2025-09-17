// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

import "@openzeppelin/contracts-upgradeable/access/AccessControlUpgradeable.sol";
import "@openzeppelin/contracts-upgradeable/security/PausableUpgradeable.sol";
import "@openzeppelin/contracts-upgradeable/proxy/utils/Initializable.sol";
import "@openzeppelin/contracts-upgradeable/proxy/utils/UUPSUpgradeable.sol";

/**
 * @title VCRegistry
 * @dev Smart contract for anchoring and managing Verifiable Credentials on-chain
 * Stores VC hashes, issuer information, and revocation status
 */
contract VCRegistry is Initializable, AccessControlUpgradeable, PausableUpgradeable, UUPSUpgradeable {
    bytes32 public constant ISSUER_ROLE = keccak256("ISSUER_ROLE");
    bytes32 public constant UPGRADER_ROLE = keccak256("UPGRADER_ROLE");
    
    struct CredentialRecord {
        address issuer;
        uint64 issuedAt;
        bool revoked;
        uint64 revokedAt;
        string uri; // Off-chain status list URI
    }
    
    // Mapping from VC hash to credential record
    mapping(bytes32 => CredentialRecord) public credentials;
    
    // Mapping to track authorized issuers
    mapping(address => bool) public authorizedIssuers;
    
    // Events
    event Anchored(
        bytes32 indexed vcHash,
        address indexed issuer,
        uint64 issuedAt,
        string uri
    );
    
    event Revoked(
        bytes32 indexed vcHash,
        address indexed issuer,
        uint64 revokedAt
    );
    
    event IssuerAdded(address indexed issuer);
    event IssuerRemoved(address indexed issuer);
    
    /// @custom:oz-upgrades-unsafe-allow constructor
    constructor() {
        _disableInitializers();
    }
    
    function initialize(address admin) public initializer {
        __AccessControl_init();
        __Pausable_init();
        __UUPSUpgradeable_init();
        
        _grantRole(DEFAULT_ADMIN_ROLE, admin);
        _grantRole(UPGRADER_ROLE, admin);
    }
    
    /**
     * @dev Anchor a Verifiable Credential on-chain
     * @param vcHash The SHA256 hash of the normalized VC
     * @param uri The URI for off-chain status checking
     */
    function anchor(bytes32 vcHash, string calldata uri) 
        external 
        whenNotPaused 
        onlyRole(ISSUER_ROLE) 
    {
        require(vcHash != bytes32(0), "VCRegistry: Invalid VC hash");
        require(bytes(uri).length > 0, "VCRegistry: URI cannot be empty");
        require(credentials[vcHash].issuer == address(0), "VCRegistry: VC already anchored");
        
        credentials[vcHash] = CredentialRecord({
            issuer: msg.sender,
            issuedAt: uint64(block.timestamp),
            revoked: false,
            revokedAt: 0,
            uri: uri
        });
        
        emit Anchored(vcHash, msg.sender, uint64(block.timestamp), uri);
    }
    
    /**
     * @dev Revoke a Verifiable Credential
     * @param vcHash The SHA256 hash of the VC to revoke
     */
    function revoke(bytes32 vcHash) 
        external 
        whenNotPaused 
    {
        CredentialRecord storage record = credentials[vcHash];
        require(record.issuer != address(0), "VCRegistry: VC not found");
        require(record.issuer == msg.sender || hasRole(DEFAULT_ADMIN_ROLE, msg.sender), 
                "VCRegistry: Only issuer or admin can revoke");
        require(!record.revoked, "VCRegistry: VC already revoked");
        
        record.revoked = true;
        record.revokedAt = uint64(block.timestamp);
        
        emit Revoked(vcHash, record.issuer, uint64(block.timestamp));
    }
    
    /**
     * @dev Check if a VC is valid (anchored and not revoked)
     * @param vcHash The SHA256 hash of the VC
     * @return valid True if VC is anchored and not revoked
     * @return revoked True if VC is revoked
     * @return issuedAt Timestamp when VC was issued
     * @return revokedAt Timestamp when VC was revoked (0 if not revoked)
     * @return uri Off-chain status URI
     */
    function isValid(bytes32 vcHash) 
        external 
        view 
        returns (
            bool valid,
            bool revoked,
            uint64 issuedAt,
            uint64 revokedAt,
            string memory uri
        ) 
    {
        CredentialRecord memory record = credentials[vcHash];
        
        if (record.issuer == address(0)) {
            return (false, false, 0, 0, "");
        }
        
        return (
            !record.revoked,
            record.revoked,
            record.issuedAt,
            record.revokedAt,
            record.uri
        );
    }
    
    /**
     * @dev Register a new issuer
     * @param issuer Address of the issuer to register
     */
    function registerIssuer(address issuer) 
        external 
        onlyRole(DEFAULT_ADMIN_ROLE) 
    {
        require(issuer != address(0), "VCRegistry: Invalid issuer address");
        
        _grantRole(ISSUER_ROLE, issuer);
        authorizedIssuers[issuer] = true;
        
        emit IssuerAdded(issuer);
    }
    
    /**
     * @dev Revoke issuer privileges
     * @param issuer Address of the issuer to revoke
     */
    function revokeIssuer(address issuer) 
        external 
        onlyRole(DEFAULT_ADMIN_ROLE) 
    {
        require(issuer != address(0), "VCRegistry: Invalid issuer address");
        
        _revokeRole(ISSUER_ROLE, issuer);
        authorizedIssuers[issuer] = false;
        
        emit IssuerRemoved(issuer);
    }
    
    /**
     * @dev Pause the contract
     */
    function pause() external onlyRole(DEFAULT_ADMIN_ROLE) {
        _pause();
    }
    
    /**
     * @dev Unpause the contract
     */
    function unpause() external onlyRole(DEFAULT_ADMIN_ROLE) {
        _unpause();
    }
    
    /**
     * @dev Get credential record by hash
     * @param vcHash The SHA256 hash of the VC
     * @return The credential record
     */
    function getCredential(bytes32 vcHash) 
        external 
        view 
        returns (CredentialRecord memory) 
    {
        return credentials[vcHash];
    }
    
    /**
     * @dev Check if an address is an authorized issuer
     * @param issuer Address to check
     * @return True if the address is an authorized issuer
     */
    function isAuthorizedIssuer(address issuer) 
        external 
        view 
        returns (bool) 
    {
        return hasRole(ISSUER_ROLE, issuer);
    }
    
    function _authorizeUpgrade(address newImplementation)
        internal
        onlyRole(UPGRADER_ROLE)
        override
    {}
}