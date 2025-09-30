// import { useKYCStore } from "@/store/kycStore";
import { Stack } from "expo-router";

export default function KYCLayout() {

  // const { resetKYC } = useKYCStore();

  // resetKYC();

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
      <Stack.Screen name="escrow-confirmation" />
      <Stack.Screen name="verification-results" />
      <Stack.Screen name="verification-tracker" />
    </Stack>
  );
}
