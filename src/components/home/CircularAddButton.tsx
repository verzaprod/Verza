import React from 'react'
import { TouchableOpacity, View } from 'react-native'
import { useTheme } from '@/theme/ThemeProvider'
import { Icon } from '@/components/ui/Icon'

export const CircularAddButton: React.FC = () => {
  const theme = useTheme()

  const handlePress = () => {
    console.log('Add account pressed')
  }

  return (
    <TouchableOpacity
      style={{
        width: 56,
        height: 56,
        borderRadius: 28,
        borderWidth: 2,
        borderColor: theme.colors.primaryGreen,
        backgroundColor: theme.colors.primaryGreen,
        alignItems: 'center',
        justifyContent: 'center',
        shadowColor: '#000',
        shadowOffset: { width: 0, height: 4 },
        shadowOpacity: 0.2,
        shadowRadius: 8,
        elevation: 4,
      }}
      onPress={handlePress}
      accessible={true}
      accessibilityLabel="Add new account"
      accessibilityRole="button"
    >
      <Icon name="plus" size={24} color="white" />
    </TouchableOpacity>
  )
}
