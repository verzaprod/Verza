import React from 'react'
import { View, Text } from 'react-native'
import { useTheme } from '@/theme/ThemeProvider'

interface PassphraseGridProps {
  words: string[]
  loading?: boolean
}

export const PassphraseGrid: React.FC<PassphraseGridProps> = ({ words, loading = false }) => {
  const theme = useTheme()

  if (loading) {
    return (
      <View className="flex-row flex-wrap gap-3">
        {Array.from({ length: 12 }).map((_, index) => (
          <View
            key={index}
            className="flex-1 min-w-[30%]"
            style={{
              height: 60,
              backgroundColor: `${theme.colors.textSecondary}1A`, // 10% opacity
              borderRadius: theme.borderRadius.md,
              justifyContent: 'center',
              alignItems: 'center',
            }}
          >
            <View
              style={{
                width: 80,
                height: 16,
                backgroundColor: theme.colors.textSecondary + '40',
                borderRadius: 8,
              }}
            />
          </View>
        ))}
      </View>
    )
  }

  return (
    <View className="flex-row flex-wrap gap-3">
      {words.map((word, index) => (
        <View
          key={index}
          className="flex-1 min-w-[30%]"
          style={{
            height: 60,
            backgroundColor: `${theme.colors.textSecondary}1A`, // 10% opacity
            borderRadius: theme.borderRadius.md,
            justifyContent: 'center',
            alignItems: 'center',
            paddingHorizontal: 12,
          }}
        >
          <Text
            style={{
              fontSize: 12,
              color: theme.colors.textSecondary,
              fontWeight: '500',
              marginBottom: 2,
            }}
          >
            {index + 1}
          </Text>
          <Text
            style={{
              fontSize: 16,
              color: theme.isDark ? theme.colors.textPrimaryDark : theme.colors.textPrimaryLight,
              fontWeight: '600',
              textAlign: 'center',
            }}
          >
            {word}
          </Text>
        </View>
      ))}
    </View>
  )
}