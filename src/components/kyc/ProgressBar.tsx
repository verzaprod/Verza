import React from 'react'
import { View, Text } from 'react-native'
import { useTheme } from '@/theme/ThemeProvider'

interface ProgressBarProps {
  progress: number
}

export const ProgressBar: React.FC<ProgressBarProps> = ({ progress }) => {
  const theme = useTheme()

  return (
    <View>
      <View className="flex-row justify-between items-center mb-3">
        <Text
          className="text-base"
          style={{
            color: theme.colors.textSecondary,
          }}
        >
          Processing...
        </Text>
        <Text
          className="text-base font-semibold"
          style={{
            color: theme.colors.textSecondary,
          }}
        >
          {progress}%
        </Text>
      </View>

      <View
        className="w-full"
        style={{
          height: 8,
          backgroundColor: `${theme.colors.textSecondary}20`,
          borderRadius: 4,
          overflow: 'hidden',
        }}
      >
        <View
          style={{
            width: `${progress}%`,
            height: '100%',
            backgroundColor: theme.colors.primaryGreen,
            borderRadius: 4,
          }}
        />
      </View>
    </View>
  )
}
