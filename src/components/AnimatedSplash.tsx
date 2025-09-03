import type React from "react"
import { useEffect } from "react"
import { View, Text } from "react-native"
import Animated, { useSharedValue, useAnimatedStyle, withTiming, withDelay } from "react-native-reanimated"
import { useTheme } from "@/theme/ThemeProvider"
import { Icon } from "./ui/Icon"

interface AnimatedSplashProps {
  onAnimationComplete: () => void
}

export const AnimatedSplash: React.FC<AnimatedSplashProps> = ({ onAnimationComplete }) => {
  const theme = useTheme()
  const logoTranslateX = useSharedValue(0)
  const logoOpacity = useSharedValue(1)
  const textOpacity = useSharedValue(0)
  const bottomImageOpacity = useSharedValue(0)
  const bottomImageTranslateY = useSharedValue(50)

  const logoAnimatedStyle = useAnimatedStyle(() => ({
    transform: [{ translateX: logoTranslateX.value }],
    opacity: logoOpacity.value,
  }))

  const textAnimatedStyle = useAnimatedStyle(() => ({
    opacity: textOpacity.value,
  }))

  const bottomImageAnimatedStyle = useAnimatedStyle(() => ({
    opacity: bottomImageOpacity.value,
    transform: [{ translateY: bottomImageTranslateY.value }],
  }))

  useEffect(() => {
    const animateSequence = () => {
      // Wait 1 second, then animate logo to left
      logoTranslateX.value = withDelay(1000, withTiming(-60, { duration: 800 }))

      // Show text after logo moves
      textOpacity.value = withDelay(1400, withTiming(1, { duration: 600 }))

      // Animate bottom image
      bottomImageOpacity.value = withDelay(1800, withTiming(1, { duration: 600 }))
      bottomImageTranslateY.value = withDelay(1800, withTiming(0, { duration: 600 }))

      // Complete animation after 3 seconds
      setTimeout(onAnimationComplete, 3000)
    }

    animateSequence()
  }, [])

  return (
    <View
      style={{
        flex: 1,
        backgroundColor: theme.isDark ? theme.colors.backgroundDark : theme.colors.backgroundLight,
        alignItems: "center",
        justifyContent: "center",
      }}
    >
      <View style={{ flexDirection: "row", alignItems: "center" }}>
        <Animated.View style={logoAnimatedStyle}>
          <Icon name="verza-logo" size={60} color={theme.colors.primaryGreen} />
        </Animated.View>
        <Animated.View style={[{ marginLeft: 16 }, textAnimatedStyle]}>
          <Text
            style={{
              fontSize: 32,
              fontWeight: "bold",
              color: theme.isDark ? theme.colors.textPrimaryDark : theme.colors.textPrimaryLight,
            }}
          >
            Verza
          </Text>
        </Animated.View>
      </View>

      <Animated.View
        style={[
          {
            position: "absolute",
            bottom: 60,
            right: 30,
          },
          bottomImageAnimatedStyle,
        ]}
      >
        <Icon name="splash-decoration" size={80} color={theme.colors.primaryGreen} />
      </Animated.View>
    </View>
  )
}
