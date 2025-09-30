import React, { useState } from 'react'
import { View, Text, SafeAreaView, Alert } from 'react-native'
import { useRouter } from 'expo-router'
import { useTheme } from '@/theme/ThemeProvider'
import { useSafeAreaInsets } from 'react-native-safe-area-context'
import { Button } from '@/components/ui/Button'
import { CameraCapture } from '@/components/kyc/CameraCapture'
import { useKYCStore } from '@/store/kycStore'
import * as ImagePicker from 'expo-image-picker'

export default function SelfieCapture() {
  const theme = useTheme()
  const router = useRouter()
  const insets = useSafeAreaInsets()
  const [selfieImage, setSelfieImageLocal] = useState<string | null>(null)
  const [isUploading, setIsUploading] = useState(false)

  const setVerificationStatus = useKYCStore((state) => state.setVerificationStatus);
  const setCurrentStep = useKYCStore((state) => state.setCurrentStep);

  const handleSelfieCapture = async () => {
    try {
      const result = await ImagePicker.launchCameraAsync({
        mediaTypes: ImagePicker.MediaTypeOptions.Images,
        allowsEditing: true,
        aspect: [3, 4],
        quality: 0.8,
        cameraType: ImagePicker.CameraType.front, 
      })

      if (!result.canceled && result.assets[0]) {
        setSelfieImageLocal(result.assets[0].uri)
      }
    } catch (error) {
      Alert.alert('Error', 'Failed to capture selfie. Please try again.')
    }
  }

  const handleContinue = async () => {
    if (!selfieImage) {
      Alert.alert('Missing Selfie', 'Please take a selfie to continue.')
      return
    }

    setIsUploading(true)
    
    try {
      // Simulate upload delay
      await new Promise(resolve => setTimeout(resolve, 2000))
      
      // setSelfieImage(selfieImage)
      setCurrentStep('processing')
      setVerificationStatus('verified')
      
      router.replace('/(kyc)/verification-tracker')
    } catch (error) {
      Alert.alert('Upload Failed', 'Please try again.')
    } finally {
      setIsUploading(false)
    }
  }

  return (
    <SafeAreaView
      style={{
        flex: 1,
        backgroundColor: theme.colors.background,
        paddingTop: insets.top + 24,
      }}
    >
      <View style={{ flex: 1, paddingHorizontal: 20, paddingTop: 0 }}>
        <View style={{ marginBottom: 32 }}>
          <Text
            style={{
              fontSize: 24,
              color: theme.colors.textPrimary,
              fontFamily: theme.fonts.welcomeHeading,
              marginBottom: 8,
            }}
          >
            Take a Selfie
          </Text>
          <Text
            style={{
              fontSize: 16,
              color: theme.colors.textSecondary,
              lineHeight: 24,
            }}
          >
            Take a clear photo of yourself for verification
          </Text>
        </View>

        <View style={{ flex: 1,  }}>
          <CameraCapture
            title=""
            subtitle="Look directly at the camera"
            image={selfieImage}
            onCapture={handleSelfieCapture}
            isCircular={true}
          />
        </View>

        <View style={{ paddingBottom: 20 }}>
          <Button
            text="Complete Verification"
            onPress={handleContinue}
            disabled={!selfieImage || isUploading}
          />
        </View>
      </View>
    </SafeAreaView>
  )
}
