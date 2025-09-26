import { Stack } from "expo-router";

export default function KYCLayout() {
  return (
    <Stack
      screenOptions={{
        animation: "slide_from_right",
        headerShown: false,
      }}
    >
      <Stack.Screen name="selection-type" />
      <Stack.Screen name="doc-capture" />
      <Stack.Screen name="selfie-capture" />
      <Stack.Screen name="selfie-note" />
      <Stack.Screen name="kyc-success" />
      <Stack.Screen name="verification-state" />
      <Stack.Screen name="escrow-confirmation" />
      <Stack.Screen name="verification-results" />
      <Stack.Screen name="verification-tracker" />
    </Stack>
  );
}
