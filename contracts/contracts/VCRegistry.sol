// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

import "@openzeppelin/contracts-upgradeable/access/AccessControlUpgradeable.sol";
import "@openzeppelin/contracts-upgradeable/security/PausableUpgradeable.sol";
import "@openzeppelin/contracts-upgradeable/proxy/utils/Initializable.sol";
import "@openzeppelin/contracts-upgradeable/proxy/utils/UUPSUpgradeable.sol";
import "@openzeppelin/contracts-upgradeable/token/ERC721/ERC721Upgradeable.sol";
import "@openzeppelin/contracts-upgradeable/token/ERC721/extensions/ERC721URIStorageUpgradeable.sol";
import "@openzeppelin/contracts-upgradeable/utils/CountersUpgradeable.sol";
import "@openzeppelin/contracts-upgradeable/security/ReentrancyGuardUpgradeable.sol";

/**
 * @title VCRegistry
 * @dev Enhanced smart contract for managing Verifiable Credentials as Soulbound Tokens (SBTs)
 * Supports Hedera DID integration, non-transferable NFTs, and advanced credential management
 */
contract VCRegistry is 
    Initializable, 
    AccessControlUpgradeable, 
    PausableUpgradeable, 
    UUPSUpgradeable,
    ERC721Upgradeable,
    ERC721URIStorageUpgradeable,
    ReentrancyGuardUpgradeable 
{
    using CountersUpgradeable for CountersUpgradeable.Counter;
    
    bytes32 public constant ISSUER_ROLE = keccak256("ISSUER_ROLE");
    bytes32 public constant UPGRADER_ROLE = keccak256("UPGRADER_ROLE");
    bytes32 public constant DID_RESOLVER_ROLE = keccak256("DID_RESOLVER_ROLE");
    
    enum CredentialType {
        Identity,
        Education,
        Professional,
        Financial,
        Health,
        Government,
        Custom
    }
    
    enum CredentialStatus {
        Active,
        Suspended,
        Revoked,
        Expired
    }
    
    struct CredentialRecord {
        uint256 tokenId;
        address issuer;
        address holder;
        string hederaDID; // Hedera DID of the credential holder
        bytes32 vcHash; // Hash of the verifiable credential
        CredentialType credentialType;
        CredentialStatus status;
        uint64 issuedAt;
        uint64 expiresAt;
        uint64 revokedAt;
        string metadataURI; // IPFS or Hedera File Service URI
        string schemaURI; // URI to credential schema
        bytes32[] claims; // Array of claim hashes
        mapping(string => string) attributes; // Dynamic attributes
    }
    
    struct DIDDocument {
        string did;
        address controller;
        string[] verificationMethods;
        string[] services;
        uint64 created;
        uint64 updated;
        bool active;
    }
    
    // Token ID counter
    CountersUpgradeable.Counter private _tokenIdCounter;
    
    // Mapping from token ID to credential record
    mapping(uint256 => CredentialRecord) public credentials;
    
    // Mapping from VC hash to token ID
    mapping(bytes32 => uint256) public vcHashToTokenId;
    
    // Mapping from Hedera DID to token IDs
    mapping(string => uint256[]) public didToTokenIds;
    
    // Mapping from Hedera DID to DID Document
    mapping(string => DIDDocument) public didDocuments;
    
    // Mapping to track authorized issuers
    mapping(address => bool) public authorizedIssuers;
    
    // Mapping from issuer to issued credential count
    mapping(address => uint256) public issuerCredentialCount;
    
    // Mapping from credential type to count
    mapping(CredentialType => uint256) public credentialTypeCount;
    
    // Mapping to track credential schemas
    mapping(string => bool) public approvedSchemas;
    
    // Contract configuration
    uint256 public maxCredentialsPerDID;
    uint256 public defaultExpirationPeriod;
    bool public transfersEnabled; // For emergency transfers only
    
    // Events
    event CredentialIssued(
        uint256 indexed tokenId,
        bytes32 indexed vcHash,
        address indexed issuer,
        address holder,
        string hederaDID,
        CredentialType credentialType,
        uint64 issuedAt,
        uint64 expiresAt
    );
    
    event CredentialRevoked(
        uint256 indexed tokenId,
        bytes32 indexed vcHash,
        address indexed issuer,
        uint64 revokedAt,
        string reason
    );
    
    event CredentialSuspended(
        uint256 indexed tokenId,
        address indexed issuer,
        string reason
    );
    
    event CredentialReactivated(
        uint256 indexed tokenId,
        address indexed issuer
    );
    
    event DIDRegistered(
        string indexed did,
        address indexed controller,
        uint64 timestamp
    );
    
    event DIDUpdated(
        string indexed did,
        address indexed controller,
        uint64 timestamp
    );
    
    event DIDDeactivated(
        string indexed did,
        address indexed controller,
        uint64 timestamp
    );
    
    event SchemaApproved(
        string indexed schemaURI,
        address indexed approver
    );
    
    event IssuerAdded(address indexed issuer);
    event IssuerRemoved(address indexed issuer);
    
    event TransferAttemptBlocked(
        address indexed from,
        address indexed to,
        uint256 indexed tokenId
    );
    
    /// @custom:oz-upgrades-unsafe-allow constructor
    constructor() {
        _disableInitializers();
    }
    
    function initialize(
        address admin,
        string memory name,
        string memory symbol
    ) public initializer {
        __AccessControl_init();
        __Pausable_init();
        __UUPSUpgradeable_init();
        __ERC721_init(name, symbol);
        __ERC721URIStorage_init();
        __ReentrancyGuard_init();
        
        _grantRole(DEFAULT_ADMIN_ROLE, admin);
        _grantRole(UPGRADER_ROLE, admin);
        _grantRole(DID_RESOLVER_ROLE, admin);
        
        // Set default configuration
        maxCredentialsPerDID = 10;
        defaultExpirationPeriod = 365 days;
        transfersEnabled = false;
        
        // Start token IDs from 1
        _tokenIdCounter.increment();
    }
    
    /**
     * @dev Register a Hedera DID
     * @param did The Hedera DID string
     * @param controller The controller address
     * @param verificationMethods Array of verification method URIs
     * @param services Array of service endpoint URIs
     */
    function registerDID(
        string calldata did,
        address controller,
        string[] calldata verificationMethods,
        string[] calldata services
    ) 
        external 
        whenNotPaused 
        onlyRole(DID_RESOLVER_ROLE) 
    {
        require(bytes(did).length > 0, "VCRegistry: Invalid DID");
        require(controller != address(0), "VCRegistry: Invalid controller");
        require(!didDocuments[did].active, "VCRegistry: DID already registered");
        
        didDocuments[did] = DIDDocument({
            did: did,
            controller: controller,
            verificationMethods: verificationMethods,
            services: services,
            created: uint64(block.timestamp),
            updated: uint64(block.timestamp),
            active: true
        });
        
        emit DIDRegistered(did, controller, uint64(block.timestamp));
    }
    
    /**
     * @dev Issue a Verifiable Credential as a Soulbound Token
     * @param holder The credential holder address
     * @param hederaDID The Hedera DID of the holder
     * @param vcHash The SHA256 hash of the normalized VC
     * @param credentialType The type of credential
     * @param metadataURI The URI for credential metadata
     * @param schemaURI The URI for credential schema
     * @param claims Array of claim hashes
     * @param expirationPeriod Expiration period in seconds (0 for no expiration)
     */
    function issueCredential(
        address holder,
        string calldata hederaDID,
        bytes32 vcHash,
        CredentialType credentialType,
        string calldata metadataURI,
        string calldata schemaURI,
        bytes32[] calldata claims,
        uint256 expirationPeriod
    ) 
        external 
        whenNotPaused 
        onlyRole(ISSUER_ROLE)
        nonReentrant
    {
        require(holder != address(0), "VCRegistry: Invalid holder address");
        require(bytes(hederaDID).length > 0, "VCRegistry: Invalid Hedera DID");
        require(vcHash != bytes32(0), "VCRegistry: Invalid VC hash");
        require(bytes(metadataURI).length > 0, "VCRegistry: Metadata URI cannot be empty");
        require(vcHashToTokenId[vcHash] == 0, "VCRegistry: VC already issued");
        require(didDocuments[hederaDID].active, "VCRegistry: DID not registered or inactive");
        require(didToTokenIds[hederaDID].length < maxCredentialsPerDID, "VCRegistry: Max credentials per DID exceeded");
        
        // Validate schema if provided
        if (bytes(schemaURI).length > 0) {
            require(approvedSchemas[schemaURI], "VCRegistry: Schema not approved");
        }
        
        uint256 tokenId = _tokenIdCounter.current();
        _tokenIdCounter.increment();
        
        // Calculate expiration
        uint64 expiresAt = 0;
        if (expirationPeriod > 0) {
            expiresAt = uint64(block.timestamp + expirationPeriod);
        } else if (defaultExpirationPeriod > 0) {
            expiresAt = uint64(block.timestamp + defaultExpirationPeriod);
        }
        
        // Create credential record
        CredentialRecord storage credential = credentials[tokenId];
        credential.tokenId = tokenId;
        credential.issuer = msg.sender;
        credential.holder = holder;
        credential.hederaDID = hederaDID;
        credential.vcHash = vcHash;
        credential.credentialType = credentialType;
        credential.status = CredentialStatus.Active;
        credential.issuedAt = uint64(block.timestamp);
        credential.expiresAt = expiresAt;
        credential.revokedAt = 0;
        credential.metadataURI = metadataURI;
        credential.schemaURI = schemaURI;
        credential.claims = claims;
        
        // Update mappings
        vcHashToTokenId[vcHash] = tokenId;
        didToTokenIds[hederaDID].push(tokenId);
        issuerCredentialCount[msg.sender]++;
        credentialTypeCount[credentialType]++;
        
        // Mint the soulbound token
        _safeMint(holder, tokenId);
        _setTokenURI(tokenId, metadataURI);
        
        emit CredentialIssued(
            tokenId,
            vcHash,
            msg.sender,
            holder,
            hederaDID,
            credentialType,
            uint64(block.timestamp),
            expiresAt
        );
    }
    
    /**
     * @dev Revoke a Verifiable Credential
     * @param tokenId The token ID of the credential to revoke
     * @param reason Reason for revocation
     */
    function revokeCredential(uint256 tokenId, string calldata reason) 
        external 
        whenNotPaused 
    {
        require(_exists(tokenId), "VCRegistry: Token does not exist");
        CredentialRecord storage credential = credentials[tokenId];
        require(
            credential.issuer == msg.sender || hasRole(DEFAULT_ADMIN_ROLE, msg.sender),
            "VCRegistry: Only issuer or admin can revoke"
        );
        require(credential.status != CredentialStatus.Revoked, "VCRegistry: Credential already revoked");
        
        credential.status = CredentialStatus.Revoked;
        credential.revokedAt = uint64(block.timestamp);
        
        emit CredentialRevoked(
            tokenId,
            credential.vcHash,
            credential.issuer,
            uint64(block.timestamp),
            reason
        );
    }
    
    /**
     * @dev Suspend a Verifiable Credential
     * @param tokenId The token ID of the credential to suspend
     * @param reason Reason for suspension
     */
    function suspendCredential(uint256 tokenId, string calldata reason) 
        external 
        whenNotPaused 
    {
        require(_exists(tokenId), "VCRegistry: Token does not exist");
        CredentialRecord storage credential = credentials[tokenId];
        require(
            credential.issuer == msg.sender || hasRole(DEFAULT_ADMIN_ROLE, msg.sender),
            "VCRegistry: Only issuer or admin can suspend"
        );
        require(credential.status == CredentialStatus.Active, "VCRegistry: Credential not active");
        
        credential.status = CredentialStatus.Suspended;
        
        emit CredentialSuspended(tokenId, credential.issuer, reason);
    }
    
    /**
     * @dev Reactivate a suspended Verifiable Credential
     * @param tokenId The token ID of the credential to reactivate
     */
    function reactivateCredential(uint256 tokenId) 
        external 
        whenNotPaused 
    {
        require(_exists(tokenId), "VCRegistry: Token does not exist");
        CredentialRecord storage credential = credentials[tokenId];
        require(
            credential.issuer == msg.sender || hasRole(DEFAULT_ADMIN_ROLE, msg.sender),
            "VCRegistry: Only issuer or admin can reactivate"
        );
        require(credential.status == CredentialStatus.Suspended, "VCRegistry: Credential not suspended");
        
        // Check if credential has expired
        if (credential.expiresAt > 0 && block.timestamp > credential.expiresAt) {
            credential.status = CredentialStatus.Expired;
        } else {
            credential.status = CredentialStatus.Active;
        }
        
        emit CredentialReactivated(tokenId, credential.issuer);
    }
    
    /**
     * @dev Set credential attributes
     * @param tokenId The token ID of the credential
     * @param key The attribute key
     * @param value The attribute value
     */
    function setCredentialAttribute(
        uint256 tokenId,
        string calldata key,
        string calldata value
    ) 
        external 
        whenNotPaused 
    {
        require(_exists(tokenId), "VCRegistry: Token does not exist");
        CredentialRecord storage credential = credentials[tokenId];
        require(
            credential.issuer == msg.sender || hasRole(DEFAULT_ADMIN_ROLE, msg.sender),
            "VCRegistry: Only issuer or admin can set attributes"
        );
        
        credential.attributes[key] = value;
    }
    
    /**
     * @dev Check if a credential is valid and active
     * @param tokenId The token ID of the credential
     * @return valid True if credential is active and not expired
     * @return status Current status of the credential
     * @return issuedAt Timestamp when credential was issued
     * @return expiresAt Timestamp when credential expires (0 if no expiration)
     * @return revokedAt Timestamp when credential was revoked (0 if not revoked)
     */
    function isCredentialValid(uint256 tokenId) 
        external 
        view 
        returns (
            bool valid,
            CredentialStatus status,
            uint64 issuedAt,
            uint64 expiresAt,
            uint64 revokedAt
        ) 
    {
        if (!_exists(tokenId)) {
            return (false, CredentialStatus.Revoked, 0, 0, 0);
        }
        
        CredentialRecord storage credential = credentials[tokenId];
        
        // Check if expired
        if (credential.expiresAt > 0 && block.timestamp > credential.expiresAt) {
            return (false, CredentialStatus.Expired, credential.issuedAt, credential.expiresAt, credential.revokedAt);
        }
        
        bool isValid = credential.status == CredentialStatus.Active;
        
        return (
            isValid,
            credential.status,
            credential.issuedAt,
            credential.expiresAt,
            credential.revokedAt
        );
    }
    
    /**
     * @dev Get credential by VC hash
     * @param vcHash The VC hash
     * @return tokenId The token ID (0 if not found)
     * @return exists Whether the credential exists
     */
    function getCredentialByHash(bytes32 vcHash) 
        external 
        view 
        returns (uint256 tokenId, bool exists) 
    {
        tokenId = vcHashToTokenId[vcHash];
        exists = tokenId != 0 && _exists(tokenId);
    }
    
    /**
     * @dev Get credentials by Hedera DID
     * @param hederaDID The Hedera DID
     * @return tokenIds Array of token IDs associated with the DID
     */
    function getCredentialsByDID(string calldata hederaDID) 
        external 
        view 
        returns (uint256[] memory tokenIds) 
    {
        return didToTokenIds[hederaDID];
    }
    
    /**
     * @dev Get credential attribute
     * @param tokenId The token ID
     * @param key The attribute key
     * @return value The attribute value
     */
    function getCredentialAttribute(uint256 tokenId, string calldata key) 
        external 
        view 
        returns (string memory value) 
    {
        require(_exists(tokenId), "VCRegistry: Token does not exist");
        return credentials[tokenId].attributes[key];
    }
    
    /**
     * @dev Get DID document
     * @param did The Hedera DID
     * @return document The DID document
     */
    function getDIDDocument(string calldata did) 
        external 
        view 
        returns (DIDDocument memory document) 
    {
        return didDocuments[did];
    }
    
    /**
     * @dev Override transfer functions to make tokens soulbound (non-transferable)
     */
    function _beforeTokenTransfer(
        address from,
        address to,
        uint256 tokenId,
        uint256 batchSize
    ) internal virtual override {
        super._beforeTokenTransfer(from, to, tokenId, batchSize);
        
        // Allow minting (from == address(0)) and burning (to == address(0))
        if (from != address(0) && to != address(0)) {
            // Only allow transfers if explicitly enabled (for emergency situations)
            if (!transfersEnabled) {
                emit TransferAttemptBlocked(from, to, tokenId);
                revert("VCRegistry: Soulbound tokens are non-transferable");
            }
            
            // Even if transfers are enabled, require admin approval
            require(hasRole(DEFAULT_ADMIN_ROLE, msg.sender), "VCRegistry: Transfer requires admin approval");
        }
    }
    
    /**
     * @dev Override approve to prevent approvals (soulbound tokens)
     */
    function approve(address to, uint256 tokenId) public virtual override(ERC721Upgradeable, IERC721Upgradeable) {
        revert("VCRegistry: Soulbound tokens cannot be approved");
    }
    
    /**
     * @dev Override setApprovalForAll to prevent approvals (soulbound tokens)
     */
    function setApprovalForAll(address operator, bool approved) public virtual override(ERC721Upgradeable, IERC721Upgradeable) {
        revert("VCRegistry: Soulbound tokens cannot be approved");
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
     * @dev Get credential record by token ID
     * @param tokenId The token ID
     * @return _tokenId The token ID
     * @return issuer The issuer address
     * @return holder The holder address
     * @return hederaDID The Hedera DID
     * @return vcHash The credential hash
     */
    function getCredential(uint256 tokenId) 
        external 
        view 
        returns (
            uint256 _tokenId,
            address issuer,
            address holder,
            string memory hederaDID,
            bytes32 vcHash,
            CredentialType credentialType,
            CredentialStatus status,
            uint64 issuedAt,
            uint64 expiresAt,
            uint64 revokedAt,
            string memory metadataURI,
            string memory schemaURI,
            bytes32[] memory claims
        ) 
    {
        require(_exists(tokenId), "VCRegistry: Token does not exist");
        CredentialRecord storage credential = credentials[tokenId];
        
        return (
            credential.tokenId,
            credential.issuer,
            credential.holder,
            credential.hederaDID,
            credential.vcHash,
            credential.credentialType,
            credential.status,
            credential.issuedAt,
            credential.expiresAt,
            credential.revokedAt,
            credential.metadataURI,
            credential.schemaURI,
            credential.claims
        );
    }
    
    /**
     * @dev Approve a credential schema
     * @param schemaURI The schema URI to approve
     */
    function approveSchema(string calldata schemaURI) 
        external 
        onlyRole(DEFAULT_ADMIN_ROLE) 
    {
        require(bytes(schemaURI).length > 0, "VCRegistry: Invalid schema URI");
        approvedSchemas[schemaURI] = true;
        
        emit SchemaApproved(schemaURI, msg.sender);
    }
    
    /**
     * @dev Revoke approval for a credential schema
     * @param schemaURI The schema URI to revoke
     */
    function revokeSchemaApproval(string calldata schemaURI) 
        external 
        onlyRole(DEFAULT_ADMIN_ROLE) 
    {
        approvedSchemas[schemaURI] = false;
    }
    
    /**
     * @dev Update DID document
     * @param did The Hedera DID
     * @param verificationMethods New verification methods
     * @param services New services
     */
    function updateDID(
        string calldata did,
        string[] calldata verificationMethods,
        string[] calldata services
    ) 
        external 
        whenNotPaused 
        onlyRole(DID_RESOLVER_ROLE) 
    {
        require(didDocuments[did].active, "VCRegistry: DID not registered or inactive");
        
        DIDDocument storage document = didDocuments[did];
        document.verificationMethods = verificationMethods;
        document.services = services;
        document.updated = uint64(block.timestamp);
        
        emit DIDUpdated(did, document.controller, uint64(block.timestamp));
    }
    
    /**
     * @dev Deactivate a DID
     * @param did The Hedera DID to deactivate
     */
    function deactivateDID(string calldata did) 
        external 
        whenNotPaused 
        onlyRole(DID_RESOLVER_ROLE) 
    {
        require(didDocuments[did].active, "VCRegistry: DID not active");
        
        didDocuments[did].active = false;
        didDocuments[did].updated = uint64(block.timestamp);
        
        emit DIDDeactivated(did, didDocuments[did].controller, uint64(block.timestamp));
    }
    
    /**
     * @dev Update contract configuration
     */
    function updateMaxCredentialsPerDID(uint256 _maxCredentialsPerDID) 
        external 
        onlyRole(DEFAULT_ADMIN_ROLE) 
    {
        maxCredentialsPerDID = _maxCredentialsPerDID;
    }
    
    function updateDefaultExpirationPeriod(uint256 _defaultExpirationPeriod) 
        external 
        onlyRole(DEFAULT_ADMIN_ROLE) 
    {
        defaultExpirationPeriod = _defaultExpirationPeriod;
    }
    
    function setTransfersEnabled(bool _transfersEnabled) 
        external 
        onlyRole(DEFAULT_ADMIN_ROLE) 
    {
        transfersEnabled = _transfersEnabled;
    }
    
    /**
     * @dev Get contract statistics
     */
    function getContractStats() 
        external 
        view 
        returns (
            uint256 totalCredentials,
            uint256 totalIssuers,
            uint256 totalDIDs,
            uint256[7] memory credentialTypeCounts
        ) 
    {
        totalCredentials = _tokenIdCounter.current() - 1;
        
        // Count active issuers
        // Note: This is a simplified count, in practice you might want to maintain a separate counter
        
        // Count active DIDs
        // Note: This would require maintaining a separate counter for efficiency
        
        // Get credential type counts
        for (uint i = 0; i < 7; i++) {
            credentialTypeCounts[i] = credentialTypeCount[CredentialType(i)];
        }
        
        return (totalCredentials, 0, 0, credentialTypeCounts); // Simplified return
    }
    
    /**
     * @dev Override tokenURI to support both ERC721URIStorage and custom logic
     */
    function tokenURI(uint256 tokenId) 
        public 
        view 
        virtual 
        override(ERC721Upgradeable, ERC721URIStorageUpgradeable) 
        returns (string memory) 
    {
        return ERC721URIStorageUpgradeable.tokenURI(tokenId);
    }
    
    /**
     * @dev Override supportsInterface
     */
    function supportsInterface(bytes4 interfaceId) 
        public 
        view 
        virtual 
        override(ERC721Upgradeable, AccessControlUpgradeable, ERC721URIStorageUpgradeable) 
        returns (bool) 
    {
        return super.supportsInterface(interfaceId);
    }
    
    /**
     * @dev Override _burn to handle URI storage
     */
    function _burn(uint256 tokenId) 
        internal 
        virtual 
        override(ERC721Upgradeable, ERC721URIStorageUpgradeable) 
    {
        super._burn(tokenId);
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