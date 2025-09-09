import { useState, useEffect } from "react"
import { View, Text, SafeAreaView } from "react-native"
import { useRouter } from "expo-router"
import { useTheme } from "@/theme/ThemeProvider"
import { VerificationIcon } from "@/components/kyc/VerificationIcon"
import { ProgressBar } from "@/components/kyc/ProgressBar"
import { VerificationSteps } from "@/components/kyc/VerificationSteps"

const verificationSteps = [
  { id: 'document-captured', label: 'Document captured', completed: true },
  { id: 'document-analyzed', label: 'Document analyzed', completed: true },
  { id: 'verifying-identity', label: 'Verifying Identity', completed: false },
]

export default function VerificationState() {
  const router = useRouter()
  const theme = useTheme()
  const [progress, setProgress] = useState(98)
  const [currentSteps, setCurrentSteps] = useState(verificationSteps)

  useEffect(() => {
    // Simulate verification completion after 3 seconds
    const timer = setTimeout(() => {
      setProgress(100)
      setCurrentSteps(prev => 
        prev.map(step => 
          step.id === 'verifying-identity' 
            ? { ...step, completed: true }
            : step
        )
      )
      
      // Redirect to success page after completing verification
      setTimeout(() => {
        router.push('/(kyc)/kyc-success')
      }, 1000)
    }, 3000)

    return () => clearTimeout(timer)
  }, [router])

  return (
    <SafeAreaView 
      style={{ 
        flex: 1, 
        backgroundColor: theme.colors.background 
      }}
    >
      <View className="flex-1 px-6 justify-center">
        {/* Verification Icon */}
        <View className="items-center mb-12">
          <VerificationIcon />
        </View>

        <View className="items-center mb-16">
          <Text
            className="text-3xl mb-4 text-center"
            style={{
              color: theme.colors.textPrimary,
              fontFamily: theme.fonts.welcomeHeading,
            }}
          >
            Verifying Your Identity
          </Text>
          <Text
            className="text-base text-center leading-6 px-4"
            style={{
              color: theme.colors.textSecondary,
              fontFamily: theme.fonts.onboardingTagline,
            }}
          >
            Please wait while we verify your documents and identity
          </Text>
        </View>

        <View className="mb-12">
          <ProgressBar progress={progress} />
        </View>

        <View>
          <VerificationSteps steps={currentSteps} />
        </View>
      </View>
    </SafeAreaView>
  )
}