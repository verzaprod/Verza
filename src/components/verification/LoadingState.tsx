import React from "react";
import { View, Text, SafeAreaView, ActivityIndicator } from "react-native";
import { useTheme } from "@/theme/ThemeProvider";

export function LoadingState() {
  const theme = useTheme();

  return (
    <SafeAreaView
      style={{ 
        flex: 1, 
        backgroundColor: theme.colors.background,
        justifyContent: 'center',
        alignItems: 'center'
      }}
    >
      <ActivityIndicator size="large" color={theme.colors.primaryGreen} />
      <Text style={{ 
        color: theme.colors.textSecondary,
        marginTop: 16,
        fontSize: 16
      }}>
        Loading results...
      </Text>
    </SafeAreaView>
  );
}
