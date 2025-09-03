interface ApiResponse<T> {
  success: boolean;
  data?: T;
  error?: string;
}

class ApiClient {
  private baseUrl: string;

  constructor(baseUrl: string = process.env.EXPO_PUBLIC_API_URL || '') {
    this.baseUrl = baseUrl;
  }

  private async request<T>(
    endpoint: string,
    options: RequestInit = {}
  ): Promise<ApiResponse<T>> {
    try {
      const response = await fetch(`${this.baseUrl}${endpoint}`, {
        headers: {
          'Content-Type': 'application/json',
          ...options.headers,
        },
        ...options,
      });

      const data = await response.json();
      return { success: response.ok, data: response.ok ? data : undefined, error: response.ok ? undefined : data.message };
    } catch (error) {
      return { success: false, error: 'Network error occurred' };
    }
  }

  async register(emailOrPhone: string): Promise<ApiResponse<{ registered: boolean }>> {
    // TODO: Implement actual API call
    return this.request("/auth/register", {
      method: "POST",
      body: JSON.stringify({ emailOrPhone }),
    });
  }

  async sendVerificationCode(email: string): Promise<ApiResponse<{ sent: boolean }>> {
    return this.request('/auth/verify-email', {
      method: 'POST',
      body: JSON.stringify({ email }),
    });
  }

  async verifyCode(email: string, code: string): Promise<ApiResponse<{ verified: boolean }>> {
    return this.request('/auth/verify-code', {
      method: 'POST',
      body: JSON.stringify({ email, code }),
    });
  }

  async generatePassphrase(): Promise<ApiResponse<{ words: string[] }>> {
    return this.request('/wallet/generate-passphrase');
  }

  async verifyPassphrase(passphrase: string[]): Promise<ApiResponse<{ valid: boolean }>> {
    return this.request('/wallet/verify-passphrase', {
      method: 'POST',
      body: JSON.stringify({ passphrase }),
    });
  }
}

export const apiClient = new ApiClient();