import React from "react";
import { View, Text } from "react-native";
import { useTheme } from "@/theme/ThemeProvider";
import { Icon } from "@/components/ui/Icon";

export const ProfileHeader: React.FC = () => {
  const theme = useTheme();

  return (
    <View style={{ alignItems: "center", marginBottom: theme.spacing.xxl }}>
      <View
        style={{
          width: 120,
          height: 120,
          borderRadius: 60,
          backgroundColor: theme.colors.primaryGreen,
          alignItems: "center",
          justifyContent: "center",
          marginBottom: theme.spacing.lg,
        }}
      >
        <Icon name="avatar" size={80} />
      </View>

      <Text
        style={{
          fontSize: 24,
          fontWeight: "bold",
          color: theme.colors.textPrimary,
          fontFamily: theme.fonts.welcomeHeading,
          marginBottom: theme.spacing.xs,
        }}
      >
        did:verza:1234abcd
      </Text>

      <Text
        style={{
          fontSize: 16,
          color: theme.colors.textSecondary,
        }}
      >
        helloworld@gmail.com
      </Text>
    </View>
  );
};
