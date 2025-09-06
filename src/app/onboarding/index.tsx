import { useState } from "react"
import { View, Text, ScrollView, Pressable, Dimensions } from "react-native"
import { useRouter } from "expo-router"
import { useSafeAreaInsets } from "react-native-safe-area-context"
import { useTheme } from "@/theme/ThemeProvider"
import { BackButton } from "@/components/ui/BackButton"
import { SkipButton } from "@/components/ui/SkipButton"
import { CircularNextButton } from "@/components/onboarding/CircularNextButton"
import { Icon } from "@/components/ui/Icon"
import { WIDTH, HEIGHT } from "@/constants"

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
  const width = WIDTH
  const height = HEIGHT

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
      className="flex-1"
      style={{
        backgroundColor: theme.isDark ? theme.colors.backgroundDark : theme.colors.backgroundLight,
        paddingTop: insets.top,
        paddingBottom: insets.bottom,
      }}
    >
      <View className="flex-row justify-between items-center px-5 py-4">
        {(currentPage != 0) && <BackButton onPress={handleBack} />}
        {(currentPage == 0) && <View />}
        <SkipButton onPress={handleSkip} />
      </View>

      {(currentPage != 1) && 
        (<ScrollView
          contentContainerStyle={{
            flex: 1,
            alignItems: "center",
          }}
        >
          <View className="w-full items-center mb-2 mt-[-20px] ml-4">
            <Icon 
              name={currentData.image} 
              style={{
                width,
                height
              }}
            />
          </View>

          <View className="items -center mb-10">
            <Text className="text-4xl font -bold text-center mb-4 px-2"
              style={{
                fontFamily: theme.fonts.onboardingHeading,
                color: theme.colors.textPrimary
              }}
            >
              {currentData.title}
            </Text>
            <Text className="text-base text-center px-5"
              style={{
                fontFamily: theme.fonts.onboardingTagline,
                color: theme.colors.textPrimary,
                textAlign: "center",
                lineHeight: 24,
              }}
            >
              {currentData.subtitle}
            </Text>
          </View>

          <CircularNextButton onPress={handleNext} progress={progress} />
        </ScrollView>)
      } 

      {(currentPage === 1) &&
        (<ScrollView
          contentContainerStyle={{
            flex: 1,
            alignItems: "center",
            justifyContent: "center",
            paddingHorizontal: 20,
          }}
        >
          <View className="items-center mb-[-42px]">
            <Text
              className="text-4xl font-bold text-center mb-4 px-5"
              style={{
                fontFamily: theme.fonts.onboardingHeading,
                color: theme.colors.textPrimary,
              }}
            >
              {currentData.title}
            </Text>

            <Text
              className="text-center px-5"
              style={{
                fontFamily: theme.fonts.onboardingTagline,
                color: theme.colors.textPrimary,
                lineHeight: 24,
              }}
            >
              {currentData.subtitle}
            </Text>
          </View>

          <View className="items-center mt-0">
            <Icon 
              name={currentData.image} 
              style={{
                width,
                height
              }} 
            />
          </View>

          <CircularNextButton onPress={handleNext} progress={progress} />
        </ScrollView>)
      }
    </View>
  )
}
