import React from "react";
import { Alert, TouchableOpacity, View } from "react-native";
import { useTheme } from "@/theme/ThemeProvider";
import { CredentialCard } from "./CredentialCard";
import Feather from "@expo/vector-icons/Feather";

interface Credential {
  id: string;
  type: string;
  status: "verified" | "pending";
  icon: string;
}

interface CredentialsListProps {
  credentials: any[];
  onRemoveCredential?: (credentialId: string) => void;
}

export const CredentialsList = ({
  credentials,
  onRemoveCredential,
}: CredentialsListProps) => {
  const theme = useTheme();
  
  const handleRemove = (credentialId: string, credentialType: string) => {
    Alert.alert(
      "Remove Credential",
      `Are you sure you want to remove ${credentialType}?`,
      [
        { text: "Cancel", style: "cancel" },
        { 
          text: "Remove", 
          style: "destructive",
          onPress: () => onRemoveCredential?.(credentialId)
        }
      ]
    );
  };

  return (
    <View style={{ gap: theme.spacing.md }}>
      {credentials.map((credential) => {
        return (
          <View key={credential.id} style={{ position: 'relative' }}>
            <CredentialCard
              type={credential.type}
              status={credential.status}
              icon={credential.icon}
            />
            
            {credential.status === "pending" && (
              <TouchableOpacity
                style={{
                  position: 'absolute',
                  top: 8,
                  right: 8,
                  backgroundColor: 'rgba(255, 255, 255, 0.9)',
                  borderRadius: 12,
                  padding: 4,
                }}
                onPress={() => handleRemove(credential.id, credential.type)}
              >
                <Feather name="x" size={16} color="#666" />
              </TouchableOpacity>
            )}
          </View>
        );
      })}
    </View>
  );
};
