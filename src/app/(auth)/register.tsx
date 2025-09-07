import React, { useState } from 'react';
import { View, Text, KeyboardAvoidingView, Platform, Alert } from 'react-native';
import { SafeAreaView } from 'react-native-safe-area-context';
import { router } from 'expo-router';
import { useTheme } from '@/theme/ThemeProvider';
import { CTAButton } from '@/components/ui/CTAButton';
import { InputBox } from '@/components/ui/InputBox';
import { Icon } from '@/components/ui/Icon';
import { useAuthStore } from '@/store/authStore';
import { apiClient } from '@/api/client';
import { WIDTH, HEIGHT } from '@/constants';

export default function RegisterScreen() {
  const theme = useTheme();
  const [emailOrPhone, setEmailOrPhone] = useState('');
  const [loading, setLoading] = useState(false);
  const { setEmail } = useAuthStore();

  const handleContinue = async () => {
    if (!emailOrPhone.trim()) return;
    
    setLoading(true);
    // setTimeout(() => { 
    //   setLoading(false);
    //   router.replace('/(auth)/verify-email');
    // }, 2000)

    try {
      const result = await apiClient.sendVerificationCode(emailOrPhone);
      // if (result.success) {
      //   setEmail(emailOrPhone);
      router.replace('/(auth)/verify-email');
      // }
    } catch (error) {
      console.error('Error sending verification code:', error);
      // TODO: show an error message to the user
      Alert.alert("Registration failed!", "Please try again.")
    } finally {
      setLoading(false);
    }
  };

  return (
    <SafeAreaView 
      className="flex-1" 
      style={{ 
        backgroundColor: theme.colors.background, 
    }}>
      <KeyboardAvoidingView 
        className="flex-1" 
        behavior={Platform.OS === 'ios' ? 'padding' : 'height'}
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
            
            <View className="mb-4"/>

            <CTAButton
              title="Continue"
              onPress={handleContinue}
              loading={loading}
              disabled={!emailOrPhone.trim()}
            />
          </View>
        </View>
      </KeyboardAvoidingView>
    </SafeAreaView>
  );
}