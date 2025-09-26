import React from "react";
import { View, Text, SafeAreaView, ScrollView } from "react-native";
import { useTheme } from "@/theme/ThemeProvider";
import { useSafeAreaInsets } from "react-native-safe-area-context";
import { ProfileHeader } from "@/components/profile/ProfileHeader";
import { CredentialsList } from "@/components/profile/CredentialsList";
import { AddCredentialButton } from "@/components/profile/AddCredentialButton";

const credentials = [
  { id: "1", type: "ID Card", status: "verified", icon: "id-card" },
  { id: "2", type: "Passport", status: "verified", icon: "passport" },
  { id: "3", type: "Proof of Address", status: "pending", icon: "home" },
];

export default function ProfileScreen() {
  const theme = useTheme();
  const insets = useSafeAreaInsets();

  return (
    <SafeAreaView
      style={{
        flex: 1,
        backgroundColor: theme.colors.background,
        paddingTop: insets.top,
      }}
    >
      <ScrollView
        className="flex-1"
        style={{ paddingHorizontal: 20 }}
        showsVerticalScrollIndicator={false}
      >
        <ProfileHeader />

        <View style={{ marginBottom: theme.spacing.xl }}>
          <Text
            style={{
              fontSize: 24,
              fontWeight: "bold",
              color: theme.colors.textPrimary,
              fontFamily: theme.fonts.welcomeHeading,
              marginBottom: theme.spacing.lg,
            }}
          >
            Enlisted Credentials
          </Text>
          <CredentialsList credentials={credentials} />
        </View>

        <View style={{ alignItems: "center", paddingBottom: theme.spacing.xl }}>
          <AddCredentialButton />
        </View>
      </ScrollView>
    </SafeAreaView>
  );
}
