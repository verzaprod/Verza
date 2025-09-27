// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

import "@openzeppelin/contracts-upgradeable/access/AccessControlUpgradeable.sol";
import "@openzeppelin/contracts-upgradeable/security/PausableUpgradeable.sol";
import "@openzeppelin/contracts-upgradeable/proxy/utils/Initializable.sol";
import "@openzeppelin/contracts-upgradeable/proxy/utils/UUPSUpgradeable.sol";
import "@openzeppelin/contracts-upgradeable/security/ReentrancyGuardUpgradeable.sol";
import "@openzeppelin/contracts/token/ERC20/IERC20.sol";

interface IVerifierMarketplace {
    function getVerifier(address verifier) external view returns (
        address verifierAddress,
        uint256 stakedAmount,
        uint256 reputationScore,
        uint256 totalVerifications,
        uint256 successfulVerifications,
        uint256 lastActivityTimestamp,
        bool isActive,
        string memory metadata,
        uint256 registrationTimestamp
    );
    
    function calculateVerificationFee(address verifier) external view returns (uint256);
}

interface IFraudDetection {
    function checkFraudScore(bytes32 requestId, bytes calldata verificationData) external returns (uint256 riskScore, bool isFraud);
}

/**
 * @title EscrowContract
 * @dev Smart contract for managing escrow funds during verification processes
 * Handles fund locking, release, refunds, and dispute resolution with AI fraud detection
 */
