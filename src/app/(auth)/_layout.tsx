import { useAuth } from "@clerk/clerk-expo";
import { Redirect, Stack } from "expo-router";
import { useAuthStore } from "@/store/authStore";

export default function AuthLayout() {
  
  const { isSignedIn } = useAuth();
  const { pinCreated, passphraseBackedUp } = useAuthStore();

  if (isSignedIn) {

    const hasCompletedOnboarding = pinCreated && passphraseBackedUp;

    if (hasCompletedOnboarding) {
      return <Redirect href={"/(tabs)/home"} />
    }
  }

  return (
    <>
      <Stack screenOptions={{ headerShown: false }}>
        <Stack.Screen name="register" options={{ headerShown: false }} />
        <Stack.Screen name="sign-in" options={{ headerShown: false }} />
        <Stack.Screen name="auth-success" options={{ headerShown: false }} />
        <Stack.Screen name="backup-passphrase" options={{ headerShown: false }} />
        <Stack.Screen name="confirm-passphrase" options={{ headerShown: false }} />
        <Stack.Screen name="create-pin" options={{ headerShown: false }} />
        <Stack.Screen name="verify-email" options={{ headerShown: false }} />
      </Stack>
      </>
  )
}
