export const MOCK_DATA = {
  escrow: {
    initiate: {
      escrow_id: 'escrow_123456789',
      status: 'pending',
      amount: 25.00,
      currency: 'HBAR',
      created_at: new Date().toISOString(),
    }
  },
  
  verification: {
    status: {
      escrowId: 'escrow_123456789',
      status: 'in_progress', // 'submitted' | 'in_progress' | 'completed' | 'failed'
      steps: [
        {
          id: '1',
          label: 'Documents Uploaded',
          status: 'completed',
          timestamp: '2024-01-15 10:30 AM'
        },
        {
          id: '2', 
          label: 'Identity Verification',
          status: 'active',
          timestamp: null
        },
        {
          id: '3',
          label: 'Final Review',
          status: 'pending',
          timestamp: null
        }
      ],
      estimatedCompletion: '2-3 minutes'
    },
    
    results: {
      status: 'verified',
      credentialId: 'cred_987654321',
      vcDetails: {
        id: 'vc_123456789abcdef',
        issuer: 'TechCorp Solutions',
        issuedDate: '2024-01-15',
        expiryDate: '2025-01-15',
        type: 'Identity Verification',
        proofHash: '0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef'
      }
    }
  }
}