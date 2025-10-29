import React, { useState } from "react";
import {
  View,
  Text,
  SafeAreaView,
  ScrollView,
  TouchableOpacity,
} from "react-native";
import { useRouter } from "expo-router";
import { useTheme } from "@/theme/ThemeProvider";
import { useSafeAreaInsets } from "react-native-safe-area-context";
import { ProfileHeader } from "@/components/profile/ProfileHeader";
import { CredentialsList } from "@/components/profile/CredentialsList";
import { AddCredentialButton } from "@/components/profile/AddCredentialButton";
import { AddCredentialModal } from "@/components/profile/AddCredentialModal";
import { useKYCStore } from "@/store/kycStore";
import { useClerk } from "@clerk/clerk-expo";
import { useAuthStore } from "@/store/authStore";

export default function ProfileScreen() {
  const theme = useTheme();
  const insets = useSafeAreaInsets();

  const [credentials, setCredentials] = useState([]);
  const [modalVisible, setModalVisible] = useState(false);
  const { signOut } = useClerk();

  const verificationStatus = useKYCStore((state) => state.verificationStatus);
  const setAuthenticated = useAuthStore((state) => state.setAuthenticated);

  const router = useRouter();

  const handleAddCredential = (credentialType) => {
    const newCredential = {
      id: Date.now().toString(),
      type: credentialType.label,
      status: "pending",
      icon: credentialType.icon,
    };
    setCredentials((prev) => [...prev, newCredential]);
  };

  const handleSignOut = async () => {
    try {
      await signOut();
      setAuthenticated(false);
      router.replace("/(auth)/sign-in");
    } catch (err) {
      console.error(JSON.stringify(err, null, 2));
    }
  };

  console.log("VerificationStatus", verificationStatus);

  const handleRemoveCredential = (credentialId) => {
    setCredentials((prev) => prev.filter((cred) => cred.id !== credentialId));
  };

  const handleAddCredentialPress = () => {
    setModalVisible(true);
  };

  return (
    <SafeAreaView
      style={{
        flex: 1,
        backgroundColor: theme.colors.background,
        paddingTop: insets.top + 24,
      }}
    >
      <ScrollView
        className="flex-1"
        style={{ paddingHorizontal: 20 }}
        showsVerticalScrollIndicator={false}
      >
        <ProfileHeader />

        {credentials.length > 0 && (
          <View style={{ marginBottom: theme.spacing.xl }}>
            <Text
              style={{
                fontSize: 24,
                fontWeight: "bold",
                color: theme.colors.textPrimary,
                fontFamily: theme.fonts.welcomeHeading,
                marginBottom: theme.spacing.lg,
                textAlign: "center",
              }}
            >
              Enlisted Credentials
            </Text>
            <CredentialsList
              credentials={credentials}
              onRemoveCredential={handleRemoveCredential}
            />
          </View>
        )}

        {credentials.length === 0 && (
          <View
            style={{
              flex: 1,
              justifyContent: "center",
              alignItems: "center",
              marginTop: theme.spacing.xl * 2,
            }}
          >
            <Text
              style={{
                fontSize: 18,
                color: theme.colors.textSecondary,
                textAlign: "center",
                marginBottom: theme.spacing.lg,
              }}
            >
              No credentials added yet
            </Text>
            <Text
              style={{
                fontSize: 14,
                color: theme.colors.textSecondary,
                textAlign: "center",
                marginBottom: theme.spacing.xl,
              }}
            >
              Add your first credential to get started
            </Text>
          </View>
        )}

        <View style={{ alignItems: "center", paddingBottom: theme.spacing.xl }}>
          <AddCredentialButton onPress={handleAddCredentialPress} />
        </View>

        { (
          <TouchableOpacity
            onPress={() => router.replace("/verifier")}
            style={{
              backgroundColor: theme.colors.background,
              paddingVertical: 14,
              borderRadius: 12,
              marginBottom: 16,
            }}
            className="w-full items-center justify-center py-4 rounded-xl mb -4"
          >
            <Text
              style={{
                fontSize: 14,
                color: theme.colors.textPrimary,
                textAlign: "center",
                marginBottom: theme.spacing.xl,
              }}
            >
              Become a Verifier
            </Text>
          </TouchableOpacity>
        )}

        
        <TouchableOpacity
          onPress={handleSignOut}
          className="w-full items-center justify-center py-4 rounded-xl mb-4"
          style={{
            backgroundColor: theme.colors.boxBorder,
            borderRadius: 12,
          }}
        >
          <Text style={{ color: theme.colors.error}}>Sign out</Text>
        </TouchableOpacity>
      </ScrollView>

      <AddCredentialModal
        visible={modalVisible}
        onClose={() => setModalVisible(false)}
        onSelect={handleAddCredential}
      />
    </SafeAreaView>
  );
}
