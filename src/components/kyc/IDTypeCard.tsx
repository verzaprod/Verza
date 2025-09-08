import React from 'react'
import { View, Text, TouchableOpacity } from 'react-native'
import { useTheme } from '@/theme/ThemeProvider'
import { Icon } from '@/components/ui/Icon'

interface IDTypeCardProps {
  title: string
  description: string
  icon: string
  selected: boolean
  onPress: () => void
}

export const IDTypeCard: React.FC<IDTypeCardProps> = ({
  title,
  description,
  icon,
  selected,
  onPress,
}) => {
  const theme = useTheme()

  return (
    <TouchableOpacity
      className="flex-row items-center p-4"
      style={{
        backgroundColor: theme.colors.background,
        borderRadius: theme.borderRadius.lg,
        borderWidth: 2,
        borderColor: selected 
          ? theme.colors.primaryGreen 
          : theme.colors.boxBorder,
        shadowColor: theme.shadows.subtle.shadowColor,
        shadowOffset: theme.shadows.subtle.shadowOffset,
        shadowOpacity: theme.shadows.subtle.shadowOpacity,
        shadowRadius: theme.shadows.subtle.shadowRadius,
        elevation: theme.shadows.subtle.elevation,
      }}
      onPress={onPress}
    >
      <View
        className="mr-4"
        style={{
          width: 48,
          height: 48,
          backgroundColor: `${theme.colors.primaryGreen}20`,
          borderRadius: theme.borderRadius.full,
          justifyContent: 'center',
          alignItems: 'center',
        }}
      >
        <Icon 
          name={icon} 
          size={24} 
          color={theme.colors.primaryGreen}
        />
      </View>

      <View className="flex-1">
        <Text
          className="text-lg font-se mibold mb-1"
          style={{
            color: theme.colors.textPrimary,
            fontFamily: theme.fonts.welcomeHeading,
          }}
        >
          {title}
        </Text>
        <Text
          className="text-sm"
          style={{
            color: theme.colors.textSecondary,
          }}
        >
          {description}
        </Text>
      </View>

    </TouchableOpacity>
  )
}
