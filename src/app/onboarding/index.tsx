import { useState } from "react"
import { View, Text, ScrollView } from "react-native"
import { useRouter } from "expo-router"
import { useSafeAreaInsets } from "react-native-safe-area-context"
import { useTheme } from "@/theme/ThemeProvider"
import { BackButton } from "@/components/ui/BackButton"
import { SkipButton } from "@/components/ui/SkipButton"
import { CircularNextButton } from "@/components/onboarding/CircularNextButton"
import { Icon } from "@/components/ui/Icon"

const onboardingData = [
  {
    title: "Onboard Smarter",
    subtitle: "Skip the long forms. With Verza, onboarding takes minutes, not hours.",
    image: "onboarding-1",
  },
  {
    title: "Identity, Simplified",
    subtitle: "One secure ID, verified once, used everywhere - safely and instantly.",
    image: "onboarding-2",
  },
  {
    title: "Seamless Access",
    subtitle: "Move freely across platforms with one account. Verza keeps it effortless.",
    image: "onboarding-3",
  },
]

export default function OnboardingScreen() {
  const router = useRouter()
  const theme = useTheme()
  const insets = useSafeAreaInsets()
  const [currentPage, setCurrentPage] = useState(0)

  const handleNext = () => {
    if (currentPage < onboardingData.length - 1) {
      setCurrentPage(currentPage + 1)
    } else {
      router.push("/register")
    }
  }

  const handleSkip = () => {
    router.push("/register")
  }

  const handleBack = () => {
    if (currentPage > 0) {
      setCurrentPage(currentPage - 1)
    } else {
      router.back()
    }
  }

  const progress = (currentPage + 1) / onboardingData.length
  const currentData = onboardingData[currentPage]

  return (
    <View
      style={{
        flex: 1,
        backgroundColor: theme.isDark ? theme.colors.backgroundDark : theme.colors.backgroundLight,
        paddingTop: insets.top,
        paddingBottom: insets.bottom,
      }}
    >
      <View
        style={{
          flexDirection: "row",
          justifyContent: "space-between",
          alignItems: "center",
          paddingHorizontal: 20,
          paddingVertical: 16,
        }}
      >
        <BackButton onPress={handleBack} />
        <SkipButton onPress={handleSkip} />
      </View>

      <ScrollView
        contentContainerStyle={{
          flex: 1,
          alignItems: "center",
          justifyContent: "center",
          paddingHorizontal: 20,
        }}
      >
        <View style={{ alignItems: "center", marginBottom: 60 }}>
          <Icon name={currentData.image} size={200} color={theme.colors.primaryGreen} />
        </View>

        <View style={{ alignItems: "center", marginBottom: 80 }}>
          <Text
            style={{
              fontSize: 28,
              fontWeight: "bold",
              color: theme.isDark ? theme.colors.textPrimaryDark : theme.colors.textPrimaryLight,
              textAlign: "center",
              marginBottom: 16,
            }}
          >
            {currentData.title}
          </Text>
          <Text
            style={{
              fontSize: 16,
              color: theme.colors.textSecondary,
              textAlign: "center",
              lineHeight: 24,
              paddingHorizontal: 20,
            }}
          >
            {currentData.subtitle}
          </Text>
        </View>

        <CircularNextButton onPress={handleNext} progress={progress} />
      </ScrollView>
    </View>
  )
}
