import React from "react";
import { View, Text } from "react-native";
import { useTheme } from "@/theme/ThemeProvider";
import { Icon } from "@/components/ui/Icon";

interface ResultsHeaderProps {
  status: "verified" | "rejected" | "pending_review";
  verifierName: string;
}

export function ResultsHeader({ status, verifierName }: ResultsHeaderProps) {
  const theme = useTheme();

  const getStatusConfig = () => {
    switch (status) {
      case "verified":
        return {
          backgroundColor: theme.colors.primaryGreen + "20",
          iconName: "check-circle",
          iconColor: theme.colors.primaryGreen,
          title: "Verification Complete!",
          description: `Your identity has been verified by ${verifierName}`,
        };
      case "rejected":
        return {
          backgroundColor: "#EF4444" + "20",
          iconName: "x-circle",
          iconColor: "#EF4444",
          title: "Verification Failed",
          description: "Please review the details below",
        };
      default:
        return {
          backgroundColor: "#F59E0B" + "20",
          iconName: "clock",
          iconColor: "#F59E0B",
          title: "Under Review",
          description: "Your verification is being processed",
        };
    }
  };

  const config = getStatusConfig();

  return (
    <View style={{ alignItems: "center", marginBottom: 32 }}>
      <View
        style={{
          width: 80,
          height: 80,
          backgroundColor: config.backgroundColor,
          borderRadius: 40,
          alignItems: "center",
          justifyContent: "center",
          marginBottom: 16,
        }}
      >
        <Icon name={config.iconName} size={40} color={config.iconColor} />
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
        {config.title}
      </Text>

      <Text
        style={{
          fontSize: 16,
          color: theme.colors.textSecondary,
          textAlign: "center",
          lineHeight: 24,
        }}
      >
        {config.description}
      </Text>
    </View>
  );
}
