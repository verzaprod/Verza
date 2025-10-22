// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

import "@openzeppelin/contracts-upgradeable/access/AccessControlUpgradeable.sol";
import "@openzeppelin/contracts-upgradeable/security/PausableUpgradeable.sol";
import "@openzeppelin/contracts-upgradeable/proxy/utils/Initializable.sol";
import "@openzeppelin/contracts-upgradeable/proxy/utils/UUPSUpgradeable.sol";
import "@openzeppelin/contracts-upgradeable/security/ReentrancyGuardUpgradeable.sol";
import "@openzeppelin/contracts/token/ERC20/IERC20.sol";

/**
 * @title VerifierMarketplace
 * @dev Smart contract for managing verifiers with reputation, staking, and dynamic pricing
 * Verifiers must stake HBAR or stablecoin to participate in the verification marketplace
 */
contract VerifierMarketplace is 
    Initializable, 
    AccessControlUpgradeable, 
    PausableUpgradeable, 
    UUPSUpgradeable,
    ReentrancyGuardUpgradeable 
{
    bytes32 public constant ADMIN_ROLE = keccak256("ADMIN_ROLE");
    bytes32 public constant UPGRADER_ROLE = keccak256("UPGRADER_ROLE");
    bytes32 public constant SLASHER_ROLE = keccak256("SLASHER_ROLE");
    
    // Minimum stake required to become a verifier (in wei)
    uint256 public minimumStake;
    
    // Base verification fee (in wei)
    uint256 public baseVerificationFee;
    
    // Reputation multiplier for pricing (basis points, 10000 = 100%)
    uint256 public reputationMultiplier;
    
    // Slashing percentage for fraud/inactivity (basis points)
    uint256 public slashingPercentage;
    
    // Inactivity threshold in seconds
    uint256 public inactivityThreshold;
    
    // Supported staking token (address(0) for native HBAR)
    IERC20 public stakingToken;
    
    struct Verifier {
        address verifierAddress;
        uint256 stakedAmount;
        uint256 reputationScore; // 0-10000 (basis points)
        uint256 totalVerifications;
        uint256 successfulVerifications;
        uint256 lastActivityTimestamp;
        bool isActive;
        string metadata; // IPFS hash or JSON metadata
        uint256 registrationTimestamp;
    }
    
    struct VerificationRequest {
        address requester;
        address assignedVerifier;
        uint256 fee;
        uint256 timestamp;
        bool completed;
        bool successful;
        string requestData; // Encrypted or hashed request data
    }
    
    // Mapping from verifier address to verifier info
    mapping(address => Verifier) public verifiers;
    
    // Array of all verifier addresses for enumeration
    address[] public verifierList;
    
    // Mapping from request ID to verification request
    mapping(bytes32 => VerificationRequest) public verificationRequests;
    
    // Mapping to track verifier earnings
    mapping(address => uint256) public verifierEarnings;
    
    // Events
    event VerifierRegistered(
        address indexed verifier,
        uint256 stakedAmount,
        string metadata
    );
    
    event VerifierStakeIncreased(
        address indexed verifier,
        uint256 additionalAmount,
        uint256 totalStake
    );
    
    event VerifierSlashed(
        address indexed verifier,
        uint256 slashedAmount,
        string reason
    );
    
    event VerificationCompleted(
        bytes32 indexed requestId,
        address indexed verifier,
        bool successful,
        uint256 fee
    );
    
    event ReputationUpdated(
        address indexed verifier,
        uint256 oldScore,
        uint256 newScore
    );
    
    event VerifierDeactivated(
        address indexed verifier,
        string reason
    );
    
    event EarningsWithdrawn(
        address indexed verifier,
        uint256 amount
    );
    
    /// @custom:oz-upgrades-unsafe-allow constructor
    constructor() {
        _disableInitializers();
    }
    
    function initialize(
        address admin,
        uint256 _minimumStake,
        uint256 _baseVerificationFee,
        address _stakingToken
    ) public initializer {
        __AccessControl_init();
        __Pausable_init();
        __UUPSUpgradeable_init();
        __ReentrancyGuard_init();
        
        _grantRole(DEFAULT_ADMIN_ROLE, admin);
        _grantRole(ADMIN_ROLE, admin);
        _grantRole(UPGRADER_ROLE, admin);
        _grantRole(SLASHER_ROLE, admin);
        
        minimumStake = _minimumStake;
        baseVerificationFee = _baseVerificationFee;
        reputationMultiplier = 5000; // 50% impact on pricing
        slashingPercentage = 1000; // 10% slashing
        inactivityThreshold = 30 days;
        
        if (_stakingToken != address(0)) {
            stakingToken = IERC20(_stakingToken);
        }
    }
    
    /**
     * @dev Register as a verifier with required stake
     * @param metadata IPFS hash or JSON metadata about the verifier
     */
    function registerVerifier(string calldata metadata) 
        external 
        payable 
        whenNotPaused 
        nonReentrant 
    {
        require(!verifiers[msg.sender].isActive, "Already registered as verifier");
        require(bytes(metadata).length > 0, "Metadata cannot be empty");
        
        uint256 stakeAmount;
        
        if (address(stakingToken) == address(0)) {
            // Native HBAR staking
            require(msg.value >= minimumStake, "Insufficient stake amount");
            stakeAmount = msg.value;
        } else {
            // ERC20 token staking
            require(msg.value == 0, "No native currency for token staking");
            stakeAmount = minimumStake;
            require(
                stakingToken.transferFrom(msg.sender, address(this), stakeAmount),
                "Stake transfer failed"
            );
        }
        
        verifiers[msg.sender] = Verifier({
            verifierAddress: msg.sender,
            stakedAmount: stakeAmount,
            reputationScore: 5000, // Start with neutral reputation (50%)
            totalVerifications: 0,
            successfulVerifications: 0,
            lastActivityTimestamp: block.timestamp,
            isActive: true,
            metadata: metadata,
            registrationTimestamp: block.timestamp
        });
        
        verifierList.push(msg.sender);
        
        emit VerifierRegistered(msg.sender, stakeAmount, metadata);
    }
    
    /**
     * @dev Increase stake amount
     */
    function increaseStake() 
        external 
        payable 
        whenNotPaused 
        nonReentrant 
    {
        require(verifiers[msg.sender].isActive, "Not a registered verifier");
        
        uint256 additionalAmount;
        
        if (address(stakingToken) == address(0)) {
            require(msg.value > 0, "Must send HBAR to increase stake");
            additionalAmount = msg.value;
        } else {
            require(msg.value == 0, "No native currency for token staking");
            additionalAmount = minimumStake; // Allow fixed increments
            require(
                stakingToken.transferFrom(msg.sender, address(this), additionalAmount),
                "Stake transfer failed"
            );
        }
        
        verifiers[msg.sender].stakedAmount += additionalAmount;
        
        emit VerifierStakeIncreased(
            msg.sender, 
            additionalAmount, 
            verifiers[msg.sender].stakedAmount
        );
    }
    
    /**
     * @dev Get comprehensive verifier information
     * @param verifier Address of the verifier
     * @return name Verifier name (extracted from metadata)
     * @return metadataURI URI pointing to verifier metadata JSON
     * @return baseFee Current verification fee for this verifier
     * @return isActive Whether the verifier is currently active
     * @return reputationScore Reputation score (0-10000 basis points)
     * @return totalVerifications Total number of verifications completed
     * @return successfulVerifications Number of successful verifications
     * @return stakedAmount Amount currently staked by verifier
     */
    function getVerifierDetails(address verifier) 
        external 
        view 
        returns (
            string memory name,
            string memory metadataURI,
            uint256 baseFee,
            bool isActive,
            uint256 reputationScore,
            uint256 totalVerifications,
            uint256 successfulVerifications,
            uint256 stakedAmount
        ) 
    {
        require(verifiers[verifier].verifierAddress != address(0), "Verifier not found");
        
        Verifier memory v = verifiers[verifier];
        
        // Extract name from metadata if it's JSON format, otherwise use address as fallback
        name = _extractNameFromMetadata(v.metadata, verifier);
        metadataURI = v.metadata;
        baseFee = calculateVerificationFee(verifier);
        isActive = v.isActive;
        reputationScore = v.reputationScore;
        totalVerifications = v.totalVerifications;
        successfulVerifications = v.successfulVerifications;
        stakedAmount = v.stakedAmount;
    }
    
    /**
     * @dev Get all active verifiers with their basic info
     * @return verifierAddresses Array of active verifier addresses
     * @return names Array of verifier names
     * @return fees Array of current verification fees
     * @return reputationScores Array of reputation scores
     */
    function getActiveVerifiersWithDetails() 
        external 
        view 
        returns (
            address[] memory verifierAddresses,
            string[] memory names,
            uint256[] memory fees,
            uint256[] memory reputationScores
        ) 
    {
        // Count active verifiers
        uint256 activeCount = 0;
        for (uint256 i = 0; i < verifierList.length; i++) {
            if (verifiers[verifierList[i]].isActive) {
                activeCount++;
            }
        }
        
        // Initialize arrays
        verifierAddresses = new address[](activeCount);
        names = new string[](activeCount);
        fees = new uint256[](activeCount);
        reputationScores = new uint256[](activeCount);
        
        // Populate arrays
        uint256 index = 0;
        for (uint256 i = 0; i < verifierList.length; i++) {
            address verifierAddr = verifierList[i];
            if (verifiers[verifierAddr].isActive) {
                verifierAddresses[index] = verifierAddr;
                names[index] = _extractNameFromMetadata(verifiers[verifierAddr].metadata, verifierAddr);
                fees[index] = calculateVerificationFee(verifierAddr);
                reputationScores[index] = verifiers[verifierAddr].reputationScore;
                index++;
            }
        }
    }
    
    /**
     * @dev Extract name from metadata string (simple JSON parsing or fallback to address)
     * @param metadata The metadata string (JSON or IPFS hash)
     * @param verifierAddr Fallback verifier address
     * @return name Extracted name or formatted address
     */
    function _extractNameFromMetadata(string memory metadata, address verifierAddr) 
        internal 
        pure 
        returns (string memory name) 
    {
        // Simple check if metadata looks like JSON (starts with '{')
        bytes memory metadataBytes = bytes(metadata);
        if (metadataBytes.length > 0 && metadataBytes[0] == 0x7B) { // '{'
            // For now, return a placeholder - in practice, you'd parse JSON
            // This is a simplified implementation
            return "Verifier"; // Could be enhanced with actual JSON parsing
        } else if (metadataBytes.length > 0) {
            // Assume it's an IPFS hash or URI, return generic name
            return "Verified Professional";
        } else {
            // Fallback to formatted address
            return string(abi.encodePacked("Verifier ", _addressToString(verifierAddr)));
        }
    }
    
    /**
     * @dev Convert address to string
     * @param addr Address to convert
     * @return String representation of address
     */
    function _addressToString(address addr) internal pure returns (string memory) {
        bytes32 value = bytes32(uint256(uint160(addr)));
        bytes memory alphabet = "0123456789abcdef";
        bytes memory str = new bytes(42);
        str[0] = '0';
        str[1] = 'x';
        for (uint256 i = 0; i < 20; i++) {
            str[2 + i * 2] = alphabet[uint8(value[i + 12] >> 4)];
            str[3 + i * 2] = alphabet[uint8(value[i + 12] & 0x0f)];
        }
        return string(str);
    }

    /**
     * @dev Calculate verification fee based on verifier reputation
     * @param verifier Address of the verifier
     * @return fee The calculated fee
     */
    function calculateVerificationFee(address verifier) 
        public 
        view 
        returns (uint256 fee) 
    {
        require(verifiers[verifier].isActive, "Verifier not active");
        
        uint256 reputation = verifiers[verifier].reputationScore;
        
        // Higher reputation = higher fee (premium pricing)
        // Fee = baseVerificationFee * (1 + (reputation - 5000) * reputationMultiplier / 10000 / 10000)
        int256 reputationAdjustment = (int256(reputation) - 5000) * int256(reputationMultiplier) / 10000;
        
        if (reputationAdjustment >= 0) {
            fee = baseVerificationFee + (baseVerificationFee * uint256(reputationAdjustment) / 10000);
        } else {
            uint256 reduction = baseVerificationFee * uint256(-reputationAdjustment) / 10000;
            fee = reduction >= baseVerificationFee ? baseVerificationFee / 2 : baseVerificationFee - reduction;
        }
    }
    
    /**
     * @dev Complete a verification and update reputation
     * @param requestId Unique identifier for the verification request
     * @param verifier Address of the verifier
     * @param successful Whether the verification was successful
     */
    function completeVerification(
        bytes32 requestId,
        address verifier,
        bool successful
    ) 
        external 
        onlyRole(ADMIN_ROLE) 
        whenNotPaused 
    {
        require(verifiers[verifier].isActive, "Verifier not active");
        require(!verificationRequests[requestId].completed, "Verification already completed");
        
        verificationRequests[requestId].completed = true;
        verificationRequests[requestId].successful = successful;
        
        // Update verifier stats
        verifiers[verifier].totalVerifications++;
        verifiers[verifier].lastActivityTimestamp = block.timestamp;
        
        if (successful) {
            verifiers[verifier].successfulVerifications++;
        }
        
        // Update reputation score
        _updateReputation(verifier);
        
        // Pay verifier
        uint256 fee = verificationRequests[requestId].fee;
        verifierEarnings[verifier] += fee;
        
        emit VerificationCompleted(requestId, verifier, successful, fee);
    }
    
    /**
     * @dev Slash a verifier for fraud or inactivity
     * @param verifier Address of the verifier to slash
     * @param reason Reason for slashing
     */
    function slashVerifier(address verifier, string calldata reason) 
        external 
        onlyRole(SLASHER_ROLE) 
        whenNotPaused 
        nonReentrant 
    {
        require(verifiers[verifier].isActive, "Verifier not active");
        
        uint256 slashAmount = verifiers[verifier].stakedAmount * slashingPercentage / 10000;
        verifiers[verifier].stakedAmount -= slashAmount;
        
        // If stake falls below minimum, deactivate verifier
        if (verifiers[verifier].stakedAmount < minimumStake) {
            verifiers[verifier].isActive = false;
            emit VerifierDeactivated(verifier, "Stake below minimum after slashing");
        }
        
        // Reduce reputation significantly
        uint256 oldScore = verifiers[verifier].reputationScore;
        verifiers[verifier].reputationScore = verifiers[verifier].reputationScore * 7 / 10; // 30% reduction
        
        emit VerifierSlashed(verifier, slashAmount, reason);
        emit ReputationUpdated(verifier, oldScore, verifiers[verifier].reputationScore);
    }
    
    /**
     * @dev Withdraw earnings
     */
    function withdrawEarnings() 
        external 
        whenNotPaused 
        nonReentrant 
    {
        uint256 earnings = verifierEarnings[msg.sender];
        require(earnings > 0, "No earnings to withdraw");
        
        verifierEarnings[msg.sender] = 0;
        
        if (address(stakingToken) == address(0)) {
            payable(msg.sender).transfer(earnings);
        } else {
            require(stakingToken.transfer(msg.sender, earnings), "Transfer failed");
        }
        
        emit EarningsWithdrawn(msg.sender, earnings);
    }
    
    /**
     * @dev Get verifier information
     * @param verifier Address of the verifier
     * @return Verifier struct
     */
    function getVerifier(address verifier) 
        external 
        view 
        returns (Verifier memory) 
    {
        return verifiers[verifier];
    }
    
    /**
     * @dev Get all active verifiers
     * @return Array of active verifier addresses
     */
    function getActiveVerifiers() 
        external 
        view 
        returns (address[] memory) 
    {
        uint256 activeCount = 0;
        
        // Count active verifiers
        for (uint256 i = 0; i < verifierList.length; i++) {
            if (verifiers[verifierList[i]].isActive) {
                activeCount++;
            }
        }
        
        // Create array of active verifiers
        address[] memory activeVerifiers = new address[](activeCount);
        uint256 index = 0;
        
        for (uint256 i = 0; i < verifierList.length; i++) {
            if (verifiers[verifierList[i]].isActive) {
                activeVerifiers[index] = verifierList[i];
                index++;
            }
        }
        
        return activeVerifiers;
    }
    
    /**
     * @dev Check for inactive verifiers and deactivate them
     */
    function checkInactiveVerifiers() 
        external 
        onlyRole(ADMIN_ROLE) 
    {
        for (uint256 i = 0; i < verifierList.length; i++) {
            address verifier = verifierList[i];
            if (verifiers[verifier].isActive && 
                block.timestamp - verifiers[verifier].lastActivityTimestamp > inactivityThreshold) {
                
                verifiers[verifier].isActive = false;
                emit VerifierDeactivated(verifier, "Inactivity threshold exceeded");
            }
        }
    }
    
    /**
     * @dev Update reputation score based on success rate
     * @param verifier Address of the verifier
     */
    function _updateReputation(address verifier) internal {
        uint256 total = verifiers[verifier].totalVerifications;
        uint256 successful = verifiers[verifier].successfulVerifications;
        
        if (total == 0) return;
        
        uint256 oldScore = verifiers[verifier].reputationScore;
        
        // Calculate success rate (0-10000 basis points)
        uint256 successRate = (successful * 10000) / total;
        
        // Weighted average with existing reputation (70% old, 30% new)
        uint256 newScore = (oldScore * 7 + successRate * 3) / 10;
        
        // Ensure score stays within bounds
        if (newScore > 10000) newScore = 10000;
        
        verifiers[verifier].reputationScore = newScore;
        
        emit ReputationUpdated(verifier, oldScore, newScore);
    }
    
    /**
     * @dev Admin functions for updating parameters
     */
    function updateMinimumStake(uint256 _minimumStake) 
        external 
        onlyRole(ADMIN_ROLE) 
    {
        minimumStake = _minimumStake;
    }
    
    function updateBaseVerificationFee(uint256 _baseVerificationFee) 
        external 
        onlyRole(ADMIN_ROLE) 
    {
        baseVerificationFee = _baseVerificationFee;
    }
    
    function updateReputationMultiplier(uint256 _reputationMultiplier) 
        external 
        onlyRole(ADMIN_ROLE) 
    {
        reputationMultiplier = _reputationMultiplier;
    }
    
    function updateSlashingPercentage(uint256 _slashingPercentage) 
        external 
        onlyRole(ADMIN_ROLE) 
    {
        require(_slashingPercentage <= 5000, "Slashing percentage too high"); // Max 50%
        slashingPercentage = _slashingPercentage;
    }
    
    function updateInactivityThreshold(uint256 _inactivityThreshold) 
        external 
        onlyRole(ADMIN_ROLE) 
    {
        inactivityThreshold = _inactivityThreshold;
    }
    
    /**
     * @dev Pause/unpause contract
     */
    function pause() external onlyRole(ADMIN_ROLE) {
        _pause();
    }
    
    function unpause() external onlyRole(ADMIN_ROLE) {
        _unpause();
    }
    
    /**
     * @dev Authorize upgrade
     */
    function _authorizeUpgrade(address newImplementation) 
        internal 
        onlyRole(UPGRADER_ROLE) 
        override 
    {}
    
    /**
     * @dev Emergency withdrawal function (only admin)
     */
    function emergencyWithdraw() 
        external 
        onlyRole(DEFAULT_ADMIN_ROLE) 
        nonReentrant 
    {
        if (address(stakingToken) == address(0)) {
            payable(msg.sender).transfer(address(this).balance);
        } else {
            uint256 balance = stakingToken.balanceOf(address(this));
            require(stakingToken.transfer(msg.sender, balance), "Transfer failed");
        }
    }
    
    /**
     * @dev Receive function for native currency deposits
     */
    receive() external payable {
        // Allow contract to receive native currency
    }
}