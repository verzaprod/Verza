import React from 'react'
import { View, TouchableOpacity } from 'react-native'
import { Icon } from '@/components/ui/Icon'
import { useTheme } from '@/theme/ThemeProvider'

export const VerifiersHeader: React.FC = () => {
  const theme = useTheme()

  return (
    <View 
      className="flex-row justify-between items-center"
      style={{ paddingVertical: 16 }}
    >
      <TouchableOpacity
        style={{
          width: 48,
          height: 48,
          borderRadius: 24,
          backgroundColor: theme.colors.primaryGreen,
          alignItems: 'center',
          justifyContent: 'center',
        }}
      >
        <Icon name="avatar" size={32} />
      </TouchableOpacity>
      
      <TouchableOpacity>
        <View
          style={{
            width: 24,
            height: 24,
            borderRadius: 12,
            alignItems: 'center',
            justifyContent: 'center',
          }}
        >
          <Icon name="notification" size={24} />
        </View>
      </TouchableOpacity>
    </View>
  )
}