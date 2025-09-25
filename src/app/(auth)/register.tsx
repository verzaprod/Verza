import React, { useState } from "react";
import {
  View,
  Text,
  KeyboardAvoidingView,
  Platform,
  Alert,
  TouchableOpacity,
} from "react-native";
import { SafeAreaView } from "react-native-safe-area-context";
import { useRouter } from "expo-router";
import { useTheme } from "@/theme/ThemeProvider";
import { CTAButton } from "@/components/ui/CTAButton";
import { InputBox } from "@/components/ui/InputBox";
import { Icon } from "@/components/ui/Icon";
import { WIDTH, HEIGHT } from "@/constants";
import { useAuth, useSignUp } from "@clerk/clerk-expo";

export default function RegisterScreen() {
  const theme = useTheme();
  const router = useRouter();
  const { isLoaded, signUp, setActive } = useSignUp();
  const [emailOrPhone, setEmailOrPhone] = useState("");
  const [loading, setLoading] = useState(false);

  const onSignUpPress = async () => {
    if (!emailOrPhone.trim()) return;

    if (!isLoaded) return;

    try {
      setLoading(true);
      await signUp.create({
        emailAddress: emailOrPhone,
      });

      // Send user an email with verification code
      await signUp.prepareEmailAddressVerification({ strategy: "email_code" });
      router.push("/(auth)/verify-email");
    } catch (err) {
      setLoading(false);
      console.error(JSON.stringify(err, null, 2));
      Alert.alert("Registration Failed", `${err.errors[0].longMessage}`);
    } finally {
      setLoading(false);
    }
  };

  return (
    <SafeAreaView
      className="flex-1"
      style={{
        backgroundColor: theme.colors.background,
      }}
    >
      <KeyboardAvoidingView
        className="flex-1"
        behavior={Platform.OS === "ios" ? "padding" : "height"}
      >
        <View className="flex-1 align-center justify-center">
          <View className="mb-8 px-6">
            <Text
              className="text-3xl text-center font-bold mb-2"
              style={{
                color: theme.colors.textPrimary,
                fontFamily: theme.fonts.welcomeHeading,
              }}
            >
              Welcome to Verza
            </Text>
            <Text
              className="text-lg text-center"
              style={{
                color: theme.colors.textSecondary,
                fontFamily: theme.fonts.onboardingTagline,
              }}
            >
              Enter your email or phone number to get started
            </Text>
          </View>

          <View className="items-center justify-center">
            <Icon
              name="welcome"
              style={{
                width: WIDTH,
                height: HEIGHT,
              }}
            />
          </View>

          <View className="space-y-6 px-6">
            <InputBox
              placeholder="Enter email or phone number"
              value={emailOrPhone}
              onChangeText={setEmailOrPhone}
              keyboardType="email-address"
              autoCapitalize="none"
            />

            <View className="mb-4" />

            <CTAButton
              title="Continue"
              onPress={onSignUpPress}
              loading={loading}
              disabled={!emailOrPhone.trim()}
            />
          </View>
        </View>

        <View className="flex-row justify-center items-center mt-6">
          <Text
            style={{
              color: theme.colors.textSecondary,
              fontSize: 16,
            }}
          >
            Already have an account?{" "}
          </Text>
          <TouchableOpacity onPress={() => router.push("/(auth)/sign-in")}>
            <Text
              style={{
                color: theme.colors.primaryGreen,
                fontSize: 16,
                fontWeight: "600",
              }}
            >
              Sign In
            </Text>
          </TouchableOpacity>
        </View>
      </KeyboardAvoidingView>
    </SafeAreaView>
  );
}

