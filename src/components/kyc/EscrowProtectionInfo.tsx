import React from "react";
import { View, Text } from "react-native";
import { useTheme } from "@/theme/ThemeProvider";
import { Icon } from "@/components/ui/Icon";

export function EscrowProtectionInfo() {
  const theme = useTheme();

  return (
    <View
      style={{
        backgroundColor: theme.colors.primaryGreen + "10",
        borderRadius: theme.borderRadius.md,
        padding: 16,
        marginBottom: 32,
      }}
    >
      <View style={{ flexDirection: "row", alignItems: "flex-start" }}>
        <Icon
          name="shield"
          size={20}
          style={{ marginTop: 2, marginRight: 12 }}
        />
        <View style={{ flex: 1 }}>
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
            Your payment is held in escrow and only released to the verifier upon successful identity verification.
          </Text>
        </View>
      </View>
    </View>
  );
}
