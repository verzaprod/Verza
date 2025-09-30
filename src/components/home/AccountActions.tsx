import React from 'react';
import { View, Text, TouchableOpacity, Alert } from 'react-native';
import { useTheme } from '@/theme/ThemeProvider';
import { Icon } from '@/components/ui/Icon';
import Feather from '@expo/vector-icons/Feather';

interface AccountActionsProps {
  onDisconnect: () => void;
}

type FeatherIconName = React.ComponentProps<typeof Feather>['name'];

export function AccountActions({ onDisconnect }: AccountActionsProps) {
  const theme = useTheme();

  const handleRefresh = () => {
    Alert.alert('Success', 'Account data refreshed successfully!');
  };

  const handleExport = () => {
    Alert.alert('Export', 'Account data will be exported to your downloads folder.');
  };

  const handleDisconnectConfirm = () => {
    Alert.alert(
      'Disconnect Account',
      'Are you sure you want to disconnect this account? This action cannot be undone.',
      [
        { text: 'Cancel', style: 'cancel' },
        { 
          text: 'Disconnect', 
          style: 'destructive',
          onPress: onDisconnect 
        }
      ]
    );
  }
  
  const actions: Array<{
    icon: FeatherIconName;
    label: string;
    onPress: () => void;
    color: string;
    background: string;
  }> = [
    {
      icon: 'refresh-cw',
      label: 'Refresh Data',
      onPress: handleRefresh,
      color: theme.colors.primaryGreen,
      background: theme.colors.background,
    },
    {
      icon: 'download',
      label: 'Export Data',
      onPress: handleExport,
      color: theme.colors.primaryGreen,
      background: theme.colors.background,
    },
    {
      icon: 'delete',
      label: 'Disconnect',
      onPress: handleDisconnectConfirm,
      color: '#F44336',
      background: 'rgba(244, 67, 54, 0.1)',
    },
  ];

  return (
    <View>
      <Text style={{
        fontSize: 16,
        fontWeight: '600',
        color: theme.colors.textPrimary,
        marginBottom: theme.spacing.md,
      }}>
        Actions
      </Text>
      
      {actions.map((action, index) => (
        <TouchableOpacity
          key={index}
          style={{
            flexDirection: 'row',
            alignItems: 'center',
            padding: theme.spacing.md,
            backgroundColor: action.background,
            borderRadius: theme.borderRadius.md,
            marginBottom: theme.spacing.sm,
          }}
          onPress={action.onPress}
        >
          <Feather 
            name={action.icon} 
            size={18} 
            color={action.color}
            style={{ marginRight: theme.spacing.md }}
          />
          <Text style={{
            fontSize: 16,
            color: action.color,
            fontWeight: '500',
          }}>
            {action.label}
          </Text>
        </TouchableOpacity>
      ))}
    </View>
  );
}
