import { useState } from "react"
import { View, Text, SafeAreaView } from "react-native"
import { useRouter } from "expo-router"
import { useTheme } from "@/theme/ThemeProvider"
import { BackButton } from "@/components/ui/BackButton"
import { CTAButton } from "@/components/ui/CTAButton"
import { IDTypeCard } from "@/components/kyc/IDTypeCard"

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
  const [selectedType, setSelectedType] = useState<IDType | null>('passport')

  const handleContinue = () => {
    if (selectedType) {
      router.push(`/(kyc)/doc-capture?type=${selectedType}`)
    }
  }

  return (
    <SafeAreaView 
      style={{ 
        flex: 1, 
        backgroundColor: theme.colors.background 
      }}
    >
      <View className="flex-1 px-5">
        {/* Header */}
        <View className="py-4">
          <BackButton />
        </View>

        {/* Title Section */}
        <View className="mb-8">
          <Text
            className="text-3xl font-bold mb-3"
            style={{
              color: theme.isDark ? theme.colors.textPrimaryDark : theme.colors.textPrimaryLight,
            }}
          >
            Select ID Type
          </Text>
          <Text
            className="text-base leading-6"
            style={{
              color: theme.colors.textSecondary,
            }}
          >
            Choose the type of identification document you'd like to use
          </Text>
        </View>

        {/* ID Type Cards */}
        <View className="flex-1 gap-4">
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

        {/* Continue Button */}
        <View className="pb-6">
          <CTAButton
            title="Continue"
            onPress={handleContinue}
            disabled={!selectedType}
          />
        </View>
      </View>
    </SafeAreaView>
  )
}
