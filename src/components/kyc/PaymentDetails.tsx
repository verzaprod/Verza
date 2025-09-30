import React from "react";
import { View, Text } from "react-native";
import { useTheme } from "@/theme/ThemeProvider";

interface VerifierDetails {
  name: string;
  fee: number;
  currency: string;
  estimatedTime: string;
}

interface PaymentDetailsProps {
  verifierDetails: VerifierDetails;
}

export function PaymentDetails({ verifierDetails }: PaymentDetailsProps) {
  const theme = useTheme();

  return (
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
        <DetailRow
          label="Verifier"
          value={verifierDetails.name}
          theme={theme}
        />
        <DetailRow
          label="Verification Fee"
          value={`${verifierDetails.fee} ${verifierDetails.currency}`}
          theme={theme}
          valueWeight="600"
        />
        <DetailRow
          label="Processing Time"
          value={verifierDetails.estimatedTime}
          theme={theme}
        />
      </View>
    </View>
  );
}

function DetailRow({ label, value, theme, valueWeight = "500" }) {
  return (
    <View style={{ flexDirection: "row", justifyContent: "space-between", alignItems: "center" }}>
      <Text style={{ color: theme.colors.textSecondary }}>{label}</Text>
      <Text style={{ color: theme.colors.textPrimary, fontWeight: valueWeight }}>
        {value}
      </Text>
    </View>
  );
}