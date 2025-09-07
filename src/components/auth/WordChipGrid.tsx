import React from 'react'
import { View, TouchableOpacity, Text } from 'react-native'
import { useTheme } from '@/theme/ThemeProvider'

interface WordChipGridProps {
  words: string[]
  selectedWords: string[]
  onWordSelect: (word: string) => void
}

export const WordChipGrid: React.FC<WordChipGridProps> = ({ 
  words, 
  selectedWords, 
  onWordSelect 
}) => {
  const theme = useTheme()

  return (
    <View className="flex-row flex-wrap gap-3">
      {words.map((word, index) => {
        const isSelected = selectedWords.includes(word)
        
        return (
          <TouchableOpacity
            key={index}
            className="flex-1 min-w-[30%]"
            style={{
              paddingVertical: 12,
              paddingHorizontal: 16,
              backgroundColor: isSelected 
                ? theme.colors.secondaryAccent 
                : `${theme.colors.textSecondary}1A`,
              borderRadius: theme.borderRadius.md,
              borderWidth: isSelected ? 2 : 1,
              borderColor: isSelected 
                ? theme.colors.secondaryAccent 
                : theme.colors.textSecondary + '30',
              alignItems: 'center',
            }}
            onPress={() => onWordSelect(word)}
          >
            <Text
              style={{
                fontSize: 16,
                fontWeight: '600',
                color: isSelected 
                  ? 'white' 
                  : theme.isDark ? theme.colors.textPrimaryDark : theme.colors.textPrimaryLight,
              }}
            >
              {word}
            </Text>
          </TouchableOpacity>
        )
      })}
    </View>
  )
}