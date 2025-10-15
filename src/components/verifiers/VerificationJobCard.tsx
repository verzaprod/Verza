import React from "react";
import { View, Text, TouchableOpacity } from "react-native";
import { Ionicons } from "@expo/vector-icons";
import { useTheme } from "@/theme/ThemeProvider";

export default function VerificationJobCard({ job, onPress }) {
  const theme = useTheme();

  return (
    <TouchableOpacity
      style={{
        marginHorizontal: 20,
        marginBottom: 16,
        borderRadius: theme.borderRadius.lg * 1.5,
        backgroundColor: theme.colors.background,
        borderWidth: theme.isDark ? 1 : 0,
        borderColor: theme.isDark ? theme.colors.boxBorder : "transparent",
        shadowColor: theme.isDark ? "#fff" : "#000",
        shadowOffset: { width: 0, height: 4 },
        shadowOpacity: theme.isDark ? 0.1 : 0.08,
        shadowRadius: 12,
        elevation: 4,
      }}
      onPress={onPress}
      activeOpacity={0.7}
    >
      <View style={{ 
        flexDirection: "row", 
        alignItems: "center", 
        padding: theme.spacing.lg 
      }}>
        {/* Avatar */}
        <View 
          style={{
            width: 64,
            height: 64,
            borderRadius: 32,
            backgroundColor: theme.isDark 
              ? `${theme.colors.textSecondary}30` 
              : `${theme.colors.textSecondary}20`,
            alignItems: "center",
            justifyContent: "center",
          }}
        >
          <Ionicons 
            name="person" 
            size={32} 
            color={theme.colors.textSecondary} 
          />
        </View>

        {/* Info */}
        <View style={{ flex: 1, marginLeft: theme.spacing.md }}>
          <Text 
            style={{
              fontSize: 20,
              fontWeight: "600",
              color: theme.colors.textPrimary,
              marginBottom: 4,
            }}
          >
            {job.requester}
          </Text>
          <Text 
            style={{
              fontSize: 15,
              color: theme.colors.textSecondary,
            }}
          >
            {job.doc}
          </Text>
        </View>

        {/* Document Icon */}
        <View 
          style={{
            width: 48,
            height: 48,
            borderRadius: theme.borderRadius.md,
            backgroundColor: theme.isDark 
              ? `${theme.colors.textSecondary}20` 
              : `${theme.colors.textSecondary}15`,
            alignItems: "center",
            justifyContent: "center",
          }}
        >
          <Ionicons 
            name="card-outline" 
            size={24} 
            color={theme.colors.textSecondary} 
          />
        </View>
      </View>
    </TouchableOpacity>
  );
}