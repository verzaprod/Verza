import React from 'react'
import { TouchableOpacity } from 'react-native'
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
        width: 80,
        height: 80,
        borderRadius: 40,
        borderWidth: 3,
        borderColor: theme.colors.primaryGreen,
        backgroundColor: 'transparent',
      }}
      onPress={handleAddAccount}
    >
      <TouchableOpacity
        className="items-center justify-center"
        style={{
          width: 60,
          height: 60,
          borderRadius: 30,
          backgroundColor: theme.colors.primaryGreen,
        }}
      >
        <Icon name="plus" size={24} color="white" />
      </TouchableOpacity>
    </TouchableOpacity>
  )
}
