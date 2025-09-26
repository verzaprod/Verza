import React, { useState, useEffect } from "react";
import {
  View,
  Text,
  SafeAreaView,
  ScrollView,
  TouchableOpacity,
} from "react-native";
import { useRouter, useLocalSearchParams } from "expo-router";
import { useTheme } from "@/theme/ThemeProvider";
import { useSafeAreaInsets } from "react-native-safe-area-context";
import { Icon } from "@/components/ui/Icon";
import { Button } from "@/components/ui/Button";
import { apiService } from "@/services/api/apiService";

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
  const { escrowId } = useLocalSearchParams();
  const [result, setResult] = useState<VerificationResult | null>(null);
  const [showProofDetails, setShowProofDetails] = useState(false);

  useEffect(() => {
    const fetchResults = async () => {
      try {
        const response = await apiService.getVerificationResults(
          escrowId as string
        );
        const data = await response.json();
        setResult(data);
      } catch (error) {
        console.error("Failed to fetch results:", error);
      }
    };

    fetchResults();
  }, [escrowId]);

  const handleCopyHash = () => {
    // Copy proof hash to clipboard
    if (result?.vcDetails?.proofHash) {
      // Clipboard.setString(result.vcDetails.proofHash)
      // Show toast notification
    }
  };

  if (!result) {
    return (
      <SafeAreaView
        style={{ flex: 1, backgroundColor: theme.colors.background }}
      >
        <View className="flex-1 justify-center items-center">
          <Text style={{ color: theme.colors.textSecondary }}>
            Loading results...
          </Text>
        </View>
      </SafeAreaView>
    );
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
          <View style={{ alignItems: "center", marginBottom: 32 }}>
            <View
              style={{
                width: 80,
                height: 80,
                backgroundColor:
                  result.status === "verified"
                    ? theme.colors.primaryGreen + "20"
                    : "#EF4444" + "20",
                borderRadius: 40,
                alignItems: "center",
                justifyContent: "center",
                marginBottom: 16,
              }}
            >
              <Icon
                name={result.status === "verified" ? "success" : "cancel"}
                size={40}
              />
            </View>

            <Text
              style={{
                fontSize: 24,
                fontWeight: "bold",
                color: theme.colors.textPrimary,
                fontFamily: theme.fonts.welcomeHeading,
                textAlign: "center",
                marginBottom: 8,
              }}
            >
              {result.status === "verified"
                ? "Verification Complete!"
                : "Verification Failed"}
            </Text>

            <Text
              style={{
                fontSize: 16,
                color: theme.colors.textSecondary,
                textAlign: "center",
              }}
            >
              {result.status === "verified"
                ? "Your digital credential is ready"
                : "Please review the details below"}
            </Text>
          </View>

          {result.status === "verified" && result.vcDetails && (
            <View
              style={{
                backgroundColor: theme.colors.backgroundLight,
                borderRadius: theme.borderRadius.lg,
                padding: 20,
                marginBottom: 24,
              }}
            >
              <Text
                style={{
                  fontSize: 18,
                  fontWeight: "600",
                  color: theme.colors.textPrimary,
                  marginBottom: 16,
                }}
              >
                Verifiable Credential
              </Text>

              <View style={{ gap: 12 }}>
                <View className="flex-row justify-between">
                  <Text style={{ color: theme.colors.textSecondary }}>ID</Text>
                  <Text
                    style={{
                      color: theme.colors.textPrimary,
                      fontFamily: "monospace",
                      fontSize: 12,
                    }}
                  >
                    {result.vcDetails.id.substring(0, 20)}...
                  </Text>
                </View>

                <View className="flex-row justify-between">
                  <Text style={{ color: theme.colors.textSecondary }}>
                    Type
                  </Text>
                  <Text
                    style={{
                      color: theme.colors.textPrimary,
                      fontWeight: "500",
                    }}
                  >
                    {result.vcDetails.type}
                  </Text>
                </View>

                <View className="flex-row justify-between">
                  <Text style={{ color: theme.colors.textSecondary }}>
                    Issuer
                  </Text>
                  <Text
                    style={{
                      color: theme.colors.textPrimary,
                      fontWeight: "500",
                    }}
                  >
                    {result.vcDetails.issuer}
                  </Text>
                </View>

                <View className="flex-row justify-between">
                  <Text style={{ color: theme.colors.textSecondary }}>
                    Issued
                  </Text>
                  <Text style={{ color: theme.colors.textPrimary }}>
                    {result.vcDetails.issuedDate}
                  </Text>
                </View>

                <View className="flex-row justify-between">
                  <Text style={{ color: theme.colors.textSecondary }}>
                    Expires
                  </Text>
                  <Text style={{ color: theme.colors.textPrimary }}>
                    {result.vcDetails.expiryDate}
                  </Text>
                </View>
              </View>

              <TouchableOpacity
                style={{
                  marginTop: 16,
                  padding: 12,
                  backgroundColor: theme.colors.primaryGreen + "10",
                  borderRadius: theme.borderRadius.md,
                  flexDirection: "row",
                  alignItems: "center",
                  justifyContent: "space-between",
                }}
                onPress={() => setShowProofDetails(!showProofDetails)}
              >
                <Text
                  style={{
                    color: theme.colors.primaryGreen,
                    fontWeight: "500",
                  }}
                >
                  Proof Hash
                </Text>
                <Icon
                  name={showProofDetails ? "chevron-up" : "chevron-down"}
                  size={16}
                />
              </TouchableOpacity>

              {showProofDetails && (
                <View
                  style={{
                    marginTop: 12,
                    padding: 12,
                    backgroundColor: theme.colors.background,
                    borderRadius: theme.borderRadius.md,
                  }}
                >
                  <Text
                    style={{
                      fontFamily: "monospace",
                      fontSize: 12,
                      color: theme.colors.textPrimary,
                      lineHeight: 16,
                    }}
                  >
                    {result.vcDetails.proofHash}
                  </Text>

                  <TouchableOpacity
                    style={{
                      marginTop: 8,
                      alignSelf: "flex-end",
                    }}
                    onPress={handleCopyHash}
                  >
                    <Icon name="copy" size={16} />
                  </TouchableOpacity>
                </View>
              )}
            </View>
          )}

          {result.status === "rejected" && (
            <View
              style={{
                backgroundColor: "#EF4444" + "10",
                borderRadius: theme.borderRadius.lg,
                padding: 20,
                marginBottom: 24,
              }}
            >
              <Text
                style={{
                  fontSize: 16,
                  fontWeight: "600",
                  color: "#EF4444",
                  marginBottom: 8,
                }}
              >
                Rejection Reason
              </Text>
              <Text
                style={{ color: theme.colors.textSecondary, lineHeight: 20 }}
              >
                {result.rejectionReason}
              </Text>
            </View>
          )}

          <Button
            text={result.status === "verified" ? "Go to Wallet" : "Try Again"}
            onPress={() => {
              if (result.status === "verified") {
                router.push("/(tabs)/profile");
              } else {
                router.push("/(tabs)/verifiers");
              }
            }}
            style={{ marginBottom: 16 }}
          />

          <Button
            text="Back to Home"
            variant="secondary"
            onPress={() => router.push("/(tabs)/home")}
          />
        </View>
      </ScrollView>
    </SafeAreaView>
  );
}
