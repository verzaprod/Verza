import React from "react";
import { View, Text, SafeAreaView, ScrollView } from "react-native";
import { useTheme } from "@/theme/ThemeProvider";
import { useSafeAreaInsets } from "react-native-safe-area-context";
import { VerifiersHeader } from "@/components/verifiers/VerifiersHeader";
import { VerifiersList } from "@/components/verifiers/VerifiersList";
import { SearchBar } from "@/components/verifiers/SearchBar";

const verifiers = [
  {
    id: "1",
    name: "TechCorp\nSolutions",
    type: "Enterprise",
    rating: 4.8,
    verified: 1240,
    logo: "shield-check",
    status: "active" as const,
    description: "Leading technology verification provider",
  },
  {
    id: "2", 
    name: "SecureID Pro",
    type: "Financial",
    rating: 4.9,
    verified: 856,
    logo: "shield",
    status: "active" as const,
    description: "Banking and financial services verification",
  },
  {
    id: "3",
    name: "VerifyNow",
    type: "General",
    rating: 4.6,
    verified: 2340,
    logo: "shield-check",
    status: "active" as const, 
    description: "Fast and reliable identity verification",
  },
  {
    id: "4",
    name: "TrustGuard",
    type: "Healthcare",
    rating: 4.7,
    verified: 445,
    logo: "shield",
    status: "busy" as const,
    description: "Healthcare industry specialist verification",
  },
];

export default function VerifiersScreen() {
  const theme = useTheme();
  const insets = useSafeAreaInsets();

  return (
    <SafeAreaView
      style={{
        flex: 1,
        backgroundColor: theme.colors.background,
        paddingTop: insets.top + 16,
      }}
    >
      <ScrollView
        className="flex-1"
        style={{ paddingHorizontal: 20 }}
        showsVerticalScrollIndicator={false}
      >

        <View style={{ marginBottom: theme.spacing.lg }}>
          <Text
            style={{
              fontSize: 24,
              color: theme.colors.textPrimary,
              fontFamily: theme.fonts.welcomeHeading,
              marginBottom: theme.spacing.sm,
            }}
          >
            Identity Verifiers
          </Text>
          <Text
            style={{
              fontSize: 16,
              color: theme.colors.textSecondary,
              marginBottom: theme.spacing.lg,
            }}
          >
            Choose a trusted verifier to authenticate your identity
          </Text>
          <SearchBar />
        </View>

        <View style={{ paddingBottom: theme.spacing.xl }}>
          <VerifiersList verifiers={verifiers} />
        </View>
      </ScrollView>
    </SafeAreaView>
  );
}