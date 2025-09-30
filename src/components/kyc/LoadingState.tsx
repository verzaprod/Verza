import React from "react";
import { SafeAreaView, ActivityIndicator, Text } from "react-native";
import { useTheme } from "@/theme/ThemeProvider";

export function LoadingState() {
  const theme = useTheme();

  return (
    <SafeAreaView
      style={{
        flex: 1,
        backgroundColor: theme.colors.background,
        justifyContent: 'center',
        alignItems: 'center',
      }}
    >
      <ActivityIndicator size="large" color={theme.colors.primaryGreen} />
      <Text style={{ 
        marginTop: 16, 
        color: theme.colors.textSecondary,
        fontSize: 16 
      }}>
        Loading verifier details...
      </Text>
    </SafeAreaView>
  );
}
