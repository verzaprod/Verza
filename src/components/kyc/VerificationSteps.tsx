import React from 'react'
import { View, Text } from 'react-native'
import { useTheme } from '@/theme/ThemeProvider'

interface VerificationStep {
  id: string
  label: string
  completed: boolean
}

interface VerificationStepsProps {
  steps: VerificationStep[]
}

export const VerificationSteps: React.FC<VerificationStepsProps> = ({ steps }) => {
  const theme = useTheme()

  return (
    <View className="gap-4">
      {steps.map((step, index) => (
        <View key={step.id} className="flex-row items-center">
          <View
            className="mr-4"
            style={{
              width: 24,
              height: 24,
              backgroundColor: step.completed 
                ? theme.colors.primaryGreen 
                : `${theme.colors.textSecondary}30`,
              borderRadius: 12,
              alignItems: 'center',
              justifyContent: 'center',
            }}
          >
            {step.completed ? (
              <Text
                style={{
                  color: 'white',
                  fontSize: 14,
                  fontWeight: 'bold',
                }}
              >
                ✓
              </Text>
            ) : (
              <View
                style={{
                  width: 8,
                  height: 8,
                  backgroundColor: theme.colors.textSecondary,
                  borderRadius: 4,
                }}
              />
            )}
          </View>

          <Text
            className="flex-1 text-base"
            style={{
              color: step.completed 
                ? theme.colors.textSecondary
                : `${theme.colors.textSecondary}80`,
              fontWeight: step.completed ? '500' : '400',
            }}
          >
            {step.label}
          </Text>
        </View>
      ))}
    </View>
  )
}
