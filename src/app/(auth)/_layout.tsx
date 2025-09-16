import { Stack } from "expo-router";

export default function AuthLayout() {
  return (
    <Stack >
      <Stack.Screen name="index" options={{ headerShown: false }} />
      <Stack.Screen name="auth-success" options={{ headerShown: false }} />
      <Stack.Screen name="backup-passphrase" options={{ headerShown: false }} />
      <Stack.Screen name="confirm-passphrase" options={{ headerShown: false }} />
      <Stack.Screen name="create-pin" options={{ headerShown: false }} />
      <Stack.Screen name="verify-email" options={{ headerShown: false }} />
    </Stack>
  )
}
