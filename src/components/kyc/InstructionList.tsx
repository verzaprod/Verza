import React from 'react'
import { View, Text } from 'react-native'
import { useTheme } from '@/theme/ThemeProvider'
// import { Icon } from '@/components/ui/Icon'

interface InstructionListProps {
  instructions: string[]
}

export const InstructionList: React.FC<InstructionListProps> = ({ instructions }) => {
  const theme = useTheme()

  return (
    <View className="gap-4">
      {instructions.map((instruction, index) => (
        <View key={index} className="flex-row items-center">
          <View
            className="mr-4"
            style={{
              width: 24,
              height: 24,
              backgroundColor: theme.colors.primaryGreen,
              borderRadius: 12,
              alignItems: 'center',
              justifyContent: 'center',
            }}
          >
            <View
              style={{
                width: 12,
                height: 12,
                backgroundColor: 'white',
                borderRadius: 8,
              }}
            />
          </View>

          <Text
            className="flex-1 text-base"
            style={{
              color: theme.colors.textSecondary,
              fontSize: 16,
            }}
          >
            {instruction}
          </Text>
        </View>
      ))}
    </View>
  )
}
