"use client"

import { useState } from "react"
import { View } from "react-native"
import { useRouter } from "expo-router"
import { useSafeAreaInsets } from "react-native-safe-area-context"
import { AnimatedSplash } from "@/components/AnimatedSplash"

export default function SplashScreen() {
  const router = useRouter()
  const insets = useSafeAreaInsets()
  const [showSplash, setShowSplash] = useState(true)

  const handleAnimationComplete = () => {
    setShowSplash(false)
    router.replace("/onboarding")
  }

  if (showSplash) {
    return (
      <View style={{ flex: 1, paddingTop: insets.top }}>
        <AnimatedSplash onAnimationComplete={handleAnimationComplete} />
      </View>
    )
  }

  return null
}
