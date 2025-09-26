import { API_CONFIG } from './config'
import { MockAPI } from './mockApi'

class APIService {
  private async makeRequest(endpoint: string, options: RequestInit = {}) {
    if (API_CONFIG.USE_MOCK) {
      // Route to mock API based on endpoint
      switch (endpoint) {
        case '/escrow/initiate':
          return MockAPI.initiateEscrow(JSON.parse(options.body as string))
        default:
          if (endpoint.includes('/escrow/status/')) {
            const escrowId = endpoint.split('/').pop()
            return MockAPI.getVerificationStatus(escrowId!)
          }
          if (endpoint.includes('/verification/results/')) {
            const escrowId = endpoint.split('/').pop()
            return MockAPI.getVerificationResults(escrowId!)
          }
          throw new Error(`Mock endpoint not implemented: ${endpoint}`)
      }
    } else {
      // Real API call
      const response = await fetch(`${API_CONFIG.BASE_URL}${endpoint}`, {
        ...options,
        headers: {
          'Content-Type': 'application/json',
          ...options.headers,
        },
      })
      return response
    }
  }

  async initiateEscrow(data: {
    verifier_id: string
    amount: number
    currency: string
    auto_release_hours: number
  }) {
    return this.makeRequest('/escrow/initiate', {
      method: 'POST',
      body: JSON.stringify(data),
      headers: {
        'Authorization': 'Bearer sess-12453'
      }
    })
  }

  async getVerificationStatus(escrowId: string) {
    return this.makeRequest(`/escrow/status/${escrowId}`)
  }

  async getVerificationResults(escrowId: string) {
    return this.makeRequest(`/verification/results/${escrowId}`)
  }
}

export const apiService = new APIService()