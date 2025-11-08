import React, { useState, useEffect, useCallback, useRef } from "react";
import { View, Text, SafeAreaView, ScrollView } from "react-native";
import { useRouter } from "expo-router";
import { useTheme } from "@/theme/ThemeProvider";
import { useSafeAreaInsets } from "react-native-safe-area-context";
import { Icon } from "@/components/ui/Icon";
import { Button } from "@/components/ui/Button";
import { apiService } from "@/services/api/apiService";
import { useKYCStore } from "@/store/kycStore";

interface VerificationStatus {
  escrowId: string;
  status: "submitted" | "in_progress" | "completed" | "failed";
  steps: {
    id: string;
    label: string;
    status: "completed" | "active" | "pending";
    timestamp?: string;
  }[];
  estimatedCompletion?: string;
}

export default function VerificationTracker() {
  const theme = useTheme();
  const router = useRouter();
  const insets = useSafeAreaInsets();
  const [verificationStatus, setVerificationStatus] =
    useState<VerificationStatus | null>(null);

  const escrowId = useKYCStore((state) => state.escrowId);
  const setCurrentStep = useKYCStore((state) => state.setCurrentStep);

  const intervalRef = useRef<NodeJS.Timeout | null>(null);
  const hasNavigatedRef = useRef(false);
  const isMountedRef = useRef(true);

  const pollStatus = useCallback(async () => {
    if (!escrowId || !isMountedRef.current || hasNavigatedRef.current) {
      return;
    }

    try {
      const response = await apiService.getVerificationStatus(escrowId);
      const data = await response.json();

      if (isMountedRef.current && !hasNavigatedRef.current) {
        setVerificationStatus(data);

        if (data.status === "completed" && !hasNavigatedRef.current) {
          hasNavigatedRef.current = true;
          setCurrentStep("complete");

          // ✅ Clear interval before navigation
          if (intervalRef.current) {
            clearInterval(intervalRef.current);
            intervalRef.current = null;
          }

          setTimeout(() => {
            if (isMountedRef.current) {
              router.replace("/(kyc)/verification-results"); 
            }
          }, 2000);
        }
      }
    } catch (error) {
      console.error("Failed to fetch status:", error);
    }
  }, [escrowId, setCurrentStep, router]);

  useEffect(() => {
    isMountedRef.current = true;
    hasNavigatedRef.current = false;

    if (escrowId) {
      // Initial call
      pollStatus();

      intervalRef.current = setInterval(pollStatus, 5000);
    }

    return () => {
      isMountedRef.current = false;
      if (intervalRef.current) {
        clearInterval(intervalRef.current);
        intervalRef.current = null;
      }
    };
  }, [escrowId, pollStatus]);

  useEffect(() => {
    return () => {
      isMountedRef.current = false;
    };
  }, []);

  const getStepIcon = useCallback((status: string) => {
    switch (status) {
      case "completed":
        return "✓";
      case "active":
        return "⋯";
      case "pending":
        return "○";
      default:
        return "○";
    }
  }, []);

  const getStepColor = useCallback(
    (status: string) => {
      switch (status) {
        case "completed":
          return theme.colors.primaryGreen;
        case "active":
          return "#F59E0B";
        case "pending":
          return theme.colors.textSecondary;
        default:
          return theme.colors.textSecondary;
      }
    },
    [theme.colors.primaryGreen, theme.colors.textSecondary]
  );

  const handleViewResults = useCallback(() => {
    if (!hasNavigatedRef.current) {
      hasNavigatedRef.current = true;
      router.replace("/(kyc)/verification-results");
    }
  }, [router]);

  if (!verificationStatus) {
    return (
      <SafeAreaView
        style={{ flex: 1, backgroundColor: theme.colors.background }}
      >
        <View
          style={{ flex: 1, justifyContent: "center", alignItems: "center" }}
        >
          <Text style={{ color: theme.colors.textSecondary }}>Loading...</Text>
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
              Verification Progress
            </Text>

            <Text
              style={{
                fontSize: 16,
                color: theme.colors.textSecondary,
                textAlign: "center",
              }}
            >
              Tracking your identity verification
            </Text>
          </View>

          <View
            style={{
              backgroundColor: theme.colors.background,
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
                marginBottom: 20,
              }}
            >
              Verification Steps
            </Text>

            <View style={{ gap: 16 }}>
              {verificationStatus.steps.map((step) => (
                <View
                  key={step.id}
                  style={{ flexDirection: "row", alignItems: "center" }}
                >
                  <View
                    style={{
                      width: 32,
                      height: 32,
                      borderRadius: 16,
                      backgroundColor: getStepColor(step.status),
                      alignItems: "center",
                      justifyContent: "center",
                      marginRight: 16,
                    }}
                  >
                    <Text
                      style={{
                        color: "white",
                        fontWeight: "bold",
                        fontSize: 16,
                      }}
                    >
                      {getStepIcon(step.status)}
                    </Text>
                  </View>

                  <View style={{ flex: 1 }}>
                    <Text
                      style={{
                        fontSize: 16,
                        fontWeight: "500",
                        color: theme.colors.textPrimary,
                        marginBottom: 2,
                      }}
                    >
                      {step.label}
                    </Text>
                    {step.timestamp && (
                      <Text
                        style={{
                          fontSize: 14,
                          color: theme.colors.textSecondary,
                        }}
                      >
                        {step.timestamp}
                      </Text>
                    )}
                  </View>
                </View>
              ))}
            </View>

            {verificationStatus.estimatedCompletion && (
              <View
                style={{
                  marginTop: 20,
                  padding: 12,
                  backgroundColor: theme.colors.primaryGreen + "10",
                  borderRadius: theme.borderRadius.md,
                }}
              >
                <Text
                  style={{
                    fontSize: 14,
                    color: theme.colors.primaryGreen,
                    textAlign: "center",
                  }}
                >
                  Estimated completion: {verificationStatus.estimatedCompletion}
                </Text>
              </View>
            )}
          </View>

          {verificationStatus.status === "completed" && (
            <Button text="View Results" onPress={handleViewResults} />
          )}
        </View>
      </ScrollView>
    </SafeAreaView>
  );
}
