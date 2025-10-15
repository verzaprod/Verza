import { create } from 'zustand'

interface KYCState {
  // Escrow & Verifier Info
  escrowId?: string
  verifierId?: string
  verifierName?: string
  
  // Document Info
  selectedDocType?: 'passport' | 'driver-license' | 'id-card'
  documentImages?: string[]
  selfieImage?: string
  
  // Process State
  currentStep: 'selection' | 'docs' | 'selfie' | 'processing' | 'complete'
  isProcessing: boolean
  verificationStatus?: 'pending' | 'verified' | 'rejected'
}

interface KYCActions {
  setEscrowInfo: (escrowId: string, verifierId: string, verifierName?: string) => void
  setDocumentType: (type: 'passport' | 'driver-license' | 'id-card') => void
  setDocumentImages: (images: string[]) => void
  setSelfieImage: (image: string) => void
  setCurrentStep: (step: KYCState['currentStep']) => void
  setProcessing: (processing: boolean) => void
  setVerificationStatus: (status: 'pending' | 'verified' | 'rejected') => void
  resetKYC: () => void
}

type KYCStore = KYCState & KYCActions

const initialState: KYCState = {
  currentStep: 'selection',
  isProcessing: false,
  verificationStatus: "verified",
}

export const useKYCStore = create<KYCStore>((set) => ({
  ...initialState,
  
  setEscrowInfo: (escrowId, verifierId, verifierName) =>
    set({ escrowId, verifierId, verifierName }),
  
  setDocumentType: (selectedDocType) =>
    set({ selectedDocType }),
  
  setDocumentImages: (documentImages) =>
    set({ documentImages }),
  
  setSelfieImage: (selfieImage) =>
    set({ selfieImage }),
  
  setCurrentStep: (currentStep) =>
    set({ currentStep }),
  
  setProcessing: (isProcessing) =>
    set({ isProcessing }),

  setVerificationStatus: (verificationStatus) =>
    set({ verificationStatus }),
  
  resetKYC: () =>
    set(initialState),
}))