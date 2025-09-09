import { useState, useEffect } from "react"
import { View, Text, KeyboardAvoidingView, Platform } from "react-native"
import { useRouter } from "expo-router"
import { useSafeAreaInsets } from "react-native-safe-area-context"
import { useTheme } from "@/theme/ThemeProvider"
import { CTAButton } from "@/components/ui/CTAButton"
import { WordChipGrid } from "@/components/auth/WordChipGrid"

const DEMO_WORDS = [
  'abandon', 'ability', 'able', 'about',
  'above', 'absent', 'absorb', 'abstract',
  'absurd', 'abuse', 'access', 'accident'
]

const CORRECT_SEQUENCE = ['abandon', 'ability', 'able'] // First 3 words for demo

export default function ConfirmPassphraseScreen() {
  const router = useRouter()
  const theme = useTheme()
  const insets = useSafeAreaInsets()
  const [shuffledWords, setShuffledWords] = useState<string[]>([])
  const [selectedWords, setSelectedWords] = useState<string[]>([])
  const [loading, setLoading] = useState(false)

  useEffect(() => {
    // Shuffle the words for selection
    const shuffled = [...DEMO_WORDS].sort(() => Math.random() - 0.5)
    setShuffledWords(shuffled)
  }, [])

  const handleWordSelect = (word: string) => {
    if (selectedWords.includes(word)) {
      // Remove word if already selected
      setSelectedWords(prev => prev.filter(w => w !== word))
    } else if (selectedWords.length < 3) {
      // Add word if less than 3 selected
      setSelectedWords(prev => [...prev, word])
    }
  }

  const isCorrectSequence = () => {
    return selectedWords.length === 3 && 
           selectedWords.every((word, index) => word === CORRECT_SEQUENCE[index])
  }

  const handleContinue = async () => {
    if (!isCorrectSequence()) return
    
    setLoading(true)
    try {
      // TODO: Implement API call to confirm passphrase
      // await authAPI.confirmPassphrase(selectedWords);
      
      setTimeout(() => {
        setLoading(false)
        router.replace('/(auth)/auth-success')
      }, 1500)
    } catch (error) {
      console.error('Error confirming passphrase:', error)
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
          paddingBottom: insets.bottom,
          paddingHorizontal: 20,
        }}
      >
        <View>
          <View style={{ alignItems: "flex-start", marginBottom: 40, marginTop: 40 }}>
            <Text
              style={{
                fontSize: 28,
                fontWeight: "bold",
                color: theme.colors.textPrimary,
                marginBottom: 8,
              }}
            >
              Confirm Passphrase
            </Text>
            <Text
              style={{
                fontSize: 16,
                color: theme.colors.textSecondary,
              }}
            >
              Select the first 3 words in the correct order to confirm your passphrase.
            </Text>
          </View>

          <View style={{ marginBottom: 40 }}>
            <WordChipGrid
              words={shuffledWords}
              selectedWords={selectedWords}
              onWordSelect={handleWordSelect}
            />
          </View>

          <CTAButton
            title="Continue"
            onPress={handleContinue}
            loading={loading}
            disabled={!isCorrectSequence()}
          />
        </View>
      </View>
    </KeyboardAvoidingView>
  )
}
