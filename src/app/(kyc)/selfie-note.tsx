import { View, Text, SafeAreaView } from "react-native"
import { useRouter } from "expo-router"
import { useTheme } from "@/theme/ThemeProvider"
// import { BackButton } from "@/components/ui/BackButton"
import { CTAButton } from "@/components/ui/CTAButton"
import { SelfieIllustration } from "@/components/kyc/SelfieIllustration"
import { InstructionList } from "@/components/kyc/InstructionList"

const instructions = [
  "Look directly at the camera",
  "Ensure good lighting on your face", 
  "Remove glasses and hat if possible",
  "Keep a neutral expression"
]

export default function SelfieNote() {
  const router = useRouter()
  const theme = useTheme()

  const handleTakeSelfie = () => {
    router.push("/(kyc)/selfie-capture")
  }

  return (
    <SafeAreaView 
      style={{ 
        flex: 1, 
        backgroundColor: theme.colors.background 
      }}
    >
      <View className="flex-1 px-6 justify-center">

        <View className="items-center">
          <View className="mb-12">
            <SelfieIllustration />
          </View>

          <View className="mb-8 items-center">
            <Text
              className="text-3xl mb-3 text-center"
              style={{
                color: theme.colors.textPrimary,
                fontFamily: theme.fonts.welcomeHeading,
              }}
            >
              Take a Selfie
            </Text>
            <Text
              className="text-base text-center leading-6"
              style={{
                color: theme.colors.textSecondary,
                fontFamily: theme.fonts.onboardingTagline,
              }}
            >
              We need to verify that you're the person in the document
            </Text>
          </View>

          <View className="mb-12 w-full">
            <InstructionList instructions={instructions} />
          </View>
        </View>

        <View className="pb-6 mt-8">
          <CTAButton
            title="Take Selfie"
            onPress={handleTakeSelfie}
          />
        </View>
      </View>
    </SafeAreaView>
  )
}
