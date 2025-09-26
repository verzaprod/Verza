import React from 'react'
import { View } from 'react-native'
import { useTheme } from '@/theme/ThemeProvider'
import { CredentialCard } from './CredentialCard'

interface Credential {
  id: string
  type: string
  status: 'verified' | 'pending'
  icon: string
}

interface CredentialsListProps {
  credentials: Credential[]
}

export const CredentialsList: React.FC<CredentialsListProps> = ({ credentials }) => {
  const theme = useTheme()

  return (
    <View style={{ gap: theme.spacing.md }}>
      {credentials.map((credential) => (
        <CredentialCard
          key={credential.id}
          type={credential.type}
          status={credential.status}
          icon={credential.icon}
        />
      ))}
    </View>
  )
}
