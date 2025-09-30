import React from 'react';
import {
  Modal,
  View,
  Text,
  TouchableOpacity,
  FlatList,
  Pressable,
} from 'react-native';
import { useTheme } from '@/theme/ThemeProvider';
import { Icon } from '@/components/ui/Icon';
import Feather from '@expo/vector-icons/Feather';

const integrations = [
  { id: '1', name: 'Google', icon: 'google', color: '#4285F4' },
  { id: '2', name: 'LinkedIn', icon: 'linkedin', color: '#0077B5' },
  { id: '3', name: 'Facebook', icon: 'facebook', color: '#1877F2' },
  { id: '4', name: 'Twitter', icon: 'twitter', color: '#1DA1F2' },
  { id: '5', name: 'Instagram', icon: 'instagram', color: '#E4405F' },
  { id: '6', name: 'GitHub', icon: 'github', color: '#333' },
  { id: '7', name: 'Discord', icon: 'discord', color: '#5865F2' },
  { id: '8', name: 'Slack', icon: 'slack', color: '#4A154B' },
];

interface AddAccountModalProps {
  visible: boolean;
  onClose: () => void;
  onSelectIntegration: (integration: any) => void;
}

export function AddAccountModal({
  visible,
  onClose,
  onSelectIntegration,
}: AddAccountModalProps) {
  const theme = useTheme();

  const renderIntegration = ({ item }) => (
    <TouchableOpacity
      style={{
        flexDirection: 'row',
        alignItems: 'center',
        padding: theme.spacing.lg,
        backgroundColor: theme.colors.background,
        marginBottom: theme.spacing.sm,
        borderRadius: theme.borderRadius.md,
        borderWidth: 1,
        borderColor: theme.colors.boxBorder,
      }}
      onPress={() => onSelectIntegration(item)}
    >
      <View
        style={{
          width: 40,
          height: 40,
          borderRadius: 20,
          backgroundColor: item.color,
          alignItems: 'center',
          justifyContent: 'center',
          marginRight: theme.spacing.md,
        }}
      >
        <Text style={{ color: 'white', fontWeight: 'bold', fontSize: 16 }}>
          {item.name.charAt(0)}
        </Text>
      </View>
      <Text
        style={{
          fontSize: 16,
          color: theme.colors.textPrimary,
          fontFamily: theme.fonts.onboardingTagline,
          flex: 1,
        }}
      >
        {item.name}
      </Text>
      <Feather
        name="chevron-right"
        size={20}
      />
    </TouchableOpacity>
  );

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
          justifyContent: 'flex-end',
        }}
        onPress={onClose}
      >
        <Pressable
          style={{
            backgroundColor: theme.colors.background,
            borderTopLeftRadius: theme.borderRadius.lg,
            borderTopRightRadius: theme.borderRadius.lg,
            paddingTop: theme.spacing.lg,
            paddingHorizontal: theme.spacing.lg,
            paddingBottom: theme.spacing.xl,
            maxHeight: '80%',
          }}
          onPress={(e) => e.stopPropagation()}
        >
          <View
            style={{
              flexDirection: 'row',
              justifyContent: 'space-between',
              alignItems: 'center',
              marginBottom: theme.spacing.lg,
            }}
          >
            <Text
              style={{
                fontSize: 20,
                fontFamily: theme.fonts.onboardingHeading,
                color: theme.colors.textPrimary,
              }}
            >
              Add Account
            </Text>
            <TouchableOpacity onPress={onClose}>
              <Icon name="remove" size={24} color={theme.colors.textSecondary} />
            </TouchableOpacity>
          </View>

          <Text
            style={{
              fontSize: 14,
              color: theme.colors.textSecondary,
              marginBottom: theme.spacing.lg,
              fontFamily: theme.fonts.onboardingTagline,
            }}
          >
            Select a platform to connect your account
          </Text>

          <FlatList
            data={integrations}
            renderItem={renderIntegration}
            keyExtractor={(item) => item.id}
            showsVerticalScrollIndicator={false}
          />
        </Pressable>
      </Pressable>
    </Modal>
  );
}
