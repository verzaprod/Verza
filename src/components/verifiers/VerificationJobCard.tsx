import React from "react";
import { View, Text, TouchableOpacity, Image } from "react-native";
import { Ionicons } from "@expo/vector-icons";
import { useTheme } from "@/theme/ThemeProvider";

export default function VerificationJobCard({ job, onPress }) {
  const { isDark } = useTheme();

  return (
    <TouchableOpacity
      className={`mx-5 mb-4 rounded-3xl overflow-hidden ${
        isDark ? 'bg-[#1C1C1E] border border-gray-800' : 'bg-white'
      }`}
      style={{
        shadowColor: isDark ? '#000' : '#000',
        shadowOffset: { width: 0, height: 4 },
        shadowOpacity: isDark ? 0.3 : 0.08,
        shadowRadius: 12,
        elevation: 4,
      }}
      onPress={onPress}
      activeOpacity={0.7}
    >
      <View className="flex-row items-center p-5">
        {/* Avatar */}
        <View 
          className={`w-16 h-16 rounded-full items-center justify-center ${
            isDark ? 'bg-gray-700' : 'bg-gray-200'
          }`}
        >
          <Ionicons 
            name="person" 
            size={32} 
            color={isDark ? '#9CA3AF' : '#6B7280'} 
          />
        </View>

        {/* Info */}
        <View className="flex-1 ml-4">
          <Text 
            className={`text-xl font-semibold mb-1 ${
              isDark ? 'text-white' : 'text-gray-900'
            }`}
          >
            {job.requester}
          </Text>
          <Text 
            className={`text-base ${
              isDark ? 'text-gray-400' : 'text-gray-500'
            }`}
          >
            {job.doc}
          </Text>
        </View>

        {/* Document Icon */}
        <View 
          className={`w-12 h-12 rounded-xl items-center justify-center ${
            isDark ? 'bg-gray-800' : 'bg-gray-100'
          }`}
        >
          <Ionicons 
            name="card-outline" 
            size={24} 
            color={isDark ? '#9CA3AF' : '#6B7280'} 
          />
        </View>
      </View>
    </TouchableOpacity>
  );
}