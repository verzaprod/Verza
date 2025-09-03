import { useState } from "react"
import { View, Text, KeyboardAvoidingView, Platform } from "react-native"
import { useRouter } from "expo-router"
import { useSafeAreaInsets } from "react-native-safe-area-context"
import { useTheme } from "@/theme/ThemeProvider"
import { BackButton } from "@/components/ui/BackButton"
import { CTAButton } from "@/components/ui/CTAButton"
import { InputBoxes } from "@/components/ui/InputBoxes"

export default function CreatePinScreen() {
  const router = useRouter()
  const theme = useTheme()
  const insets = useSafeAreaInsets()
  const [pin, setPin] = useState("")
  const [loading, setLoading] = useState(false)

  const handleCreatePin = async () => {
    if (pin.length !== 4) return

    setLoading(true)
    // TODO: Implement API call to create PIN
    // await authAPI.createPin(pin);

    setTimeout(() => {
      setLoading(false)
      router.push("/backup-passphrase")
    }, 2000)
  }

  return (
    <KeyboardAvoidingView style={{ flex: 1 }} behavior={Platform.OS === "ios" ? "padding" : "height"}>
      <View
        style={{
          flex: 1,
          backgroundColor: theme.isDark ? theme.colors.backgroundDark : theme.colors.backgroundLight,
          paddingTop: insets.top,
          paddingBottom: insets.bottom,
          paddingHorizontal: 20,
        }}
      >
        <View style={{ paddingVertical: 16 }}>
          <BackButton />
        </View>

        <View style={{ flex: 1, justifyContent: "center" }}>
          <View style={{ alignItems: "flex-start", marginBottom: 40 }}>
            <Text
              style={{
                fontSize: 28,
                fontWeight: "bold",
                color: theme.isDark ? theme.colors.textPrimaryDark : theme.colors.textPrimaryLight,
                marginBottom: 8,
              }}
            >
              Create Your PIN
            </Text>
            <Text
              style={{
                fontSize: 16,
                color: theme.colors.textSecondary,
              }}
            >
              Create a 4-digit pin to secure your wallet
            </Text>
          </View>

          <View style={{ marginBottom: 40 }}>
            <InputBoxes value={pin} onChangeText={setPin} length={4} type="pin" />
          </View>

          <CTAButton title="Create PIN" onPress={handleCreatePin} loading={loading} disabled={pin.length !== 4} />
        </View>
      </View>
    </KeyboardAvoidingView>
  )
}
