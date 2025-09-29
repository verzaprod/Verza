import React from 'react'
import { View, Text, TouchableOpacity } from 'react-native'
import { useRouter } from 'expo-router'
import { useTheme } from '@/theme/ThemeProvider'
import { Icon } from '@/components/ui/Icon'

interface Verifier {
  id: string
  name: string
  type: string
  rating: number
  verified: number
  logo: string
  status: 'active' | 'busy' | 'offline'
  description: string
}

interface VerifierCardProps {
  verifier: Verifier
}

export const VerifierCard: React.FC<VerifierCardProps> = ({ verifier }) => {
  const theme = useTheme()
  const router = useRouter()

  const getStatusColor = () => {
    switch (verifier.status) {
      case 'active': return theme.colors.primaryGreen
      case 'busy': return '#F59E0B'
      case 'offline': return theme.colors.textSecondary
      default: return theme.colors.textSecondary
    }
  } 

  const handlePress = () => {
    router.push(`/(kyc)/escrow-confirmation?verifierId=${verifier.id}`)
  }

  return (
    <TouchableOpacity
      style={{
        backgroundColor: theme.colors.background,
        borderRadius: theme.borderRadius.lg,
        padding: theme.spacing.lg,
        shadowColor: theme.isDark ? "#fff" : "#000",
        shadowOffset: { width: 0, height: 2 },
        shadowOpacity: 0.1,
        shadowRadius: 8,
        elevation: 4,
        borderWidth: 1,
        borderColor: theme.colors.backgroundLight,
      }}
      onPress={handlePress}
    >
      <View className="flex-row items-start justify-between mb-3">
        <View className="flex-row items-center flex-1">
          <View
            style={{
              width: 48,
              height: 48,
              backgroundColor: theme.colors.primaryGreen + '20',
              borderRadius: 24,
              alignItems: 'center',
              justifyContent: 'center',
              marginRight: theme.spacing.md,
            }}
          >
            <Icon name={verifier.logo} size={24} />
          </View>

          <View className="flex-1">
            <View className="flex-row items-center mb-1">
              <Text
                style={{
                  fontSize: 18,
                  fontWeight: '600',
                  color: theme.colors.textPrimary,
                  marginRight: theme.spacing.sm,
                }}
              >
                {verifier.name}
              </Text>
              <View
                style={{
                  width: 8,
                  height: 8,
                  borderRadius: 4,
                  backgroundColor: getStatusColor(),
                }}
              />
            </View>
            
            <Text
              style={{
                fontSize: 14,
                color: theme.colors.textSecondary,
                marginBottom: 2,
              }}
            >
              {verifier.type} • {verifier.verified.toLocaleString()} verified
            </Text>
          </View>
        </View>

        <View className="flex-row items-center">
          <Text
            style={{
              fontSize: 14,
              color: theme.colors.primaryGreen,
              fontWeight: '600',
              marginRight: 4,
            }}
          >
            ⭐ {verifier.rating}
          </Text>
        </View>
      </View>

      <Text
        style={{
          fontSize: 14,
          color: theme.colors.textSecondary,
          lineHeight: 20,
          marginBottom: theme.spacing.sm,
        }}
      >
        {verifier.description}
      </Text>

      <View className="flex-row justify-between items-center">
        <Text
          style={{
            fontSize: 12,
            color: getStatusColor(),
            fontWeight: '500',
            textTransform: 'capitalize',
          }}
        >
          {verifier.status}
        </Text>

        <TouchableOpacity
          style={{
            paddingVertical: 8,
            paddingHorizontal: 16,
            backgroundColor: theme.colors.primaryGreen,
            borderRadius: theme.borderRadius.md,
          }}
          onPress={handlePress}
        >
          <Text
            style={{
              fontSize: 14,
              color: 'white',
              fontWeight: '600',
            }}
          >
            Select
          </Text>
        </TouchableOpacity>
      </View>
    </TouchableOpacity>
  )
}