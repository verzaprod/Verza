import React from 'react'
import { View, TouchableOpacity, Text } from 'react-native'
import { useTheme } from '@/theme/ThemeProvider'
import { Icon } from '@/components/ui/Icon'

interface PassphraseActionsProps {
  onCopy: () => void
  onSave: () => void
  saving?: boolean
}

export const PassphraseActions: React.FC<PassphraseActionsProps> = ({ 
  onCopy, 
  onSave, 
  saving = false 
}) => {
  const theme = useTheme()

  return (
    <View className="flex-row gap-4">
      <TouchableOpacity
        className="flex-1 flex-row items-center justify-center"
        style={{
          paddingVertical: 16,
          backgroundColor: theme.isDark ? theme.colors.backgroundLight + '10' : theme.colors.backgroundDark + '10',
          borderRadius: theme.borderRadius.md,
          borderWidth: 1,
          borderColor: theme.colors.textSecondary + '30',
        }}
        onPress={onCopy}
      >
        <Icon name="copy" size={20} color={theme.colors.textSecondary} />
        <Text
          style={{
            marginLeft: 8,
            fontSize: 16,
            fontWeight: '600',
            color: theme.colors.textSecondary,
          }}
        >
          Copy
        </Text>
      </TouchableOpacity>

      <TouchableOpacity
        className="flex-1 flex-row items-center justify-center"
        style={{
          paddingVertical: 16,
          backgroundColor: theme.isDark ? theme.colors.backgroundLight + '10' : theme.colors.backgroundDark + '10',
          borderRadius: theme.borderRadius.md,
          borderWidth: 1,
          borderColor: theme.colors.textSecondary + '30',
          opacity: saving ? 0.6 : 1,
        }}
        onPress={onSave}
        disabled={saving}
      >
        <Icon name="save" size={20} color={theme.colors.textSecondary} />
        <Text
          style={{
            marginLeft: 8,
            fontSize: 16,
            fontWeight: '600',
            color: theme.colors.textSecondary,
          }}
        >
          {saving ? 'Saving...' : 'Save'}
        </Text>
      </TouchableOpacity>
    </View>
  )
}