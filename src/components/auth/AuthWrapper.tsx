import React, { useEffect } from 'react';
import { useRouter, useSegments } from 'expo-router';
import { useAuth } from '@clerk/clerk-expo';
import { useAuthStore } from '@/store/authStore';

interface AuthWrapperProps {
  children: React.ReactNode;
}

export const AuthWrapper: React.FC<AuthWrapperProps> = ({ children }) => {
  const { isSignedIn, isLoaded } = useAuth();
  const router = useRouter();
  const segments = useSegments();
  
  const {
    onboardingComplete,
    pinCreated,
    passphraseBackedUp,
    isFirstTimeUser,
    setAuthenticated,
    setFirstTimeUser
  } = useAuthStore();

  useEffect(() => {
    if (!isLoaded) return;

    const inAuthGroup = segments[0] === '(auth)';
    const inTabsGroup = segments[0] === '(tabs)';
    const inKYCGroup = segments[0] === '(kyc)';

    console.log('Auth state:', {
      isSignedIn,
      onboardingComplete,
      pinCreated,
      passphraseBackedUp,
      isFirstTimeUser,
      segments
    });

    if (isSignedIn) {
      setAuthenticated(true);
      
      // Check if user has completed full onboarding
      const hasCompletedOnboarding = pinCreated && passphraseBackedUp;
      
      if (hasCompletedOnboarding) {
        // Returning user - go to home
        if (inAuthGroup) {
          router.replace('/(tabs)/home');
        }
      } else {
        // First time user or incomplete onboarding
        if (isFirstTimeUser) {
          // Start onboarding process
          if (!pinCreated) {
            router.replace('/(auth)/create-pin');
          } else if (!passphraseBackedUp) {
            router.replace('/(auth)/backup-passphrase');
          }
        }
      }
    } else {
      setAuthenticated(false);
      // Not signed in - redirect to auth
      if (inTabsGroup || inKYCGroup) {
        router.replace('/(auth)/sign-in');
      }
    }
  }, [
    isSignedIn, 
    isLoaded, 
    segments, 
    onboardingComplete, 
    pinCreated, 
    passphraseBackedUp,
    isFirstTimeUser
  ]);

  return <>{children}</>;
};