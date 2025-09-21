import React from 'react'
import { View, Text, TouchableOpacity } from 'react-native'
import { useTheme } from '@/theme/ThemeProvider'
import { Icon } from '@/components/ui/Icon'

interface VerifiedAccount {
  id: string
  name: string
  status: string
}

interface OverlappingAccountsListProps {
  accounts: VerifiedAccount[]
}

export const AccountsList: React.FC<OverlappingAccountsListProps> = ({ 
  accounts 
}) => {
  const theme = useTheme()

  return (
    <View>
      {accounts.map((account, index) => {
        const isLast = index === accounts.length - 1
        const zIndex = accounts.length - index
        const opacity = 1 - (index * 0.05)
        const scale = 1 - (index * 0.02)
        const marginTop = index === 0 ? 0 : -8

        return (
          <TouchableOpacity
            key={account.id}
            style={{
              backgroundColor: theme.isDark 
                ? 'rgba(255,255,255,0.1)' 
                : 'rgba(0,0,0,0.05)',
              borderRadius: theme.borderRadius.lg,
              padding: 16,
              flexDirection: 'row',
              alignItems: 'center',
              justifyContent: 'space-between',
              zIndex,
              opacity,
              transform: [{ scale }],
              marginTop,
              shadowColor: '#000',
              shadowOffset: { width: 0, height: 2 * index + 2 },
              shadowOpacity: 0.1 * (zIndex / accounts.length),
              shadowRadius: 4,
              elevation: zIndex,
            }}
            hitSlop={{ top: 8, bottom: 8, left: 8, right: 8 }}
          >
            <View className="flex-row items-center flex-1">
              <View
                style={{
                  width: 40,
                  height: 40,
                  backgroundColor: theme.colors.textPrimary,
                  borderRadius: 8,
                  alignItems: 'center',
                  justifyContent: 'center',
                  marginRight: 16,
                }}
              >
                <Icon name="cancel" size={20} />
              </View>
              
              <View className="flex-1">
                <Text
                  style={{
                    fontSize: 18,
                    fontWeight: '600',
                    color: theme.colors.textPrimary,
                    marginBottom: 4,
                  }}
                >
                  {account.name}
                </Text>
                {isLast && (
                  <View className="flex-row justify-between items-center">
                    <Text
                      style={{
                        fontSize: 14,
                        color: theme.colors.textSecondary,
                        fontStyle: 'italic',
                      }}
                    >
                      Verified
                    </Text>
                    <TouchableOpacity>
                      <Text
                        style={{
                          fontSize: 14,
                          fontWeight: '600',
                          color: theme.colors.primaryGreen,
                        }}
                      >
                        View Details
                      </Text>
                    </TouchableOpacity>
                  </View>
                )}
              </View>
            </View>

            <TouchableOpacity
              style={{
                width: 32,
                height: 32,
                borderRadius: 16,
                backgroundColor: theme.colors.textSecondary + '30',
                alignItems: 'center',
                justifyContent: 'center',
              }}
              hitSlop={{ top: 8, bottom: 8, left: 8, right: 8 }}
            >
              <Icon name="remove" size={16} />
            </TouchableOpacity>
          </TouchableOpacity>
        )
      })}
    </View>
  )
}