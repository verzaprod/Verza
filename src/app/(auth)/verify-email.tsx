import { useEffect, useState } from "react"
import { View, Text, TouchableOpacity, KeyboardAvoidingView, Platform, Alert } from "react-native"
import { useRouter } from "expo-router"
import { useSafeAreaInsets } from "react-native-safe-area-context"
import { useTheme } from "@/theme/ThemeProvider"
import { BackButton } from "@/components/ui/BackButton"
import { CTAButton } from "@/components/ui/CTAButton"
import { InputBoxes } from "@/components/ui/InputBoxes"
import { useSignUp } from "@clerk/clerk-expo"

export default function VerifyEmailScreen() {
  const router = useRouter()
  const theme = useTheme()
  const insets = useSafeAreaInsets()
  const { signUp, setActive, isLoaded } = useSignUp()

  const [otp, setOtp] = useState("")
  const [loading, setLoading] = useState(false)
  const [newAttempt, setNewAttempt] = useState(null)
  
  const onVerifyPress = async () => {
    if (otp.length !== 6) return
    if (!isLoaded || !signUp) return

    try {
      setLoading(true)

      const signUpAttempt = await signUp.attemptEmailAddressVerification({
        code: otp,
      })

      console.log("Verification attempt status:", signUpAttempt.status)

      if (signUpAttempt.status === 'complete') {
        await setActive({ session: signUpAttempt.createdSessionId })
        router.replace('/(auth)/create-pin')
        
      } else if (signUpAttempt.status === 'missing_requirements') {
        console.log("Missing requirements:", signUpAttempt.missingFields)
        
        // Check what's missing and handle accordingly
        const missingFields = signUpAttempt.missingFields || []
        
        if (missingFields.includes('password') || missingFields.includes('phone_number')) {
          // Navigate to a screen to collect missing information
          router.push('/(auth)/complete-signup')
        } else {
          Alert.alert("Verification Error", "Additional information required to complete verification.")
        }
      }

    } catch (err) {
      console.error("Verification error:", err)
      
      // Handle the "already verified" error specifically
      if (err.message?.includes('already been verified')) {
        console.log("Already verified, checking signup status...")
        
        // Check if we can complete with missing requirements
        if (signUp.status === 'missing_requirements') {
          router.push('/(auth)/complete-signup')
        } else {
          router.replace('/(auth)/create-pin')
        }
      } else if (err?.errors && err.errors.length > 0) {
        const errorMessage = err.errors[0]?.longMessage || err.errors[0]?.message || "Verification failed"
        Alert.alert("Verification Failed", errorMessage)
      } else {
        Alert.alert("Verification Failed", "An unexpected error occurred. Please try again.")
      }
    } finally {
      setLoading(false)
    }
  }

  // Check if user is already verified on component mount
  useEffect(() => {
    if (isLoaded && signUp) {
      console.log("Current signUp status:", signUp.status)
      
      // If already complete, redirect immediately
      if (signUp.status === 'complete') {
        console.log("SignUp already complete, redirecting...")
        router.replace('/(auth)/create-pin')
      }
    }
  }, [isLoaded, signUp, router])

  const handleResend = async () => {
    if (!isLoaded || !signUp) return

    try {
      setLoading(true)
      
      // Resend the verification email
      await signUp.prepareEmailAddressVerification()
      
      Alert.alert("Code Resent", "A new verification code has been sent to your email.")
      
    } catch (err) {
      console.error("Resend error:", err)
      Alert.alert("Resend Failed", "Unable to resend verification code. Please try again.")
    } finally {
      setLoading(false)
    }
  }

  return (
    <KeyboardAvoidingView style={{ flex: 1 }} behavior={Platform.OS === "ios" ? "padding" : "height"}>
      <View
        style={{
          flex: 1,
          backgroundColor: theme.colors.background,
          paddingTop: insets.top,
          paddingHorizontal: 20,
        }}
      >
        <View style={{ paddingVertical: 16 }}>
          <BackButton />
        </View>

        <View>
          <View className="items-start mb-20">
            <Text className="mb-2"
              style={{
                fontSize: 28,
                fontFamily: theme.fonts.welcomeHeading,
                color: theme.colors.textPrimary,
              }}
            >
              Verify Your Email
            </Text>
            <Text
              style={{
                fontSize: 16,
                color: theme.colors.textSecondary,
              }}
            >
              We have sent the verification code to your email address.
            </Text>
          </View>

          <View style={{ marginBottom: 40 }}>
            <InputBoxes value={otp} onChangeText={setOtp} length={6} type="otp" />
          </View>

          <View style={{ marginBottom: 20 }}>
            <CTAButton title="Verify Code" onPress={onVerifyPress} loading={loading} disabled={otp.length !== 6} />
          </View>

          <View style={{ alignItems: "center" }}>
            <Text style={{ color: theme.colors.textSecondary }}>
              Didn't receive code?{" "}
                <Text style={{ color: theme.colors.primaryGreen }} onPress={handleResend}>
                  Resend
                </Text>
            </Text>
          </View>
        </View>
      </View>
    </KeyboardAvoidingView>
  )
}
