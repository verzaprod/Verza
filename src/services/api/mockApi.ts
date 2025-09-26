import { MOCK_DATA } from "./mockData";
import { API_CONFIG } from "./config";

const delay = (ms: number) => new Promise((resolve) => setTimeout(resolve, ms));

export class MockAPI {
  private static verificationProgress = 0;

  static async initiateEscrow(data: any) {
    await delay(API_CONFIG.MOCK_DELAY);

    // Simulate success/failure
    if (Math.random() > 0.1) {
      // 90% success rate
      return {
        ok: true,
        json: async () => MOCK_DATA.escrow.initiate,
      };
    } else {
      return {
        ok: false,
        json: async () => ({ error: "Payment processing failed" }),
      };
    }
  }

  static async getVerificationStatus(escrowId: string) {
    await delay(1000); // Shorter delay for polling

    // Simulate progress
    this.verificationProgress += 1;

    let status = MOCK_DATA.verification.status;

    if (this.verificationProgress > 3) {
      status = {
        ...status,
        status: "completed",
        steps: status.steps.map((step, index) => ({
          ...step,
          status: index <= 2 ? "completed" : "pending",
          timestamp:
            index <= 2 ? `2024-01-15 ${10 + index}:${30 + index * 5} AM` : null,
        })),
      };
    } else if (this.verificationProgress > 1) {
      status = {
        ...status,
        steps: status.steps.map((step, index) => ({
          ...step,
          status:
            index === 0 ? "completed" : index === 1 ? "active" : "pending",
          timestamp: index === 0 ? "2024-01-15 10:30 AM" : null,
        })),
      };
    }

    return {
      ok: true,
      json: async () => status,
    };
  }

  static async getVerificationResults(escrowId: string) {
    await delay(API_CONFIG.MOCK_DELAY);

    // Simulate 80% verification success
    if (Math.random() > 0.2) {
      return {
        ok: true,
        json: async () => MOCK_DATA.verification.results,
      };
    } else {
      return {
        ok: true,
        json: async () => ({
          status: "rejected",
          rejectionReason:
            "Document quality insufficient. Please retake photos with better lighting and ensure all text is clearly visible.",
        }),
      };
    }
  }
}
