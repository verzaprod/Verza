import React from 'react'
import { View, Text } from 'react-native'
import Svg, { Path } from 'react-native-svg'
import { useTheme } from '@/theme/ThemeProvider'
import { Icon } from '@/components/ui/Icon'

export const PatternedIDCard: React.FC = () => {
  const theme = useTheme()

  return (
    <View
      style={{
        borderRadius: theme.borderRadius.lg,
        backgroundColor: theme.colors.primaryGreen,
        padding: 24,
        overflow: 'hidden',
        position: 'relative',
        ...theme.shadows.subtle,
      }}
    >
      <Svg
        style={{ position: 'absolute', top: 0, left: 0 }}
        width="100"
        height="60"
        viewBox="0 0 100 60"
      >
        <Path
          d="M0,0 Q25,15 50,10 T100,20 L100,0 Z"
          fill={theme.colors.primaryGreen}
          opacity={0.8}
        />
      </Svg>
      
      <Svg
        style={{ position: 'absolute', bottom: 0, right: 0 }}
        width="120"
        height="80"
        viewBox="0 0 120 80"
      >
        <Path
          d="M120,80 Q90,65 60,70 T0,60 L0,80 Z"
          fill={theme.colors.primaryGreen}
          opacity={0.9}
        />
      </Svg>

      <View className="flex-row justify-between items-start mb-6">
        <Text style={{ color: 'rgba(255,255,255,0.8)', fontSize: 16 }}>
          UserID
        </Text>
        <View className="flex-row gap-1">
          {[...Array(3)].map((_, i) => (
            <View
              key={i}
              style={{
                width: 3,
                height: 8,
                backgroundColor: 'white',
                borderRadius: 2,
                opacity: 1 - (i * 0.2),
              }}
            />
          ))}
        </View>
      </View>

      <Text
        style={{
          color: 'white',
          fontSize: 28,
          fontWeight: 'bold',
          fontFamily: theme.fonts.welcomeHeading,
          marginBottom: 32,
        }}
      >
        did:verza:1234abcd
      </Text>

      <View className="flex-row justify-end">
        <Text
          style={{
            color: 'white',
            fontSize: 18,
            fontWeight: '600',
            opacity: 0.9,
          }}
        >
          Verza
        </Text>
      </View>
    </View>
  )
}
