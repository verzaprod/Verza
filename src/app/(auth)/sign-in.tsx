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
import { useAuth, useSignIn, useUser } from "@clerk/clerk-expo";
import { useSafeAreaInsets } from "react-native-safe-area-context";
import { useAuthStore } from "@/store/authStore";

export default function SignInScreen() {
  const theme = useTheme();
  const router = useRouter();
  const insets = useSafeAreaInsets();
  const { signIn, setActive, isLoaded } = useSignIn();
  const { isSignedIn } = useAuth();
  const [emailAddress, setEmailAddress] = useState("");
  const [password, setPassword] = useState("");
  const [loading, setLoading] = useState(false);

  const { user } = useUser();

  const {
    pinCreated,
    passphraseBackedUp,
    setAuthenticated,
    setFirstTimeUser,
    isFirstTimeUser,
    isAuthenticated,
  } = useAuthStore();

  console.log("Firsttimeuser", isFirstTimeUser, user);

  const onSignInPress = async () => {
    if (!emailAddress.trim() || !password.trim()) {
      Alert.alert("Error", "Please fill in all fields");
      return;
    }

    if (!isLoaded) return;

    try {
      setLoading(true);

      const signInAttempt = await signIn.create({
        identifier: emailAddress,
        password,
      });

      if (signInAttempt.status === "complete") {
        await setActive({ session: signInAttempt.createdSessionId });

        router.replace("/(tabs)/home");
        // const hasCompletedOnboarding = pinCreated && passphraseBackedUp;

        // if (hasCompletedOnboarding) {
        //   setFirstTimeUser(false);
        //   setAuthenticated(true);
        //   router.replace("/(tabs)/home");
        // } else {
        //   setFirstTimeUser(true);
        //   if (!pinCreated) {
        //     router.replace("/(auth)/create-pin");
        //   } else if (!passphraseBackedUp) {
        //     router.replace("/(auth)/backup-passphrase");
        //   } else {
        //     router.replace("/(tabs)/home");
        //   }
        // }
      } else {
        console.log("Sign-in incomplete:", signInAttempt.status);
        Alert.alert(
          "Sign In Failed",
          "Unable to complete sign in. Please try again."
        );
      }
    } catch (err) {
      console.error("Sign-in error:", err);
      const errorMessage =
        err?.errors?.[0]?.longMessage ||
        "Sign in failed. Please check your credentials.";
      Alert.alert("Sign In Failed", errorMessage);
    } finally {
      setLoading(false);
    }
  };

  const handleForgotPassword = () => {
    // TODO: Implement forgot password flow
    Alert.alert(
      "Forgot Password",
      "Password reset functionality will be implemented soon."
    );
  };

  const handleCreateAccount = () => {
    router.push("/(auth)/register");
  };

  return (
    <SafeAreaView
      className="flex-1"
      style={{
        backgroundColor: theme.colors.background,
        paddingBottom: insets.bottom,
      }}
    >
      <KeyboardAvoidingView
        className="flex-1"
        behavior={Platform.OS === "ios" ? "padding" : "height"}
      >
        <View className="flex-1 justify-center">
          <View className="px-6">
            <Text
              className="text-3xl text-center font-bold mb-2"
              style={{
                color: theme.colors.textPrimary,
                fontFamily: theme.fonts.welcomeHeading,
              }}
            >
              Welcome Back
            </Text>
            <Text
              className="text-lg text-center"
              style={{
                color: theme.colors.textSecondary,
                fontFamily: theme.fonts.onboardingTagline,
              }}
            >
              Sign in to your Verza account
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
              placeholder="Email e.g. michaeljackson@gmail.com"
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
              secureTextEntry
              autoCapitalize="none"
              returnKeyType="done"
            />

            <View className="mb-4" />

            <CTAButton
              title="Sign In"
              onPress={onSignInPress}
              loading={loading}
              disabled={!emailAddress.trim() || !password.trim()}
            />

            {user && (
              <View className="flex-row justify-center items-center mt-4">
                <Text
                  style={{
                    color: theme.colors.textSecondary,
                    fontSize: 16,
                  }}
                >
                  Don't have an account?{" "}
                </Text>
                <TouchableOpacity onPress={handleCreateAccount}>
                  <Text
                    style={{
                      color: theme.colors.primaryGreen,
                      fontSize: 16,
                      fontWeight: "600",
                    }}
                  >
                    Create Account
                  </Text>
                </TouchableOpacity>
              </View>
            )}
          </View>
        </View>
      </KeyboardAvoidingView>
    </SafeAreaView>
  );
}
