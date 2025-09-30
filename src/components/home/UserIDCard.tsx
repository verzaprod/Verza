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
        shadowColor: 'black',
        shadowOffset: { width: 0, height: 0 },
        shadowOpacity: 0.15,
        shadowRadius: 40,
        elevation: 40,
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
        className="text-3xl mb-8 text-center"
        style={{ 
          color: 'white',
          fontFamily: theme.fonts.welcomeHeading,
        }}
      >
        did:verza:z6Mha...2doK
      </Text>

      <View className="flex-row justify-end">
        <Text
          className="text-xl font- opacity-90"
          style={{ color: 'white', fontFamily: theme.fonts.welcomeHeading }}
        >
          Verza
        </Text>
      </View>
    </View>
  )
}