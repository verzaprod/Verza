"use client"

import type React from "react"
import { TouchableOpacity, Text } from "react-native"
import { useTheme } from "@/theme/ThemeProvider"

interface SkipButtonProps {
  onPress: () => void
}

export const SkipButton: React.FC<SkipButtonProps> = ({ onPress }) => {
  const theme = useTheme()

  return (
    <TouchableOpacity onPress={onPress}>
      <Text
        style={{
          color: theme.isDark ? theme.colors.textPrimaryDark : theme.colors.textPrimaryLight,
          fontSize: 16,
        }}
      >
        Skip
      </Text>
    </TouchableOpacity>
  )
}
