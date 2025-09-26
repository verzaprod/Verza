import React from 'react'
import { TouchableOpacity, Text, ViewStyle } from 'react-native'
import { useTheme } from '@/theme/ThemeProvider'

interface ButtonProps {
  text: string
  onPress: () => void
  variant?: 'primary' | 'secondary'
  disabled?: boolean
  style?: ViewStyle
}

export const Button: React.FC<ButtonProps> = ({
  text,
  onPress,
  variant = 'primary',
  disabled = false,
  style,
}) => {
  const theme = useTheme()

  const isPrimary = variant === 'primary'

  return (
    <TouchableOpacity
      style={[
        {
          paddingVertical: 16,
          paddingHorizontal: 24,
          borderRadius: theme.borderRadius.lg,
          alignItems: 'center',
          backgroundColor: isPrimary ? theme.colors.primaryGreen : 'transparent',
          borderWidth: isPrimary ? 0 : 1,
          borderColor: theme.colors.primaryGreen,
          opacity: disabled ? 0.6 : 1,
        },
        style,
      ]}
      onPress={onPress}
      disabled={disabled}
    >
      <Text
        style={{
          color: isPrimary ? 'white' : theme.colors.primaryGreen,
          fontSize: 16,
          fontWeight: '600',
        }}
      >
        {text}
      </Text>
    </TouchableOpacity>
  )
}
