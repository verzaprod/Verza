import React from 'react';
import {
  Modal,
  View,
  Text,
  TouchableOpacity,
  Pressable,
  FlatList,
} from 'react-native';
import { useTheme } from '@/theme/ThemeProvider';
import { Icon } from '@/components/ui/Icon';
import Feather from '@expo/vector-icons/Feather';

const credentialTypes = [
  { 
    id: 'driver-license', 
    label: "Driver's\nLicense", 
    icon: 'driver-license',
    description: 'Government issued driving permit',
    color: '#3B82F6'
  },
  { 
    id: 'passport', 
    label: "Passport", 
    icon: 'passport',
    description: 'International travel document',
    color: '#10B981'
  },
  { 
    id: 'id-card', 
    label: "Student ID", 
    icon: 'id-card',
    description: 'Educational institution identification',
    color: '#F59E0B'
  },
];

interface AddCredentialModalProps {
  visible: boolean;
  onClose: () => void;
  onSelect: (credentialType: any) => void;
}

export function AddCredentialModal({
  visible,
  onClose,
  onSelect,
}: AddCredentialModalProps) {
  const theme = useTheme();

  const handleSelect = (credentialType) => {
    onSelect(credentialType);
    onClose();
  };

  const renderCredentialType = ({ item }) => (
    <TouchableOpacity
      style={{
        flexDirection: 'row',
        alignItems: 'center',
        padding: theme.spacing.lg,
        backgroundColor: theme.colors.background,
        marginBottom: theme.spacing.sm,
        borderRadius: theme.borderRadius.md,
        borderWidth: 1,
        borderColor: theme.colors.background,
      }}
      onPress={() => handleSelect(item)}
    >
      <View
        style={{
          width: 50,
          height: 50,
          borderRadius: 25,
          backgroundColor: item.color,
          alignItems: 'center',
          justifyContent: 'center',
          marginRight: theme.spacing.md,
        }}
      >
        <Icon name={item.icon} size={24} color="white" />
      </View>
      <View style={{ flex: 1 }}>
        <Text
          style={{
            fontSize: 16,
            color: theme.colors.textPrimary,
            fontWeight: '600',
            marginBottom: 4,
          }}
        >
          {item.label}
        </Text>
        <Text
          style={{
            fontSize: 12,
            color: theme.colors.textSecondary,
          }}
        >
          {item.description}
        </Text>
      </View>
      <Feather
        name="chevron-right"
        size={20}
        color={theme.colors.textSecondary}
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
            maxHeight: '70%',
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
                fontWeight: 'bold',
                color: theme.colors.textPrimary,
              }}
            >
              Add Credential
            </Text>
            <TouchableOpacity onPress={onClose}>
              <Feather name="x" size={24} color={theme.colors.textSecondary} />
            </TouchableOpacity>
          </View>

          <Text
            style={{
              fontSize: 14,
              color: theme.colors.textSecondary,
              marginBottom: theme.spacing.lg,
            }}
          >
            Select a credential type to add to your profile
          </Text>

          <FlatList
            data={credentialTypes}
            renderItem={renderCredentialType}
            keyExtractor={(item) => item.id}
            showsVerticalScrollIndicator={false}
          />
        </Pressable>
      </Pressable>
    </Modal>
  );
}
