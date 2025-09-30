import React from 'react';
import { View, Text } from 'react-native';
import { useTheme } from '@/theme/ThemeProvider';
import FontAwesome5 from '@expo/vector-icons/FontAwesome5';

interface VerifiedAccount {
  id: string;
  name: string;
  status: string;
}

interface AccountHeaderProps {
  account: VerifiedAccount;
}

export function AccountHeader({ account }: AccountHeaderProps) {
  const theme = useTheme();

  const mockData = {
    verificationDate: 'September 25, 2025',
    lastActivity: '3 days ago',
  };

  return (
    <View style={{ marginBottom: theme.spacing.xl }}>
      <View style={{
        flexDirection: 'row',
        alignItems: 'center',
        marginBottom: theme.spacing.lg,
        padding: theme.spacing.lg,
        backgroundColor: theme.colors.background,
        borderRadius: theme.borderRadius.lg,
      }}>
        <View style={{
          width: 50,
          height: 50,
          backgroundColor: '#4CAF50',
          borderRadius: 10,
          alignItems: 'center',
          justifyContent: 'center',
          marginRight: theme.spacing.md,
        }}>
          <FontAwesome5 name="check" size={20} color="#fff" />
        </View>
        <View style={{ flex: 1 }}>
          <Text style={{
            fontSize: 18,
            fontWeight: '600',
            color: theme.colors.textPrimary,
            marginBottom: 4,
          }}>
            {account.name}
          </Text>
          <Text style={{
            fontSize: 14,
            color: '#4CAF50',
            fontWeight: '500',
          }}>
            Verified & Active
          </Text>
        </View>
      </View>

      <View style={{
        backgroundColor: theme.colors.background,
        borderRadius: theme.borderRadius.md,
        padding: theme.spacing.md,
      }}>
        <InfoRow 
          label="Verified" 
          value={mockData.verificationDate} 
        />
        <InfoRow 
          label="Last Activity" 
          value={mockData.lastActivity} 
          isLast 
        />
      </View>
    </View>
  );
}

function InfoRow({ label, value, isLast = false }) {
  const theme = useTheme();
  
  return (
    <View style={{
      flexDirection: 'row',
      justifyContent: 'space-between',
      paddingVertical: theme.spacing.sm,
      borderBottomWidth: isLast ? 0 : 1,
      borderBottomColor: theme.colors.boxBorder,
    }}>
      <Text style={{
        fontSize: 14,
        color: theme.colors.textSecondary,
      }}>
        {label}
      </Text>
      <Text style={{
        fontSize: 14,
        color: theme.colors.textPrimary,
        fontWeight: '500',
      }}>
        {value}
      </Text>
    </View>
  );
}
