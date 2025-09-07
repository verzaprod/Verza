import { useState } from "react"
import { View, Text, TouchableOpacity, KeyboardAvoidingView, Platform } from "react-native"
import { useRouter } from "expo-router"
import { useSafeAreaInsets } from "react-native-safe-area-context"
import { useTheme } from "@/theme/ThemeProvider"
import { BackButton } from "@/components/ui/BackButton"
import { CTAButton } from "@/components/ui/CTAButton"
import { InputBoxes } from "@/components/ui/InputBoxes"

export default function VerifyEmailScreen() {
  const router = useRouter()
  const theme = useTheme()
  const insets = useSafeAreaInsets()
  const [otp, setOtp] = useState("")
  const [loading, setLoading] = useState(false)

  const handleVerify = async () => {
    if (otp.length !== 4) return

    setLoading(true)
    // TODO: Implement API call to verify OTP
    // await authAPI.verifyOTP(otp);

    setTimeout(() => {
      setLoading(false)
      router.push("/(auth)/create-pin")
    }, 2000)
  }

  const handleResend = async () => {
    // TODO: Implement resend OTP API call
    // await authAPI.resendOTP();
  }

  return (
    <KeyboardAvoidingView style={{ flex: 1 }} behavior={Platform.OS === "ios" ? "padding" : "height"}>
      <View
        style={{
          flex: 1,
          backgroundColor: theme.colors.background,
          paddingTop: insets.top,
          paddingBottom: insets.bottom,
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
            <InputBoxes value={otp} onChangeText={setOtp} length={4} type="otp" />
          </View>

          <View style={{ marginBottom: 20 }}>
            <CTAButton title="Verify Code" onPress={handleVerify} loading={loading} disabled={otp.length !== 4} />
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
