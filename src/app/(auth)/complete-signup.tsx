import { useState } from "react"
import { View, Text, KeyboardAvoidingView, Platform, Alert } from "react-native"
import { useRouter } from "expo-router"
import { useSafeAreaInsets } from "react-native-safe-area-context"
import { useTheme } from "@/theme/ThemeProvider"
import { BackButton } from "@/components/ui/BackButton"
import { CTAButton } from "@/components/ui/CTAButton"
import { InputBox } from "@/components/ui/InputBox"
import { useSignUp } from "@clerk/clerk-expo"

export default function CompleteSignupScreen() {
  const router = useRouter()
  const theme = useTheme()
  const insets = useSafeAreaInsets()
  const { signUp, setActive, isLoaded } = useSignUp()

  const [phoneNumber, setPhoneNumber] = useState("")
  const [password, setPassword] = useState("")
  const [loading, setLoading] = useState(false)

  const handleComplete = async () => {
    if (!isLoaded || !signUp) return
    if (!password.trim()) {
      Alert.alert("Error", "Password is required")
      return
    }

    try {
      setLoading(true)

      // Update the signup with missing fields
      const updatedSignUp = await signUp.update({
        password,
        // phoneNumber: phoneNumber || undefined, // Only include if provided
      })

      console.log("Updated signup status:", updatedSignUp.status)

      if (updatedSignUp.status === 'complete') {
        await setActive({ session: updatedSignUp.createdSessionId })
        router.replace('/(auth)/create-pin')
      } else {
        console.log("Still missing requirements:", updatedSignUp.missingFields)
        Alert.alert("Error", "Unable to complete signup. Please check all required fields.")
      }

    } catch (err) {
      console.error("Complete signup error:", err)
      const errorMessage = err?.errors?.[0]?.longMessage || "Failed to complete signup"
      Alert.alert("Error", errorMessage)
    } finally {
      setLoading(false)
    }
  }

  return (
    <KeyboardAvoidingView 
      style={{ flex: 1 }} 
      behavior={Platform.OS === "ios" ? "padding" : "height"}
    >
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

        <View style={{ flex: 1, justifyContent: "center" }}>
          <View className="items-start mb-8">
            <Text 
              style={{
                fontSize: 28,
                fontFamily: theme.fonts.welcomeHeading,
                color: theme.colors.textPrimary,
                marginBottom: 8,
              }}
            >
              Complete Your Account
            </Text>
            <Text
              style={{
                fontSize: 16,
                color: theme.colors.textSecondary,
              }}
            >
              Please provide the following information to complete your account setup.
            </Text>
          </View>

          <View style={{ marginBottom: 20 }}>
            <InputBox
              // label="Password"
              value={password}
              onChangeText={setPassword}
              placeholder="Enter your password"
              secureTextEntry
            />
          </View>

          <View style={{ marginBottom: 40 }}>
            <InputBox
              // label="Phone Number (Optional)"
              value={phoneNumber}
              onChangeText={setPhoneNumber}
              placeholder="Enter your phone number"
              keyboardType="phone-pad"
            />
          </View>

          <CTAButton
            title="Complete Account"
            onPress={handleComplete}
            loading={loading}
            disabled={!password.trim() || loading}
          />
        </View>
      </View>
    </KeyboardAvoidingView>
  )
}
