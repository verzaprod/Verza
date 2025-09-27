import React from 'react'
import { TouchableOpacity, View } from 'react-native'
import { useTheme } from '@/theme/ThemeProvider'
import { Icon } from '@/components/ui/Icon'

export const AddAccountButton: React.FC = () => {
  const theme = useTheme()

  const handleAddAccount = () => {
    // TODO: Implement add account functionality
    console.log('Add account pressed')
  }

  return (
    <TouchableOpacity
      className="items-center justify-center"
      style={{
        width: 90,
        height: 90,
        borderRadius: 50,
        borderWidth: 3,
        borderColor: theme.colors.primaryGreen,
        backgroundColor: 'transparent',
      }}
      onPress={handleAddAccount}
    >
      <View
        className="items-center justify-center"
        style={{
          width: 50,
          height: 50,
          borderRadius: 30,
          backgroundColor: theme.colors.primaryGreen,
        }}
      >
        <Icon name="plus" size={24} color="white" />
      </View>
    </TouchableOpacity>
  )
}
