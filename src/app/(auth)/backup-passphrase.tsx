import { useState, useEffect } from "react"
import { View, Text, KeyboardAvoidingView, Platform, Alert } from "react-native"
import { useRouter } from "expo-router"
import { useSafeAreaInsets } from "react-native-safe-area-context"
import { useTheme } from "@/theme/ThemeProvider"
import { BackButton } from "@/components/ui/BackButton"
import { CTAButton } from "@/components/ui/CTAButton"
import { Icon } from "@/components/ui/Icon"
import { PassphraseGrid } from "@/components/auth/PassphraseGrid"
import { PassphraseActions } from "@/components/auth/PassphraseActions"
import { apiClient } from "@/api/client"
import * as FileSystem from 'expo-file-system'
import * as Sharing from 'expo-sharing'
import * as Clipboard from "expo-clipboard"

export default function BackupPassphraseScreen() {
  const router = useRouter()
  const theme = useTheme()
  const insets = useSafeAreaInsets()
  const [passphrase, setPassphrase] = useState<string[]>([])
  const [loading, setLoading] = useState(true)
  const [saving, setSaving] = useState(false)

  useEffect(() => {
    generatePassphrase()
  }, [])

  const generatePassphrase = async () => {
    try {
      const result = await apiClient.generatePassphrase()
      if (result.success && result.data) {
        setPassphrase(result.data.words)
      } else {
        // Fallback demo words
        setPassphrase([
          'abandon', 'ability', 'able', 'about',
          'above', 'absent', 'absorb', 'abstract',
          'absurd', 'abuse', 'access', 'accident'
        ])
      }
    } catch (error) {
      console.error('Error generating passphrase:', error)
    } finally {
      setLoading(false)
    }
  }

  const handleCopy = async () => {
    const passphraseText = passphrase.join(' ')
    await Clipboard.setStringAsync(passphraseText)
    Alert.alert('Copied', 'Passphrase copied to clipboard')
  }

  const handleSave = async () => {
    setSaving(true)
    try {
      const passphraseText = passphrase.join(' ')
      const fileName = 'verza-passphrase.txt'
      const fileUri = FileSystem.documentDirectory + fileName
      
      await FileSystem.writeAsStringAsync(fileUri, passphraseText)
      await Sharing.shareAsync(fileUri)
    } catch (error) {
      Alert.alert('Error', 'Failed to save passphrase')
    } finally {
      setSaving(false)
    }
  }

  const handleContinue = () => {
    router.push('/(auth)/confirm-passphrase')
  }

  return (
    <KeyboardAvoidingView 
      style={{ flex: 1 }} 
      behavior={Platform.OS === "ios" ? "padding" : "height"}
    >
      <View
        style={{
          flex: 1,
          backgroundColor: theme.isDark ? theme.colors.backgroundDark : theme.colors.backgroundLight,
          paddingTop: insets.top,
          paddingBottom: insets.bottom,
          paddingHorizontal: 20,
        }}
      >
        {/* Header */}
        <View style={{ paddingVertical: 16 }}>
          <BackButton />
        </View>

        {/* Content */}
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
              Backup Passphrase
            </Text>
            <Text
              style={{
                fontSize: 16,
                color: theme.colors.textSecondary,
                marginBottom: 4,
              }}
            >
              Write down these 12 words in order.
            </Text>
            <Text
              style={{
                fontSize: 16,
                color: theme.colors.textSecondary,
              }}
            >
              You'll need them to recover your wallet.
            </Text>
          </View>

          <View style={{ marginBottom: 32 }}>
            <PassphraseGrid words={passphrase} loading={loading} />
          </View>

          <View style={{ marginBottom: 32 }}>
            <PassphraseActions 
              onCopy={handleCopy}
              onSave={handleSave}
              saving={saving}
            />
          </View>

          <CTAButton
            title="I've Saved My Passphrase"
            onPress={handleContinue}
            disabled={loading || passphrase.length === 0}
          />
        </View>
      </View>
    </KeyboardAvoidingView>
  )
}