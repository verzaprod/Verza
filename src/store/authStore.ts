import { create } from 'zustand';
import { persist, createJSONStorage } from 'zustand/middleware';
import * as SecureStore from 'expo-secure-store';

interface AuthState {
  isAuthenticated: boolean;
  userEmail: string | null;
  onboardingComplete: boolean;
  pinCreated: boolean;
  passphraseBackedUp: boolean;
  isFirstTimeUser: boolean;
  setEmail: (email: string) => void;
  setOnboardingComplete: (complete: boolean) => void;
  setPinCreated: (created: boolean) => void;
  setPassphraseBackedUp: (backed: boolean) => void;
  setFirstTimeUser: (isFirstTime: boolean) => void;
  setAuthenticated: (authenticated: boolean) => void;
  completeOnboarding: () => void;
  reset: () => void;
}

const secureStorage = {
  getItem: async (name: string): Promise<string | null> => {
    try {
      return await SecureStore.getItemAsync(name);
    } catch (error) {
      console.warn("SecureStore getItem erorr", error);
    } 
  },
  setItem: async (name: string, value: string): Promise<void> => {
    try {
      await SecureStore.setItemAsync(name, value);
    } catch (error) {
      console.warn("SecureStore getItem erorr", error);
    } 
  },
  removeItem: async (name: string): Promise<void> => {
    try {
      await SecureStore.deleteItemAsync(name);
    } catch (error) {
      console.warn("SecureStore getItem erorr", error);
    } ;
  },
};

export const useAuthStore = create<AuthState>()(
  persist(
    (set, get) => ({
      isAuthenticated: false,
      userEmail: null,
      onboardingComplete: false,
      pinCreated: false,
      passphraseBackedUp: false,
      isFirstTimeUser: true,
      setEmail: (email) => set({ userEmail: email }),
      setOnboardingComplete: (complete) => set({ onboardingComplete: complete }),
      setPinCreated: (created) => {
        console.log("Setting pinCreated:", created);
        set({ pinCreated: created });
      },
      setPassphraseBackedUp: (backed) => {
        console.log("Setting PassphraseBackedup:", backed);
        set({ passphraseBackedUp: backed });
      },
      setFirstTimeUser: (isFirstTime) => set({ isFirstTimeUser: isFirstTime}),
      setAuthenticated: (authenticated) => set({isAuthenticated: authenticated}),
      completeOnboarding: () => set({
        onboardingComplete: true, 
        pinCreated: true,
        passphraseBackedUp: true,
        isFirstTimeUser: false,
      }),
      reset: () => set({
        isAuthenticated: false,
        userEmail: null,
        onboardingComplete: false,
        pinCreated: false,
        passphraseBackedUp: false,
        isFirstTimeUser: true,
      }),
    }),
    {
      name: 'auth-storage',
      storage: createJSONStorage(() => secureStorage),
    }
  )
);