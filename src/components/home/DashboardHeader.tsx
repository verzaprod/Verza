import React from 'react'
import { View, TouchableOpacity } from 'react-native'
import { Icon } from '@/components/ui/Icon'

export const DashboardHeader: React.FC = () => {
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
          backgroundColor: '#16A34A',
          alignItems: 'center',
          justifyContent: 'center',
        }}
      >
        <Icon name="avatar" size={32} />
      </TouchableOpacity>
      
      <TouchableOpacity style={{ position: 'relative' }}>
        <View
          style={{
            width: 24,
            height: 24,
            // backgroundColor: '#16A34A',
            borderRadius: 12,
            alignItems: 'center',
            justifyContent: 'center',
          }}
        >
          <Icon name="notification" size={24} />
        </View>
        {/* <View
          style={{
            position: 'absolute',
            top: -2,
            right: -2,
            width: 8,
            height: 8,
            backgroundColor: '#EF4444',
            borderRadius: 4,
          }}
        /> */}
      </TouchableOpacity>
    </View>
  )
}
