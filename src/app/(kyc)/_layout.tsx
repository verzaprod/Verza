import { Stack } from "expo-router";

export default function KYCLayout() {
  return (
    <Stack>
      <Stack.Screen name="selection-type" options={{ headerShown: false }} />
      <Stack.Screen name="doc-capture" options={{ headerShown: false }} />
      <Stack.Screen name="selfie-capture" options={{ headerShown: false }} />
      <Stack.Screen name="selfie-note" options={{ headerShown: false }} />
      <Stack.Screen name="kyc-success" options={{ headerShown: false }} />
      <Stack.Screen
        name="verification-state"
        options={{ headerShown: false }}
      />
    </Stack>
  );
}
