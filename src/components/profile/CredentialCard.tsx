import React from "react";
import { View, Text, TouchableOpacity } from "react-native";
import { useTheme } from "@/theme/ThemeProvider";
import Feather from "@expo/vector-icons/Feather";
import { Icon } from "../ui/Icon";

interface CredentialCardProps {
  type: string;
  status: "verified" | "pending";
  // icon: React.ComponentProps<typeof Feather>['name'];
  icon: string;
}

export const CredentialCard: React.FC<CredentialCardProps> = ({
  type,
  status,
  icon,
}) => {
  const theme = useTheme();
  const isVerified = status === "verified";

  return (
    <TouchableOpacity
      style={{
        flexDirection: "row",
        alignItems: "center",
        justifyContent: "space-between",
        padding: theme.spacing.lg,
        backgroundColor: theme.colors.background,
        borderRadius: theme.borderRadius.lg,
        shadowColor: theme.isDark ? "#fff" : "#000",
        shadowOffset: { width: 0, height: 2 },
        shadowOpacity: 0.2,
        shadowRadius: 4,
        elevation: 4,
      }}
    >
      <View className="flex-row items-center flex-1">
        <View
          style={{
            width: 48,
            height: 48,
            backgroundColor: theme.colors.textPrimary,
            borderRadius: 8,
            alignItems: "center",
            justifyContent: "center",
            marginRight: theme.spacing.md,
          }}
        >
          <Icon name={icon} size={24} />
        </View>

        <Text
          style={{
            fontSize: 18,
            fontWeight: "600",
            color: theme.colors.textPrimary,
          }}
        >
          {type}
        </Text>
      </View>

      <View className="flex-row items-center">
        <Text
          style={{
            fontSize: 16,
            color: isVerified ? theme.colors.primaryGreen : "#F59E0B",
            fontWeight: "500",
            marginRight: theme.spacing.sm,
            fontStyle: "italic",
          }}
        >
          {isVerified ? "Verified" : "Pending"}
        </Text>

        {isVerified && (
          <View
            style={{
              width: 24,
              height: 24,
              backgroundColor: theme.colors.primaryGreen,
              borderRadius: 12,
              alignItems: "center",
              justifyContent: "center",
            }}
          >
            <Text style={{ color: "white", fontSize: 14, fontWeight: "bold" }}>
              ✓
            </Text>
          </View>
        )}

        {!isVerified && (
          <View
            style={{
              width: 24,
              height: 24,
              backgroundColor: "#F59E0B",
              borderRadius: 12,
              alignItems: "center",
              justifyContent: "center",
            }}
          >
            <Text style={{ color: "white", fontSize: 14, fontWeight: "bold" }}>
              !
            </Text>
          </View>
        )}
      </View>
    </TouchableOpacity>
  );
};
