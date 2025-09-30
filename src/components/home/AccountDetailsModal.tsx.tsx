import React from 'react';
import {
  Modal,
  View,
  Text,
  TouchableOpacity,
  Pressable,
} from 'react-native';
import { useTheme } from '@/theme/ThemeProvider';
import { Icon } from '@/components/ui/Icon';
import { AccountHeader } from './AccountHeader';
import { AccountActions } from './AccountActions';
import Feather from '@expo/vector-icons/Feather';

interface VerifiedAccount {
  id: string;
  name: string;
  status: string;
}

interface AccountDetailsModalProps {
  visible: boolean;
  account: VerifiedAccount | null;
  onClose: () => void;
  onDisconnect: (accountId: string) => void;
}

export function AccountDetailsModal({
  visible,
  account,
  onClose,
  onDisconnect,
}: AccountDetailsModalProps) {
  const theme = useTheme();

  if (!account) return null;

  const handleDisconnect = () => {
    onDisconnect(account.id);
    onClose();
  };

  return (
    <Modal
      visible={visible}
      animationType="slide"
      transparent
      onRequestClose={onClose}
    >
      <Pressable
        style={{
          flex: 1,
          backgroundColor: 'rgba(0, 0, 0, 0.5)',
          justifyContent: 'center',
          alignItems: 'center',
        }}
        onPress={onClose}
      >
        <Pressable
          style={{
            backgroundColor: theme.colors.background,
            borderRadius: theme.borderRadius.md,
            padding: theme.spacing.md,
            width: '90%',
            maxWidth: 400,
            shadowColor: theme.isDark ? "#fff" : "#000",
            shadowOffset: { width: 0, height: 10 },
            shadowOpacity: 0.25,
            shadowRadius: 10,
            elevation: 10,
          }}
          onPress={(e) => e.stopPropagation()}
        >
\          <View style={{
            flexDirection: 'row',
            justifyContent: 'space-between',
            alignItems: 'center',
            marginBottom: theme.spacing.xl,
          }}>
            <Text style={{
              fontSize: 24,
              fontWeight: 'bold',
              color: theme.colors.textPrimary,
            }}>
              Account Details
            </Text>
            <TouchableOpacity onPress={onClose}>
              <Feather name="x" size={24} color={theme.colors.textSecondary} />
            </TouchableOpacity>
          </View>

          <AccountHeader account={account} />
          <AccountActions onDisconnect={handleDisconnect} />
        </Pressable>
      </Pressable>
    </Modal>
  );
}
