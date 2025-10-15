import React from "react";
import { View, Text, Modal, TouchableOpacity, Image, Pressable } from "react-native";
import { Ionicons } from "@expo/vector-icons";
import { useTheme } from "@/theme/ThemeProvider";

export default function DocumentReviewModal({ visible, job, onClose, onApprove, onReject }) {
  const theme = useTheme();

  if (!job) return null;

  return (
    <Modal
      visible={visible}
      transparent
      animationType="fade"
      onRequestClose={onClose}
      statusBarTranslucent
    >
      <Pressable 
        style={{
          flex: 1,
          backgroundColor: "rgba(0, 0, 0, 0.6)",
        }}
        onPress={onClose}
      >
        <Pressable 
          style={{
            flex: 1,
            justifyContent: "center",
            paddingHorizontal: 20,
          }}
          onPress={(e) => e.stopPropagation()}
        >
          <View 
            style={{
              borderRadius: theme.borderRadius.lg * 2,
              backgroundColor: theme.isDark 
                ? theme.colors.backgroundDark 
                : theme.colors.backgroundLight,
              overflow: "hidden",
              shadowColor: "#000",
              shadowOffset: { width: 0, height: 20 },
              shadowOpacity: 0.3,
              shadowRadius: 30,
              elevation: 20,
            }}
          >
            {/* Document Preview */}
            <View style={{ padding: theme.spacing.xl }}>
              <View 
                style={{
                  backgroundColor: theme.colors.primaryGreen,
                  borderRadius: theme.borderRadius.lg * 1.5,
                  padding: theme.spacing.lg,
                  overflow: "hidden",
                  aspectRatio: 1.6,
                  shadowColor: theme.colors.primaryGreen,
                  shadowOffset: { width: 0, height: 8 },
                  shadowOpacity: 0.3,
                  shadowRadius: 16,
                  elevation: 8,
                }}
              >
                <Image
                  source={{ uri: job.documentImage || "https://via.placeholder.com/600x375" }}
                  style={{
                    width: "100%",
                    height: "100%",
                    borderRadius: theme.borderRadius.md,
                  }}
                  resizeMode="contain"
                />
              </View>
            </View>

            {/* Action Buttons */}
            <View style={{
              flexDirection: "row",
              justifyContent: "center",
              gap: 40,
              paddingBottom: theme.spacing.xl,
            }}>
              {/* Reject Button */}
              <TouchableOpacity
                style={{
                  width: 80,
                  height: 80,
                  backgroundColor: theme.colors.error,
                  borderRadius: 40,
                  alignItems: "center",
                  justifyContent: "center",
                  shadowColor: theme.colors.error,
                  shadowOffset: { width: 0, height: 6 },
                  shadowOpacity: 0.4,
                  shadowRadius: 12,
                  elevation: 8,
                }}
                onPress={() => onReject(job)}
                activeOpacity={0.8}
              >
                <Ionicons name="close" size={40} color="#FFFFFF" />
              </TouchableOpacity>

              {/* Approve Button */}
              <TouchableOpacity
                style={{
                  width: 80,
                  height: 80,
                  backgroundColor: theme.colors.primaryGreen,
                  borderRadius: 40,
                  alignItems: "center",
                  justifyContent: "center",
                  shadowColor: theme.colors.primaryGreen,
                  shadowOffset: { width: 0, height: 6 },
                  shadowOpacity: 0.4,
                  shadowRadius: 12,
                  elevation: 8,
                }}
                onPress={() => onApprove(job)}
                activeOpacity={0.8}
              >
                <Ionicons name="checkmark" size={40} color="#FFFFFF" />
              </TouchableOpacity>
            </View>

            {/* Requester Info (Dimmed) */}
            <View style={{ 
              paddingHorizontal: theme.spacing.xl, 
              paddingBottom: theme.spacing.xl,
              opacity: 0.6 
            }}>
              <View 
                style={{
                  flexDirection: "row",
                  alignItems: "center",
                  borderRadius: theme.borderRadius.lg,
                  padding: theme.spacing.md,
                  backgroundColor: theme.isDark 
                    ? "rgba(255, 255, 255, 0.1)" 
                    : "rgba(255, 255, 255, 0.6)",
                }}
              >
                <View 
                  style={{
                    width: 48,
                    height: 48,
                    borderRadius: 24,
                    backgroundColor: theme.isDark 
                      ? `${theme.colors.textSecondary}30` 
                      : `${theme.colors.textSecondary}20`,
                    alignItems: "center",
                    justifyContent: "center",
                  }}
                >
                  <Ionicons 
                    name="person" 
                    size={24} 
                    color={theme.colors.textSecondary} 
                  />
                </View>
                
                <View style={{ flex: 1, marginLeft: theme.spacing.sm }}>
                  <Text 
                    style={{
                      fontSize: 18,
                      fontWeight: "600",
                      color: theme.colors.textPrimary,
                    }}
                  >
                    {job.requester}
                  </Text>
                  <Text 
                    style={{
                      fontSize: 14,
                      color: theme.colors.textSecondary,
                    }}
                  >
                    {job.doc}
                  </Text>
                </View>
                
                <View style={{ 
                  width: 40, 
                  height: 40, 
                  alignItems: "center", 
                  justifyContent: "center" 
                }}>
                  <Ionicons 
                    name="card-outline" 
                    size={22} 
                    color={theme.colors.textSecondary} 
                  />
                </View>
              </View>
            </View>
          </View>
        </Pressable>
      </Pressable>
    </Modal>
  );
}