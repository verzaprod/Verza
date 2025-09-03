"use client"

import type React from "react"
import { TouchableOpacity, View } from "react-native"
import Svg, { Circle } from "react-native-svg"
import { useTheme } from "@/theme/ThemeProvider"
import { Icon } from "@/components/ui/Icon"

interface CircularNextButtonProps {
  onPress: () => void
  progress: number // 0 to 1
}

export const CircularNextButton: React.FC<CircularNextButtonProps> = ({ onPress, progress }) => {
  const theme = useTheme()
  const size = 60
  const strokeWidth = 3
  const radius = (size - strokeWidth) / 2
  const circumference = radius * 2 * Math.PI
  const strokeDasharray = circumference
  const strokeDashoffset = circumference - progress * circumference

  return (
    <TouchableOpacity
      onPress={onPress}
      style={{
        width: size,
        height: size,
        alignItems: "center",
        justifyContent: "center",
      }}
    >
      <Svg width={size} height={size} style={{ position: "absolute" }}>
        <Circle cx={size / 2} cy={size / 2} r={radius} stroke="#E2E8F0" strokeWidth={strokeWidth} fill="none" />
        <Circle
          cx={size / 2}
          cy={size / 2}
          r={radius}
          stroke={theme.colors.primaryGreen}
          strokeWidth={strokeWidth}
          fill="none"
          strokeDasharray={strokeDasharray}
          strokeDashoffset={strokeDashoffset}
          strokeLinecap="round"
          transform={`rotate(-90 ${size / 2} ${size / 2})`}
        />
      </Svg>
      <View
        style={{
          width: size - 12,
          height: size - 12,
          borderRadius: (size - 12) / 2,
          backgroundColor: theme.colors.primaryGreen,
          alignItems: "center",
          justifyContent: "center",
        }}
      >
        <Icon name="chevron-right" size={20} color="#FFFFFF" />
      </View>
    </TouchableOpacity>
  )
}
