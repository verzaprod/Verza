import React, { useState, useEffect } from "react";
import { View, SafeAreaView, ScrollView } from "react-native";
import { useRouter, useLocalSearchParams } from "expo-router";
import { useTheme } from "@/theme/ThemeProvider";
import { useSafeAreaInsets } from "react-native-safe-area-context";
import { apiService } from "@/services/api/apiService";
import { useKYCStore } from "@/store/kycStore";
import { useRatingStore } from "@/store/ratingStore";
import { MOCK_DATA } from "@/services/api/mockData";
import { ResultsHeader } from "@/components/verification/ResultsHeader";
import { CredentialDetails } from "@/components/verification/CredentialDetails";
import { RejectionDetails } from "@/components/verification/RejectionDetails";
import { VerifierRating } from "@/components/verifiers/VerifierRating";
import { ResultsActions } from "@/components/verification/ResultsActions";
import { LoadingState } from "@/components/verification/LoadingState";

interface VerificationResult {
  status: "verified" | "rejected" | "pending_review";
  credentialId?: string;
  vcDetails?: {
    id: string;
    issuer: string;
    issuedDate: string;
    expiryDate: string;
    type: string;
    proofHash: string;
  };
  rejectionReason?: string;
}

export default function VerificationResults() {
  const theme = useTheme();
  const router = useRouter();
  const insets = useSafeAreaInsets();
  const { verifierId } = useLocalSearchParams();
  const [result, setResult] = useState<VerificationResult | null>(null);
  const escrowId = useKYCStore((state) => state.escrowId);
  const { loadRatings } = useRatingStore();

  useEffect(() => {
    loadRatings();
    fetchResults();
  }, [escrowId]);

  const fetchResults = async () => {
    try {
      if (!escrowId) return;
      const response = await apiService.getVerificationResults(escrowId as string);
      const data = await response.json();
      setResult(data);
    } catch (error) {
      console.error("Failed to fetch results:", error);
    }
  };

  const verifier = MOCK_DATA.verifiers.find(v => v.id === verifierId);
  const verifierName = verifier ? verifier?.name : 'Unknown Verifier';

  if (!result) {
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
        <View style={{ paddingTop: 20, paddingBottom: 40 }}>
          <ResultsHeader 
            status={result.status} 
            verifierName={verifierName} 
          />

          {result.status === "verified" && result.vcDetails && (
            <CredentialDetails vcDetails={result.vcDetails} />
          )}

          {result.status === "rejected" && (
            <RejectionDetails rejectionReason={result.rejectionReason} />
          )}

          {result.status === "verified" && (
            <VerifierRating 
              verifierId={verifierId as string}
              verifierName={verifierName}
            />
          )}

          <ResultsActions 
            status={result.status}
            onNavigate={(route) => router.push(route as any)}
          />
        </View>
      </ScrollView>
    </SafeAreaView>
  );
}
