import React from "react";
import { View, Text } from "react-native";
import { useTheme } from "@/theme/ThemeProvider";

interface RejectionDetailsProps {
  rejectionReason?: string;
}

export function RejectionDetails({ rejectionReason }: RejectionDetailsProps) {
  const theme = useTheme();

  return (
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
        style={{ 
          color: theme.colors.textSecondary, 
          lineHeight: 20 
        }}
      >
        {rejectionReason || "No specific reason provided."}
      </Text>
    </View>
  );
}
