import React, { useState } from 'react';
import { View, Text, KeyboardAvoidingView, Platform, Alert, TouchableOpacity } from 'react-native';
import { SafeAreaView } from 'react-native-safe-area-context';
import { useRouter } from 'expo-router';
import { useTheme } from '@/theme/ThemeProvider';
import { CTAButton } from '@/components/ui/CTAButton';
import { InputBox } from '@/components/ui/InputBox';
import { Icon } from '@/components/ui/Icon';
import { WIDTH, HEIGHT } from '@/constants';
import { useSignIn } from '@clerk/clerk-expo';

export default function SignInScreen() {
  const theme = useTheme();
  const router = useRouter();
  const { signIn, setActive, isLoaded } = useSignIn();
  const [emailOrPhone, setEmailOrPhone] = useState('');
  const [password, setPassword] = useState('');
  const [loading, setLoading] = useState(false);

  const onSignInPress = async () => {
    if (!emailOrPhone.trim() || !password.trim()) {
      Alert.alert("Error", "Please fill in all fields");
      return;
    }

    if (!isLoaded) return;

    try {
      setLoading(true);
      
      const signInAttempt = await signIn.create({
        identifier: emailOrPhone,
        // password,
      });

      if (signInAttempt.status === 'complete') {
        await setActive({ session: signInAttempt.createdSessionId });
        router.replace("/(tabs)/home");
      } else {
        console.log("Sign-in incomplete:", signInAttempt.status);
        Alert.alert("Sign In Failed", "Unable to complete sign in. Please try again.");
      }
    } catch (err) {
      console.error("Sign-in error:", err);
      const errorMessage = err?.errors?.[0]?.longMessage || "Sign in failed. Please check your credentials.";
      Alert.alert("Sign In Failed", errorMessage);
    } finally {
      setLoading(false);
    }
  };

  const handleForgotPassword = () => {
    // TODO: Implement forgot password flow
    Alert.alert("Forgot Password", "Password reset functionality will be implemented soon.");
  };

  const handleCreateAccount = () => {
    router.push("/(auth)/register");
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
        behavior={Platform.OS === 'ios' ? 'padding' : 'height'}
      >
        <View className="flex-1 justify-center">
          {/* Header Section */}
          <View className="mb-8 px-6">
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

          {/* Illustration */}
          <View className="items-center justify-center mb-8">
            <Icon 
              name="welcome" 
              style={{
                width: WIDTH * 0.8, 
                height: HEIGHT * 0.8,
              }}
            />
          </View>

          {/* Form Section */}
          <View className="space-y-6 px-6">
            <InputBox
              placeholder="Enter email or phone number"
              value={emailOrPhone}
              onChangeText={setEmailOrPhone}
              keyboardType="email-address"
              autoCapitalize="none"
            />
            
            <InputBox
              placeholder="Enter your password"
              value={password}
              onChangeText={setPassword}
              secureTextEntry
              autoCapitalize="none"
            />

            {/* Forgot Password Link */}
            <TouchableOpacity 
              onPress={handleForgotPassword}
              style={{ alignSelf: 'flex-end' }}
            >
              <Text
                style={{
                  color: theme.colors.primaryGreen,
                  fontSize: 14,
                  fontWeight: '600',
                }}
              >
                Forgot Password?
              </Text>
            </TouchableOpacity>
            
            <View className="mb-4"/>

            <CTAButton
              title="Sign In"
              onPress={onSignInPress}
              loading={loading}
              disabled={!emailOrPhone.trim() || !password.trim()}
            />

            {/* Create Account Link */}
            <View className="flex-row justify-center items-center mt-6">
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
                    fontWeight: '600',
                  }}
                >
                  Create Account
                </Text>
              </TouchableOpacity>
            </View>
          </View>
        </View>
      </KeyboardAvoidingView>
    </SafeAreaView>
  );
}