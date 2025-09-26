import React, { useState } from 'react'
import { View, Text, SafeAreaView, Alert, Image } from 'react-native'
import { useRouter } from 'expo-router'
import { useTheme } from '@/theme/ThemeProvider'
import { useSafeAreaInsets } from 'react-native-safe-area-context'
import { Button } from '@/components/ui/Button'
import { CameraCapture } from '@/components/kyc/CameraCapture'
import { useKYCStore } from '@/store/kycStore'
import * as ImagePicker from 'expo-image-picker'

export default function DocCapture() {
  const theme = useTheme()
  const router = useRouter()
  const insets = useSafeAreaInsets()
  const [frontImage, setFrontImage] = useState<string | null>(null)
  const [backImage, setBackImage] = useState<string | null>(null)
  const [isUploading, setIsUploading] = useState(false)

  const { selectedDocType, setDocumentImages, setCurrentStep } = useKYCStore()

  const docInfo = {
    'passport': { title: 'Passport', requiresBack: false },
    'driver-license': { title: "Driver's License", requiresBack: true },
    'id-card': { title: 'National ID Card', requiresBack: true }
  }

  const currentDoc = docInfo[selectedDocType || 'passport']
  const needsBothSides = currentDoc.requiresBack

  const handleImageCapture = async (side: 'front' | 'back') => {
    try {
      const result = await ImagePicker.launchCameraAsync({
        mediaTypes: ImagePicker.MediaTypeOptions.Images,
        allowsEditing: true,
        aspect: [16, 10],
        quality: 0.8,
      })

      if (!result.canceled && result.assets[0]) {
        const imageUri = result.assets[0].uri
        if (side === 'front') {
          setFrontImage(imageUri)
        } else {
          setBackImage(imageUri)
        }
      }
    } catch (error) {
      Alert.alert('Error', 'Failed to capture image. Please try again.')
    }
  }

  const handleContinue = async () => {
    if (!frontImage || (needsBothSides && !backImage)) {
      Alert.alert('Missing Images', 'Please capture all required document photos.')
      return
    }

    setIsUploading(true)
    
    try {
      // Simulate upload delay
      await new Promise(resolve => setTimeout(resolve, 2000))
      
      const images = needsBothSides ? [frontImage, backImage!] : [frontImage]
      setDocumentImages(images)
      setCurrentStep('selfie')
      
      router.push('/(kyc)/selfie-note')
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
        paddingTop: insets.top,
      }}
    >
      <View style={{ flex: 1, paddingHorizontal: 20, paddingTop: 20 }}>
        <View style={{ marginBottom: 32 }}>
          <Text
            style={{
              fontSize: 24,
              fontWeight: 'bold',
              color: theme.colors.textPrimary,
              fontFamily: theme.fonts.welcomeHeading,
              marginBottom: 8,
            }}
          >
            Document Photos
          </Text>
          <Text
            style={{
              fontSize: 16,
              color: theme.colors.textSecondary,
              lineHeight: 24,
            }}
          >
            Take clear photos of your {currentDoc.title.toLowerCase()}
          </Text>
        </View>

        <View style={{ flex: 1, gap: 20 }}>
          {/* Front Side */}
          <CameraCapture
            title="Front Side"
            subtitle="Capture the front of your document"
            image={frontImage}
            onCapture={() => handleImageCapture('front')}
          />

          {/* Back Side (if required) */}
          {needsBothSides && (
            <CameraCapture
              title="Back Side"
              subtitle="Capture the back of your document"
              image={backImage}
              onCapture={() => handleImageCapture('back')}
            />
          )}
        </View>

        <View style={{ paddingBottom: 20 }}>
          <Button
            text={isUploading ? "Uploading..." : "Continue"}
            onPress={handleContinue}
            disabled={!frontImage || (needsBothSides && !backImage) || isUploading}
          />
        </View>
      </View>
    </SafeAreaView>
  )
}
