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
import { useSafeAreaInsets, } from "react-native-safe-area-context";

export default function RegisterScreen() {
  const theme = useTheme();
  const router = useRouter();
const insets = useSafeAreaInsets();

  const { isLoaded, signUp, setActive } = useSignUp();
  const [emailAddress, setEmailAddress] = useState("");
  const [password, setPassword] = useState("");
  const [loading, setLoading] = useState(false);

  const onSignUpPress = async () => {
    if (!emailAddress.trim()) return;

    if (!isLoaded) return;

    try {
      setLoading(true);
      await signUp.create({
        emailAddress,
        password,
      });

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
        paddingTop: 0,
        paddingBottom: insets.bottom,
      }}
    >
      <KeyboardAvoidingView
        className="flex-1"
        behavior={Platform.OS === "ios" ? "padding" : "height"}
      >
        <View className="flex-1 align-center justify-center">
          <View className="px-6">
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
              Enter your email to get started
            </Text>
          </View>

          <View className="items-center justify-center mb-10">
            <Icon
              name="welcome"
              style={{
                width: WIDTH * 0.8,
                height: HEIGHT * 0.8,
              }}
            />
          </View>

          <View className="space-y-6 px-6">
            <InputBox
              placeholder="Email e.g. michealjackson@gmail.com"
              value={emailAddress}
              onChangeText={setEmailAddress}
              keyboardType="email-address"
              autoCapitalize="none"
              returnKeyType="next"
            />

            <View className="mb-4" />

            <InputBox
              placeholder="Password e.g. X8df!90EO"
              value={password}
              onChangeText={setPassword}
              keyboardType="default"
              secureTextEntry
              returnKeyType="done"
            />

            <View className="mb-4" />

            <CTAButton
              title="Continue"
              onPress={onSignUpPress}
              loading={loading}
              disabled={!emailAddress.trim() || !password.trim()}
            />
          </View>

          <View className="flex-row justify-center items-center mt-4">
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
          
        </View>

      </KeyboardAvoidingView>
    </SafeAreaView>
  );
}
