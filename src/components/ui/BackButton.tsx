"use client"

import type React from "react"
import { TouchableOpacity } from "react-native"
import { useRouter } from "expo-router"
import { useTheme } from "@/theme/ThemeProvider"
import { Icon } from "./Icon"

interface BackButtonProps {
  onPress?: () => void
}

export const BackButton: React.FC<BackButtonProps> = ({ onPress }) => {
  const router = useRouter()
  const theme = useTheme()

  const handlePress = () => {
    if (onPress) {
      onPress()
    } else {
      router.back()
    }
  }

  return (
    <TouchableOpacity
      onPress={handlePress}
      style={{
        width: 40,
        height: 40,
        borderRadius: 20,
        backgroundColor: theme.colors.primaryGreen,
        alignItems: "center",
        justifyContent: "center",
      }}
    >
      <Icon name="chevron-left" size={12} color="#ffffff"/>
    </TouchableOpacity>
  )
}
