import React from 'react'
import { View } from 'react-native'
import { Icon } from '@/components/ui/Icon'

export const SelfieIllustration: React.FC = () => {
  return (
    <View className="relative items-center justify-center">
      <View 
        className="absolute left-0"
        style={{ top: 0, left: -120 }}
      >
        <Icon 
          name="left-pattern" 
          style={{ width: 120, height: 200 }}
        />
      </View>

      <View 
        className="absolute right-0"
        style={{ top: 40, right: -120 }}
      >
        <Icon 
          name="right-pattern" 
          style={{ width: 120, height: 200 }}
        />
      </View>

      <View
        className="bg-gray-200 rounded-3xl p-6 mx-4"
        style={{
          width: 140,
          height: 160,
          alignItems: 'center',
          justifyContent: 'center',
        }}
      >
        <Icon 
          name="face" 
          style={{ width: 140, height: 140 }}
        />
      </View>
    </View>
  )
}