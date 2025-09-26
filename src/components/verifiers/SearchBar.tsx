import React, { useState } from 'react'
import { View } from 'react-native'
import { InputBox } from '@/components/ui/InputBox'
import { useTheme } from '@/theme/ThemeProvider'

export const SearchBar: React.FC = () => {
  const [searchQuery, setSearchQuery] = useState('')
  const theme = useTheme()

  return (
    <View>
      <InputBox
        placeholder="Search verifiers..."
        value={searchQuery}
        onChangeText={setSearchQuery}
        style={{
          backgroundColor: theme.colors.backgroundLight,
        }}
      />
    </View>
  )
}