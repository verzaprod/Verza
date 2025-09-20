import React from 'react'
import { View, Text, TouchableOpacity } from 'react-native'
import { useTheme } from '@/theme/ThemeProvider'
import { Icon } from '@/components/ui/Icon'

interface VerifiedAccount {
  id: string
  name: string
  status: string
}

interface VerifiedAccountsListProps {
  accounts: VerifiedAccount[]
}

export const VerifiedAccountsList: React.FC<VerifiedAccountsListProps> = ({ accounts }) => {
  const theme = useTheme()

  return (
    <View
      className="p-4 rounded-2xl"
      style={{
        backgroundColor: `${theme.colors.textSecondary}10`,
      }}
    >
      {accounts.map((account, index) => (
        <View
          key={account.id}
          className="flex-row items-center justify-between py-4"
          style={{
            borderBottomWidth: index < accounts.length - 1 ? 1 : 0,
            borderBottomColor: `${theme.colors.textSecondary}20`,
          }}
        >
          <View className="flex-row items-center flex-1">
            <View
              className="mr-4 p-2 rounded-lg"
              style={{
                backgroundColor: theme.colors.textPrimary,
              }}
            >
              <Icon name="cancel" size={20} />
            </View>
            
            <View className="flex-1">
              <Text
                className="text-lg font-semibold mb-1"
                style={{
                  color: theme.colors.textPrimary,
                }}
              >
                {account.name}
              </Text>
              {index === accounts.length - 1 && (
                <View className="flex-row justify-between items-center">
                  <Text
                    className="text-sm"
                    style={{
                      color: theme.colors.textSecondary,
                      fontStyle: 'italic',
                    }}
                  >
                    Verified
                  </Text>
                  <TouchableOpacity>
                    <Text
                      className="text-sm font-semibold"
                      style={{
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

          <TouchableOpacity className="p-2">
            <Icon name="remove" size={20} />
          </TouchableOpacity>
        </View>
      ))}
    </View>
  )
}
