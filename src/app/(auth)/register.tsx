import React, { useState } from 'react';
import { View, Text, KeyboardAvoidingView, Platform } from 'react-native';
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
    try {
      const result = await apiClient.sendVerificationCode(emailOrPhone);
      if (result.success) {
        setEmail(emailOrPhone);
        router.push('/(auth)/verify-email');
      }
    } catch (error) {
      // Handle error
    } finally {
      setLoading(false);
    }
  };

  return (
    <SafeAreaView className="flex-1" style={{ backgroundColor: theme.colors.background }}>
      <KeyboardAvoidingView 
        className="flex-1" 
        behavior={Platform.OS === 'ios' ? 'padding' : 'height'}
      >
        <View className="flex-1 px-6 align-center justify-center">
          <View className="mb-8">
            <Text 
              className="text-3xl text-center font-bold mb-2"
              style={{ 
                color: theme.colors.textPrimary 
              }}
            >
              Welcome to Verza
            </Text>
            <Text 
              className="text-lg text-center"
              style={{ color: theme.colors.textSecondary }}
            >
              Enter your email or phone number to get started
            </Text>
          </View>

          <View className="mb-8">
            <Icon 
              name="welcome" 
              style={{
                width: WIDTH, 
                height: HEIGHT
              }}
            />
          </View>

          <View className="space-y-4">
            <InputBox
              placeholder="Enter email or phone number"
              value={emailOrPhone}
              onChangeText={setEmailOrPhone}
              keyboardType="email-address"
              autoCapitalize="none"
            />
            
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