contract EscrowContract is 
    Initializable, 
    AccessControlUpgradeable, 
    PausableUpgradeable, 
    UUPSUpgradeable,
    ReentrancyGuardUpgradeable 
{
    bytes32 public constant ADMIN_ROLE = keccak256("ADMIN_ROLE");
    bytes32 public constant UPGRADER_ROLE = keccak256("UPGRADER_ROLE");
    bytes32 public constant ORACLE_ROLE = keccak256("ORACLE_ROLE");
    bytes32 public constant DISPUTE_RESOLVER_ROLE = keccak256("DISPUTE_RESOLVER_ROLE");
    
    enum EscrowStatus {
        Created,
        FundsLocked,
        VerificationSubmitted,
        FraudCheckPending,
        DisputeRaised,
        Completed,
        Refunded,
        Cancelled
    }
    
    enum DisputeStatus {
        None,
        Raised,
        UnderReview,
        Resolved
    }
    
    struct EscrowRequest {
        bytes32 requestId;
        address user;
        address verifier;
        uint256 amount;
        uint256 createdAt;
        uint256 expiresAt;
        EscrowStatus status;
        string verificationData; // Encrypted or hashed verification data
        uint256 fraudScore;
        bool fraudDetected;
        DisputeStatus disputeStatus;
        string disputeReason;
        address disputeResolver;
        uint256 resolvedAt;
    }
    
    struct DisputeInfo {
        bytes32 requestId;
        address initiator;
        string reason;
        uint256 createdAt;
        uint256 resolvedAt;
        address resolver;
        bool userFavored; // True if dispute resolved in favor of user
        string resolution;
    }
    
    // Contract references
    IVerifierMarketplace public verifierMarketplace;
    IFraudDetection public fraudDetection;
    IERC20 public paymentToken; // address(0) for native HBAR
    
    // Configuration
    uint256 public escrowTimeout; // Default timeout for escrow requests
    uint256 public disputeTimeout; // Timeout for dispute resolution
    uint256 public fraudThreshold; // Risk score threshold for fraud detection (0-100)
    uint256 public platformFeePercentage; // Platform fee in basis points
    address public feeRecipient; // Address to receive platform fees
    
    // Storage
    mapping(bytes32 => EscrowRequest) public escrowRequests;
    mapping(bytes32 => DisputeInfo) public disputes;
    mapping(address => uint256) public userRefunds; // Pending refunds for users
    mapping(address => uint256) public verifierEarnings; // Pending earnings for verifiers
    
    // Statistics
    uint256 public totalEscrowsCreated;
    uint256 public totalAmountEscrowed;
    uint256 public totalDisputesRaised;
    uint256 public totalFraudDetected;
    
    // Events
    event EscrowCreated(
        bytes32 indexed requestId,
        address indexed user,
        address indexed verifier,
        uint256 amount
    );
    
    event FundsLocked(
        bytes32 indexed requestId,
        uint256 amount,
        uint256 expiresAt
    );
    
    event VerificationSubmitted(
        bytes32 indexed requestId,
        address indexed verifier,
        string verificationData
    );
    
    event FraudCheckCompleted(
        bytes32 indexed requestId,
        uint256 fraudScore,
        bool fraudDetected
    );
    
    event DisputeRaised(
        bytes32 indexed requestId,
        address indexed initiator,
        string reason
    );
    
    event DisputeResolved(
        bytes32 indexed requestId,
        address indexed resolver,
        bool userFavored,
        string resolution
    );
    
    event FundsReleased(
        bytes32 indexed requestId,
        address indexed verifier,
        uint256 amount,
        uint256 platformFee
    );
    
    event RefundIssued(
        bytes32 indexed requestId,
        address indexed user,
        uint256 amount
    );
    
    event EscrowCancelled(
        bytes32 indexed requestId,
        string reason
    );
    
    /// @custom:oz-upgrades-unsafe-allow constructor
    constructor() {
        _disableInitializers();
    }
    
    function initialize(
        address admin,
        address _verifierMarketplace,
        address _fraudDetection,
        address _paymentToken,
        address _feeRecipient
    ) public initializer {
        __AccessControl_init();
        __Pausable_init();
        __UUPSUpgradeable_init();
        __ReentrancyGuard_init();
        
        _grantRole(DEFAULT_ADMIN_ROLE, admin);
        _grantRole(ADMIN_ROLE, admin);
        _grantRole(UPGRADER_ROLE, admin);
        _grantRole(ORACLE_ROLE, admin);
        _grantRole(DISPUTE_RESOLVER_ROLE, admin);
        
        verifierMarketplace = IVerifierMarketplace(_verifierMarketplace);
        fraudDetection = IFraudDetection(_fraudDetection);
        paymentToken = IERC20(_paymentToken);
        feeRecipient = _feeRecipient;
        
        escrowTimeout = 7 days;
        disputeTimeout = 3 days;
        fraudThreshold = 70; // 70% risk score threshold
        platformFeePercentage = 250; // 2.5% platform fee
    }
    
    /**
     * @dev Create a new escrow request
     * @param requestId Unique identifier for the verification request
     * @param verifier Address of the assigned verifier
     */
    function createEscrow(
        bytes32 requestId,
        address verifier
    ) 
        external 
        payable 
        whenNotPaused 
        nonReentrant 
    {
        require(requestId != bytes32(0), "Invalid request ID");
        require(verifier != address(0), "Invalid verifier address");
        require(escrowRequests[requestId].requestId == bytes32(0), "Escrow already exists");
        
        // Verify verifier is active in marketplace
        (, , , , , , bool isActive, ,) = verifierMarketplace.getVerifier(verifier);
        require(isActive, "Verifier not active");
        
        // Calculate verification fee
        uint256 verificationFee = verifierMarketplace.calculateVerificationFee(verifier);
        
        uint256 amount;
        if (address(paymentToken) == address(0)) {
            // Native HBAR payment
            require(msg.value >= verificationFee, "Insufficient payment");
            amount = msg.value;
        } else {
            // ERC20 token payment
            require(msg.value == 0, "No native currency for token payment");
            amount = verificationFee;
            require(
                paymentToken.transferFrom(msg.sender, address(this), amount),
                "Payment transfer failed"
            );
        }
        
        escrowRequests[requestId] = EscrowRequest({
            requestId: requestId,
            user: msg.sender,
            verifier: verifier,
            amount: amount,
            createdAt: block.timestamp,
            expiresAt: block.timestamp + escrowTimeout,
            status: EscrowStatus.Created,
            verificationData: "",
            fraudScore: 0,
            fraudDetected: false,
            disputeStatus: DisputeStatus.None,
            disputeReason: "",
            disputeResolver: address(0),
            resolvedAt: 0
        });
        
        totalEscrowsCreated++;
        totalAmountEscrowed += amount;
        
        emit EscrowCreated(requestId, msg.sender, verifier, amount);
    }
    
    /**
     * @dev Lock funds for verification process
     * @param requestId The escrow request ID
     */
    function lockFunds(bytes32 requestId) 
        external 
        onlyRole(ORACLE_ROLE) 
        whenNotPaused 
    {
        EscrowRequest storage request = escrowRequests[requestId];
        require(request.status == EscrowStatus.Created, "Invalid status for locking funds");
        require(block.timestamp <= request.expiresAt, "Escrow expired");
        
        request.status = EscrowStatus.FundsLocked;
        
        emit FundsLocked(requestId, request.amount, request.expiresAt);
    }
    
    /**
     * @dev Submit verification result
     * @param requestId The escrow request ID
     * @param verificationData Encrypted or hashed verification data
     */
    function submitVerification(
        bytes32 requestId,
        string calldata verificationData
    ) 
        external 
        whenNotPaused 
    {
        EscrowRequest storage request = escrowRequests[requestId];
        require(msg.sender == request.verifier, "Only assigned verifier can submit");
        require(request.status == EscrowStatus.FundsLocked, "Invalid status for submission");
        require(block.timestamp <= request.expiresAt, "Escrow expired");
        require(bytes(verificationData).length > 0, "Verification data cannot be empty");
        
        request.verificationData = verificationData;
        request.status = EscrowStatus.VerificationSubmitted;
        
        emit VerificationSubmitted(requestId, msg.sender, verificationData);
        
        // Trigger fraud detection
        _triggerFraudCheck(requestId);
    }
    
    /**
     * @dev Complete fraud check and update escrow status
     * @param requestId The escrow request ID
     * @param riskScore Risk score from fraud detection (0-100)
     * @param isFraud Whether fraud was detected
     */
    function completeFraudCheck(
        bytes32 requestId,
        uint256 riskScore,
        bool isFraud
    ) 
        external 
        onlyRole(ORACLE_ROLE) 
        whenNotPaused 
    {
        EscrowRequest storage request = escrowRequests[requestId];
        require(request.status == EscrowStatus.FraudCheckPending, "Invalid status for fraud check");
        
        request.fraudScore = riskScore;
        request.fraudDetected = isFraud;
        
        if (isFraud) {
            totalFraudDetected++;
            // Automatically refund user if fraud detected
            _issueRefund(requestId, "Fraud detected by AI system");
        } else {
            // Release funds to verifier
            _releaseFunds(requestId);
        }
        
        emit FraudCheckCompleted(requestId, riskScore, isFraud);
    }
    
    /**
     * @dev Raise a dispute
     * @param requestId The escrow request ID
     * @param reason Reason for the dispute
     */
    function raiseDispute(
        bytes32 requestId,
        string calldata reason
    ) 
        external 
        whenNotPaused 
    {
        EscrowRequest storage request = escrowRequests[requestId];
        require(
            msg.sender == request.user || msg.sender == request.verifier,
            "Only user or verifier can raise dispute"
        );
        require(
            request.status == EscrowStatus.VerificationSubmitted || 
            request.status == EscrowStatus.Completed,
            "Invalid status for dispute"
        );
        require(request.disputeStatus == DisputeStatus.None, "Dispute already raised");
        require(bytes(reason).length > 0, "Dispute reason cannot be empty");
        
        request.disputeStatus = DisputeStatus.Raised;
        request.disputeReason = reason;
        request.status = EscrowStatus.DisputeRaised;
        
        disputes[requestId] = DisputeInfo({
            requestId: requestId,
            initiator: msg.sender,
            reason: reason,
            createdAt: block.timestamp,
            resolvedAt: 0,
            resolver: address(0),
            userFavored: false,
            resolution: ""
        });
        
        totalDisputesRaised++;
        
        emit DisputeRaised(requestId, msg.sender, reason);
    }
    
    /**
     * @dev Resolve a dispute
     * @param requestId The escrow request ID
     * @param userFavored Whether the dispute is resolved in favor of the user
     * @param resolution Resolution details
     */
    function resolveDispute(
        bytes32 requestId,
        bool userFavored,
        string calldata resolution
    ) 
        external 
        onlyRole(DISPUTE_RESOLVER_ROLE) 
        whenNotPaused 
    {
        EscrowRequest storage request = escrowRequests[requestId];
        require(request.status == EscrowStatus.DisputeRaised, "No active dispute");
        require(request.disputeStatus == DisputeStatus.Raised, "Dispute not in correct state");
        
        DisputeInfo storage dispute = disputes[requestId];
        dispute.resolver = msg.sender;
        dispute.resolvedAt = block.timestamp;
        dispute.userFavored = userFavored;
        dispute.resolution = resolution;
        
        request.disputeStatus = DisputeStatus.Resolved;
        request.disputeResolver = msg.sender;
        request.resolvedAt = block.timestamp;
        
        if (userFavored) {
            _issueRefund(requestId, "Dispute resolved in favor of user");
        } else {
            _releaseFunds(requestId);
        }
        
        emit DisputeResolved(requestId, msg.sender, userFavored, resolution);
    }
    
    /**
     * @dev Cancel an expired escrow
     * @param requestId The escrow request ID
     */
    function cancelExpiredEscrow(bytes32 requestId) 
        external 
        whenNotPaused 
    {
        EscrowRequest storage request = escrowRequests[requestId];
        require(block.timestamp > request.expiresAt, "Escrow not expired");
        require(
            request.status == EscrowStatus.Created || 
            request.status == EscrowStatus.FundsLocked,
            "Cannot cancel escrow in current status"
        );
        
        _issueRefund(requestId, "Escrow expired");
        
        emit EscrowCancelled(requestId, "Escrow expired");
    }
    
    /**
     * @dev Withdraw pending refunds
     */
    function withdrawRefund() 
        external 
        whenNotPaused 
        nonReentrant 
    {
        uint256 refundAmount = userRefunds[msg.sender];
        require(refundAmount > 0, "No refund available");
        
        userRefunds[msg.sender] = 0;
        
        _transferFunds(msg.sender, refundAmount);
    }
    
    /**
     * @dev Withdraw verifier earnings
     */
    function withdrawEarnings() 
        external 
        whenNotPaused 
        nonReentrant 
    {
        uint256 earnings = verifierEarnings[msg.sender];
        require(earnings > 0, "No earnings available");
        
        verifierEarnings[msg.sender] = 0;
        
        _transferFunds(msg.sender, earnings);
    }
    
    /**
     * @dev Get escrow request details
     * @param requestId The escrow request ID
     * @return EscrowRequest struct
     */
    function getEscrowRequest(bytes32 requestId) 
        external 
        view 
        returns (EscrowRequest memory) 
    {
        return escrowRequests[requestId];
    }
    
    /**
     * @dev Get dispute information
     * @param requestId The escrow request ID
     * @return DisputeInfo struct
     */
    function getDispute(bytes32 requestId) 
        external 
        view 
        returns (DisputeInfo memory) 
    {
        return disputes[requestId];
    }
    
    /**
     * @dev Internal function to trigger fraud check
     * @param requestId The escrow request ID
     */
    function _triggerFraudCheck(bytes32 requestId) internal {
        EscrowRequest storage request = escrowRequests[requestId];
        request.status = EscrowStatus.FraudCheckPending;
        
        // In a real implementation, this would trigger an external fraud detection service
        // For now, we'll emit an event that can be picked up by off-chain services
    }
    
    /**
     * @dev Internal function to release funds to verifier
     * @param requestId The escrow request ID
     */
    function _releaseFunds(bytes32 requestId) internal {
        EscrowRequest storage request = escrowRequests[requestId];
        
        uint256 platformFee = (request.amount * platformFeePercentage) / 10000;
        uint256 verifierAmount = request.amount - platformFee;
        
        verifierEarnings[request.verifier] += verifierAmount;
        
        if (platformFee > 0 && feeRecipient != address(0)) {
            _transferFunds(feeRecipient, platformFee);
        }
        
        request.status = EscrowStatus.Completed;
        
        emit FundsReleased(requestId, request.verifier, verifierAmount, platformFee);
    }
    
    /**
     * @dev Internal function to issue refund to user
     * @param requestId The escrow request ID
     * @param reason Reason for refund
     */
    function _issueRefund(bytes32 requestId, string memory reason) internal {
        EscrowRequest storage request = escrowRequests[requestId];
        
        userRefunds[request.user] += request.amount;
        request.status = EscrowStatus.Refunded;
        
        emit RefundIssued(requestId, request.user, request.amount);
    }
    
    /**
     * @dev Internal function to transfer funds
     * @param to Recipient address
     * @param amount Amount to transfer
     */
    function _transferFunds(address to, uint256 amount) internal {
        if (address(paymentToken) == address(0)) {
            payable(to).transfer(amount);
        } else {
            require(paymentToken.transfer(to, amount), "Transfer failed");
        }
    }
    
    /**
     * @dev Admin functions
     */
    function updateEscrowTimeout(uint256 _escrowTimeout) 
        external 
        onlyRole(ADMIN_ROLE) 
    {
        escrowTimeout = _escrowTimeout;
    }
    
    function updateDisputeTimeout(uint256 _disputeTimeout) 
        external 
        onlyRole(ADMIN_ROLE) 
    {
        disputeTimeout = _disputeTimeout;
    }
    
    function updateFraudThreshold(uint256 _fraudThreshold) 
        external 
        onlyRole(ADMIN_ROLE) 
    {
        require(_fraudThreshold <= 100, "Invalid fraud threshold");
        fraudThreshold = _fraudThreshold;
    }
    
    function updatePlatformFee(uint256 _platformFeePercentage) 
        external 
        onlyRole(ADMIN_ROLE) 
    {
        require(_platformFeePercentage <= 1000, "Platform fee too high"); // Max 10%
        platformFeePercentage = _platformFeePercentage;
    }
    
    function updateFeeRecipient(address _feeRecipient) 
        external 
        onlyRole(ADMIN_ROLE) 
    {
        feeRecipient = _feeRecipient;
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
        if (address(paymentToken) == address(0)) {
            payable(msg.sender).transfer(address(this).balance);
        } else {
            uint256 balance = paymentToken.balanceOf(address(this));
            require(paymentToken.transfer(msg.sender, balance), "Transfer failed");
        }
    }
    
    /**
     * @dev Receive function for native currency deposits
     */
    receive() external payable {
        // Allow contract to receive native currency
    }
}