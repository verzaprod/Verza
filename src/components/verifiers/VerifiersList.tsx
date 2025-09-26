import React from 'react'
import { View } from 'react-native'
import { useTheme } from '@/theme/ThemeProvider'
import { VerifierCard } from './VerifierCard'

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

interface VerifiersListProps {
  verifiers: Verifier[]
}

export const VerifiersList: React.FC<VerifiersListProps> = ({ verifiers }) => {
  const theme = useTheme()

  return (
    <View style={{ gap: theme.spacing.md }}>
      {verifiers.map((verifier) => (
        <VerifierCard key={verifier.id} verifier={verifier} />
      ))}
    </View>
  )
}