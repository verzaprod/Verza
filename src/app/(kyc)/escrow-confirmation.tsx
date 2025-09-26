import React, { useState } from "react";
import {
  View,
  Text,
  SafeAreaView,
  ScrollView,
  Alert,
  ActivityIndicator,
  ToastAndroid,
} from "react-native";
import { useRouter, useLocalSearchParams } from "expo-router";
import { useTheme } from "@/theme/ThemeProvider";
import { useSafeAreaInsets } from "react-native-safe-area-context";
import { Icon } from "@/components/ui/Icon";
import { Button } from "@/components/ui/Button";
import { apiService } from "@/services/api/apiService";
// import { LoadingSpinner } from '@/components/ui/LoadingSpinner'

interface VerifierDetails {
  id: string;
  name: string;
  fee: number;
  currency: string;
  description: string;
  estimatedTime: string;
}

export default function EscrowConfirmation() {
  const theme = useTheme();
  const router = useRouter();
  const insets = useSafeAreaInsets();
  const { verifierId } = useLocalSearchParams();
  const [isProcessing, setIsProcessing] = useState(false);

  // Mock verifier data - replace with API call
  const verifierDetails: VerifierDetails = {
    id: verifierId as string,
    name: "TechCorp Solutions",
    fee: 25.0,
    currency: "HBAR",
    description: "Enterprise-grade identity verification with 99.9% accuracy",
    estimatedTime: "2-5 minutes",
  };

  const handleConfirmPayment = async () => {
    setIsProcessing(true);

    try {
      // Initialize escrow transaction
      const escrowResponse = await apiService.initiateEscrow({
        verifier_id: verifierId as string,
        amount: verifierDetails.fee,
        currency: verifierDetails.currency,
        auto_release_hours: 24,
      });

      const escrowData = await escrowResponse.json();

      if (escrowResponse.ok) {
        // Navigate to KYC flow with escrow ID
        ToastAndroid.show("Escrow initiated successfully!", ToastAndroid.SHORT);
        router.push(
          `/(kyc)/selection-type?escrowId=${escrowData.escrow_id}&verifierId=${verifierId}`
        );
      } else {
        Alert.alert(
          "Payment Failed",
          escrowData.error || "Unable to process payment"
        );
      }
    } catch (error) {
      Alert.alert("Error", "Network error occurred. Please try again.");
    } finally {
      setIsProcessing(false);
    }
  };

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
                backgroundColor: theme.colors.primaryGreen + "20",
                borderRadius: 40,
                alignItems: "center",
                justifyContent: "center",
                marginBottom: 16,
              }}
            >
              <Icon name="shield-check" size={40} />
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
              Escrow Payment
            </Text>

            <Text
              style={{
                fontSize: 16,
                color: theme.colors.textSecondary,
                textAlign: "center",
                lineHeight: 24,
              }}
            >
              Secure your verification with {verifierDetails.name}
            </Text>
          </View>

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
              Payment Details
            </Text>

            <View style={{ gap: 12 }}>
              <View className="flex-row justify-between items-center">
                <Text style={{ color: theme.colors.textSecondary }}>
                  Verifier
                </Text>
                <Text
                  style={{ color: theme.colors.textPrimary, fontWeight: "500" }}
                >
                  {verifierDetails.name}
                </Text>
              </View>

              <View className="flex-row justify-between items-center">
                <Text style={{ color: theme.colors.textSecondary }}>
                  Verification Fee
                </Text>
                <Text
                  style={{ color: theme.colors.textPrimary, fontWeight: "600" }}
                >
                  {verifierDetails.fee} {verifierDetails.currency}
                </Text>
              </View>

              <View className="flex-row justify-between items-center">
                <Text style={{ color: theme.colors.textSecondary }}>
                  Processing Time
                </Text>
                <Text
                  style={{ color: theme.colors.textPrimary, fontWeight: "500" }}
                >
                  {verifierDetails.estimatedTime}
                </Text>
              </View>
            </View>
          </View>

          <View
            style={{
              backgroundColor: theme.colors.primaryGreen + "10",
              borderRadius: theme.borderRadius.md,
              padding: 16,
              marginBottom: 32,
            }}
          >
            <View className="flex-row items-start">
              <Icon
                name="shield"
                size={20}
                style={{ marginTop: 2, marginRight: 12 }}
              />
              <View className="flex-1">
                <Text
                  style={{
                    fontSize: 14,
                    fontWeight: "600",
                    color: theme.colors.primaryGreen,
                    marginBottom: 4,
                  }}
                >
                  Escrow Protection
                </Text>
                <Text
                  style={{
                    fontSize: 14,
                    color: theme.colors.textSecondary,
                    lineHeight: 20,
                  }}
                >
                  Your payment is held in escrow and only released to the
                  verifier upon successful identity verification.
                </Text>
              </View>
            </View>
          </View>

          <Button
            text={
              isProcessing
                ? "Processing..."
                : `Pay ${verifierDetails.fee} ${verifierDetails.currency}`
            }
            onPress={handleConfirmPayment}
            disabled={isProcessing}
            style={{ marginBottom: 16 }}
          />

          <Button
            text="Cancel"
            variant="secondary"
            onPress={() => router.back()}
            disabled={isProcessing}
          />
        </View>

        {isProcessing && <ActivityIndicator />}
      </ScrollView>
    </SafeAreaView>
  );
}
