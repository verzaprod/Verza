import React from 'react'
import { TouchableOpacity, Text, View } from 'react-native'
import { useTheme } from '@/theme/ThemeProvider'
import { Icon } from '@/components/ui/Icon'

interface AddCredentialButtonProps {
  onPress?: () => void
}

export const AddCredentialButton = ({ onPress }: AddCredentialButtonProps) => {
  const theme = useTheme()

  const handleAddCredential = () => {
    console.log('Add credential pressed')
  }

  return (
    <TouchableOpacity
      style={{
        flexDirection: 'row',
        alignItems: 'center',
        justifyContent: 'center',
        paddingVertical: theme.spacing.lg,
        paddingHorizontal: theme.spacing.xl,
        borderWidth: 2,
        borderColor: theme.colors.primaryGreen,
        borderStyle: 'dashed',
        borderRadius: theme.borderRadius.lg,
        backgroundColor: 'transparent',
        width: '100%',
      }}
      onPress={onPress}
    >
      <View
        style={{
          width: 32,
          height: 32,
          backgroundColor: theme.colors.primaryGreen,
          borderRadius: 16,
          alignItems: 'center',
          justifyContent: 'center',
          marginRight: theme.spacing.sm,
        }}
      >
        <Icon name="plus" size={20} color="white" />
      </View>
      
      <Text
        style={{
          fontSize: 18,
          fontWeight: '600',
          color: theme.colors.primaryGreen,
        }}
      >
        Add Credential
      </Text>
    </TouchableOpacity>
  )
}
