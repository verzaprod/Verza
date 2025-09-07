import React from 'react'
import { View, TouchableOpacity, Text } from 'react-native'
import { useTheme } from '@/theme/ThemeProvider'

interface SelectedWordsDisplayProps {
  selectedWords: string[]
  onRemoveWord: (word: string) => void
}

export const SelectedWordsDisplay: React.FC<SelectedWordsDisplayProps> = ({ 
  selectedWords, 
  onRemoveWord 
}) => {
  const theme = useTheme()

  return (
    <View>
      <Text
        style={{
          fontSize: 16,
          color: theme.colors.textSecondary,
          marginBottom: 12,
        }}
      >
        Selected words ({selectedWords.length}/3):
      </Text>
      
      <View className="flex-row gap-3">
        {Array.from({ length: 3 }).map((_, index) => {
          const word = selectedWords[index]
          
          return (
            <TouchableOpacity
              key={index}
              className="flex-1"
              style={{
                height: 50,
                backgroundColor: word 
                  ? theme.colors.primaryGreen 
                  : `${theme.colors.textSecondary}20`,
                borderRadius: theme.borderRadius.md,
                borderWidth: 2,
                borderColor: word 
                  ? theme.colors.primaryGreen 
                  : theme.colors.textSecondary + '40',
                borderStyle: word ? 'solid' : 'dashed',
                justifyContent: 'center',
                alignItems: 'center',
              }}
              onPress={() => word && onRemoveWord(word)}
              disabled={!word}
            >
              <Text
                style={{
                  fontSize: 14,
                  fontWeight: '600',
                  color: word ? 'white' : theme.colors.textSecondary,
                }}
              >
                {word || `${index + 1}`}
              </Text>
            </TouchableOpacity>
          )
        })}
      </View>
    </View>
  )
}