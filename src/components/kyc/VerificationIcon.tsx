import React from 'react'
import { View } from 'react-native'
import { useTheme } from '@/theme/ThemeProvider'
import { Icon } from '@/components/ui/Icon'

export const VerificationIcon: React.FC = () => {
  const theme = useTheme()

  return (
    <View
      className="items-center justify-center"
      style={{
        width: 100,
        height: 100,
        backgroundColor: 'white',
        borderRadius: 60,
        ...theme.shadows.subtle,
      }}
    >
      <View
        style={{
          width: 80,
          height: 80,
          // backgroundColor: `${theme.colors.primaryGreen}20`,
          borderRadius: 40,
          alignItems: 'center',
          justifyContent: 'center',
        }}
      >
        <Icon 
          name="shield" 
          size={48}
          color={theme.colors.primaryGreen}
        />
      </View>
    </View>
  )
}
