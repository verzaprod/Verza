import React from "react";
import { SafeAreaView, ScrollView } from "react-native";
import { useRouter, useLocalSearchParams } from "expo-router";
import { useTheme } from "@/theme/ThemeProvider";
import { useSafeAreaInsets } from "react-native-safe-area-context";
import { useVerifierDetails } from "@/hooks/useVerifierDetails";
import { EscrowHeader } from "@/components/kyc/EscrowHeader";
import { PaymentDetails } from "@/components/kyc/PaymentDetails";
import { EscrowProtectionInfo } from "@/components/kyc/EscrowProtectionInfo";
import { PaymentActions } from "@/components/kyc/PaymentActions";
import { LoadingState } from "@/components/kyc/LoadingState";

export default function EscrowConfirmation() {
  const theme = useTheme();
  const router = useRouter();
  const insets = useSafeAreaInsets();
  const { verifierId } = useLocalSearchParams();
  
  const {
    verifierDetails,
    isLoading,
    isProcessing,
    selectedToken,
    updateSelectedToken,
    handleConfirmPayment
  } = useVerifierDetails(verifierId as string);

  if (isLoading || !verifierDetails) {
    return <LoadingState />;
  }

  return (
    <SafeAreaView
      style={{
        flex: 1,
        backgroundColor: theme.colors.background,
        paddingTop: insets.top,
      }}
    >
      <ScrollView
        style={{ paddingHorizontal: 20 }}
        showsVerticalScrollIndicator={false}
      >
        <EscrowHeader verifierName={verifierDetails.name} />
        
        <PaymentDetails verifierDetails={verifierDetails} />
        
        <EscrowProtectionInfo />
        
        <PaymentActions
          verifierDetails={verifierDetails}
          isProcessing={isProcessing}
          selectedToken={selectedToken}
          onTokenChange={updateSelectedToken}
          onConfirmPayment={handleConfirmPayment}
          onCancel={() => router.back()}
        />
      </ScrollView>
    </SafeAreaView>
  );
}