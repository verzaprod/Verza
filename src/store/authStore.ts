import { create } from 'zustand';
import { persist, createJSONStorage } from 'zustand/middleware';
import * as SecureStore from 'expo-secure-store';

interface AuthState {
  isAuthenticated: boolean;
  userEmail: string | null;
  onboardingComplete: boolean;
  pinCreated: boolean;
  passphraseBackedUp: boolean;
  setEmail: (email: string) => void;
  setOnboardingComplete: (complete: boolean) => void;
  setPinCreated: (created: boolean) => void;
  setPassphraseBackedUp: (backed: boolean) => void;
  reset: () => void;
}

const secureStorage = {
  getItem: async (name: string): Promise<string | null> => {
    return await SecureStore.getItemAsync(name);
  },
  setItem: async (name: string, value: string): Promise<void> => {
    await SecureStore.setItemAsync(name, value);
  },
  removeItem: async (name: string): Promise<void> => {
    await SecureStore.deleteItemAsync(name);
  },
};

export const useAuthStore = create<AuthState>()(
  persist(
    (set) => ({
      isAuthenticated: false,
      userEmail: null,
      onboardingComplete: false,
      pinCreated: false,
      passphraseBackedUp: false,
      setEmail: (email) => set({ userEmail: email }),
      setOnboardingComplete: (complete) => set({ onboardingComplete: complete }),
      setPinCreated: (created) => set({ pinCreated: created }),
      setPassphraseBackedUp: (backed) => set({ passphraseBackedUp: backed }),
      reset: () => set({
        isAuthenticated: false,
        userEmail: null,
        onboardingComplete: false,
        pinCreated: false,
        passphraseBackedUp: false,
      }),
    }),
    {
      name: 'auth-storage',
      storage: createJSONStorage(() => secureStorage),
    }
  )
);