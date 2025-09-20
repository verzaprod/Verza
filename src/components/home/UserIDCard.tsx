import React from 'react'
import { View, Text } from 'react-native'
import { useTheme } from '@/theme/ThemeProvider'
import { Icon } from '@/components/ui/Icon'

export const UserIDCard: React.FC = () => {
  const theme = useTheme()

  return (
    <View
      className="p-6 rounded-3xl"
      style={{
        backgroundColor: theme.colors.primaryGreen,
        shadowColor: '#000',
        shadowOffset: { width: 0, height: 8 },
        shadowOpacity: 0.15,
        shadowRadius: 20,
        elevation: 8,
      }}
    >
      <View className="flex-row justify-between items-start mb-6">
        <Text
          className="text-lg opacity-80"
          style={{ color: 'white' }}
        >
          UserID
        </Text>
        <Icon name="wifi" size={24} />
      </View>

      <Text
        className="text-3xl font-bold mb-8"
        style={{ 
          color: 'white',
          fontFamily: theme.fonts.welcomeHeading,
        }}
      >
        did:verza:1234abcd
      </Text>

      <View className="flex-row justify-end">
        <Text
          className="text-xl font-semibold opacity-90"
          style={{ color: 'white' }}
        >
          Verza
        </Text>
      </View>
    </View>
  )
}