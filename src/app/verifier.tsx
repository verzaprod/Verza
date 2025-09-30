import React from "react";
import { SafeAreaView, View, Text, TouchableOpacity, ToastAndroid } from "react-native";
import { useSafeAreaInsets } from "react-native-safe-area-context";
import { useRouter } from "expo-router";
import { useTheme } from "@/theme/ThemeProvider";
import { useKYCStore } from "@/store/kycStore";

export default function VerifierJobDetail() {
  const insets = useSafeAreaInsets();
  const theme = useTheme();
  const router = useRouter();

  const verificationStatus = useKYCStore((state) => state.verificationStatus);
  const setVerificationStatus = useKYCStore((state) => state.setVerificationStatus);

  const handleApprove = () => {
    setVerificationStatus("verified");
    ToastAndroid.show("Verified", ToastAndroid.SHORT);
    router.replace("/(tabs)/home");
  };

  const handleReject = () => {
    setVerificationStatus("rejected");
    ToastAndroid.show("Rejected", ToastAndroid.SHORT);
    router.replace("/(tabs)/home");
  };

  return (
    <SafeAreaView
      style={{
        flex: 1,
        paddingTop: insets.top,
        backgroundColor: theme.colors.background,
        paddingHorizontal: 20,
      }}
    >
      <View style={{ marginTop: 40 }}>
        <Text
          style={{
            fontSize: 22,
            fontWeight: "bold",
            color: theme.colors.textPrimary,
            marginBottom: 20,
            textAlign: "center",
          }}
        >
          Job Detail
        </Text>

        {/* Mock submission details */}
        <View
          style={{
            padding: 16,
            backgroundColor: theme.colors.card,
            borderRadius: 12,
            marginBottom: 30,
          }}
        >
          <Text style={{ color: theme.colors.textSecondary, marginBottom: 8 }}>
            Requester submitted documents for verification.
          </Text>
          <Text style={{ color: theme.colors.textPrimary }}>
            Current Status:{" "}
            <Text style={{ fontWeight: "bold" }}>{verificationStatus}</Text>
          </Text>
        </View>

        {/* Approve / Reject buttons */}
        <TouchableOpacity
          style={{
            backgroundColor: "green",
            padding: 14,
            borderRadius: 12,
            marginBottom: 16,
          }}
          onPress={handleApprove}
        >
          <Text style={{ color: "#fff", textAlign: "center", fontWeight: "bold" }}>
            Approve
          </Text>
        </TouchableOpacity>

        <TouchableOpacity
          style={{
            backgroundColor: "red",
            padding: 14,
            borderRadius: 12,
          }}
          onPress={handleReject}
        >
          <Text style={{ color: "#fff", textAlign: "center", fontWeight: "bold" }}>
            Reject
          </Text>
        </TouchableOpacity>
      </View>
    </SafeAreaView>
  );
}
