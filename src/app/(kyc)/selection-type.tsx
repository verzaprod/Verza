import { useState } from "react"
import { View, Text, SafeAreaView, KeyboardAvoidingView, Platform } from "react-native"
import { useRouter } from "expo-router"
import { useTheme } from "@/theme/ThemeProvider"
import { CTAButton } from "@/components/ui/CTAButton"
import { IDTypeCard } from "@/components/kyc/IDTypeCard"
import { useSafeAreaInsets } from "react-native-safe-area-context"
import { useKYCStore } from "@/store/kycStore"

type IDType = 'passport' | 'driver-license' | 'id-card'

const idTypes = [
  {
    id: 'passport' as IDType,
    title: 'Passport',
    description: 'International travel document',
    icon: 'passport',
  },
  {
    id: 'driver-license' as IDType,
    title: "Driver's License",
    description: 'Government-issued driving permit',
    icon: 'driver-license',
  },
  {
    id: 'id-card' as IDType,
    title: 'National ID Card',
    description: 'Government-issued identity card',
    icon: 'id-card',
  },
]

export default function SelectionType() {
  const router = useRouter()
  const theme = useTheme()
  const insets = useSafeAreaInsets()
  const [selectedType, setSelectedType] = useState<IDType | null>('passport')

  const { setCurrentStep, setDocumentType } = useKYCStore();

  const handleContinue = () => {
    if (selectedType) {
      setDocumentType(selectedType)
      setCurrentStep('docs')
      router.push(`/(kyc)/doc-capture`)
    }
  }

  return (
    <SafeAreaView 
      style={{
        flex: 1, 
        paddingTop: insets.top, 
        backgroundColor: theme.colors.background,
      }}
    >
      <View className="flex-1 px-6"
        style={{
          paddingTop: 20,
        }}
      >
        <View className="mb-8">
          <Text
            className="text-3xl mb-3"
            style={{
              color: theme.colors.textPrimary,
              fontFamily: theme.fonts.welcomeHeading
            }}
          >
            Select ID Type
          </Text>
          <Text
            className="text-base leading-6"
            style={{
              color: theme.colors.textSecondary,
              fontFamily: theme.fonts.onboardingTagline
            }}
          >
            Choose the type of identification document you'd like to use
          </Text>
        </View>

        <View className="gap-8 mb-20">
          {idTypes.map((idType) => (
            <IDTypeCard
              key={idType.id}
              title={idType.title}
              description={idType.description}
              icon={idType.icon}
              selected={selectedType === idType.id}
              onPress={() => setSelectedType(idType.id)}
            />
          ))}
        </View>

        <CTAButton
          title="Continue"
          onPress={handleContinue}
          disabled={!selectedType}
        />
      </View>
    </SafeAreaView>
  )
}
