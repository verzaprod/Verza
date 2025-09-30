import React, { useState } from "react";
import { View, Text, TouchableOpacity } from "react-native";
import { useTheme } from "@/theme/ThemeProvider";
import { Icon } from "@/components/ui/Icon";

interface VCDetails {
  id: string;
  issuer: string;
  issuedDate: string;
  expiryDate: string;
  type: string;
  proofHash: string;
}

interface CredentialDetailsProps {
  vcDetails: VCDetails;
}

export function CredentialDetails({ vcDetails }: CredentialDetailsProps) {
  const theme = useTheme();
  const [showProofDetails, setShowProofDetails] = useState(false);

  const handleCopyHash = () => {
    console.log("Proof hash copied to clipboard", vcDetails.proofHash);
  };

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
        Verifiable Credential
      </Text>

      <View style={{ gap: 12 }}>
        <DetailRow label="ID" value={`${vcDetails.id.substring(0, 20)}...`} theme={theme} />
        <DetailRow label="Type" value={vcDetails.type} theme={theme} />
        <DetailRow label="Issuer" value={vcDetails.issuer} theme={theme} />
        <DetailRow label="Issued" value={vcDetails.issuedDate} theme={theme} />
        <DetailRow label="Expires" value={vcDetails.expiryDate} theme={theme} />
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
        <Text style={{ color: theme.colors.primaryGreen, fontWeight: "500" }}>
          Proof Hash
        </Text>
        <Icon
          name={showProofDetails ? "chevron-up" : "chevron-down"}
          size={16}
          color={theme.colors.primaryGreen}
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
            {vcDetails.proofHash}
          </Text>

          <TouchableOpacity
            style={{ marginTop: 8, alignSelf: "flex-end" }}
            onPress={handleCopyHash}
          >
            <Icon name="copy" size={16} color={theme.colors.textSecondary} />
          </TouchableOpacity>
        </View>
      )}
    </View>
  );
}

function DetailRow({ label, value, theme }) {
  return (
    <View style={{ flexDirection: "row", justifyContent: "space-between" }}>
      <Text style={{ color: theme.colors.textSecondary }}>{label}</Text>
      <Text style={{ color: theme.colors.textPrimary, fontWeight: "500" }}>
        {value}
      </Text>
    </View>
  );
}
