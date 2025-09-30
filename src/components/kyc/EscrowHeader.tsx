import React from "react";
import { View, Text } from "react-native";
import { useTheme } from "@/theme/ThemeProvider";
import { Icon } from "@/components/ui/Icon";
import Feather from "@expo/vector-icons/Feather";

interface EscrowHeaderProps {
  verifierName: string;
}

export function EscrowHeader({ verifierName }: EscrowHeaderProps) {
  const theme = useTheme();

  return (
    <View style={{ paddingTop: 20, alignItems: "center", marginBottom: 32 }}>
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
        <Feather name="shield" size={40} />
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
        Secure your verification with {verifierName}
      </Text>
    </View>
  );
}