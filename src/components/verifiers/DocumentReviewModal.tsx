import React from "react";
import { View, Text, Modal, TouchableOpacity, Image, Pressable } from "react-native";
// import { BlurView } from "expo-blur";
import { Ionicons } from "@expo/vector-icons";
import { useTheme } from "@/theme/ThemeProvider";

export default function DocumentReviewModal({ visible, job, onClose, onApprove, onReject }) {
  const { isDark } = useTheme();

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
        className="flex-1 bg-black/60" 
        onPress={onClose}
      >
        <Pressable 
          className="flex-1 justify-center px-5"
          onPress={(e) => e.stopPropagation()}
        >
          <View 
            className={`rounded-[32px] overflow-hidden ${
              isDark ? 'bg-[#1C1C1E]' : 'bg-gray-50'
            }`}
            style={{
              shadowColor: '#000',
              shadowOffset: { width: 0, height: 20 },
              shadowOpacity: 0.3,
              shadowRadius: 30,
              elevation: 20,
            }}
          >
            {/* Document Preview */}
            <View className="p-6 pb-8">
              <View 
                className="bg-green-500 rounded-3xl p-5 overflow-hidden"
                style={{
                  aspectRatio: 1.6,
                  shadowColor: '#22C55E',
                  shadowOffset: { width: 0, height: 8 },
                  shadowOpacity: 0.3,
                  shadowRadius: 16,
                  elevation: 8,
                }}
              >
                <Image
                  source={{ uri: job.documentImage || "https://via.placeholder.com/600x375" }}
                  className="w-full h-full rounded-2xl"
                  resizeMode="contain"
                />
              </View>
            </View>

            {/* Action Buttons */}
            <View className="flex-row justify-center gap-10 pb-6">
              {/* Reject Button */}
              <TouchableOpacity
                className="w-20 h-20 bg-red-500 rounded-full items-center justify-center active:scale-95"
                style={{
                  shadowColor: '#EF4444',
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
                className="w-20 h-20 bg-green-500 rounded-full items-center justify-center active:scale-95"
                style={{
                  shadowColor: '#22C55E',
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
            <View className="px-6 pb-6 opacity-60">
              <View 
                className={`flex-row items-center rounded-2xl p-4 ${
                  isDark ? 'bg-white/10' : 'bg-white/60'
                }`}
              >
                <View 
                  className={`w-12 h-12 rounded-full items-center justify-center ${
                    isDark ? 'bg-gray-700' : 'bg-gray-200'
                  }`}
                >
                  <Ionicons 
                    name="person" 
                    size={24} 
                    color={isDark ? '#9CA3AF' : '#6B7280'} 
                  />
                </View>
                
                <View className="flex-1 ml-3">
                  <Text 
                    className={`text-lg font-semibold ${
                      isDark ? 'text-white' : 'text-gray-900'
                    }`}
                  >
                    {job.requester}
                  </Text>
                  <Text 
                    className={`text-sm ${
                      isDark ? 'text-gray-400' : 'text-gray-600'
                    }`}
                  >
                    {job.doc}
                  </Text>
                </View>
                
                <View className="w-10 h-10 items-center justify-center">
                  <Ionicons 
                    name="card-outline" 
                    size={22} 
                    color={isDark ? '#9CA3AF' : '#6B7280'} 
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