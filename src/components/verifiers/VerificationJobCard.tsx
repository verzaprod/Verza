import React from "react";
import { View, Text, TouchableOpacity } from "react-native";

export default function VerificationJobCard({ job, onPress }) {
  return (
    <TouchableOpacity
      className="bg-white rounded-3xl mb-4 shadow-md active:opacity-70"
      onPress={onPress}
      activeOpacity={0.7}
    >
      <View className="flex-row items-center p-5">
        {/* Avatar */}
        <View className="w-14 h-14 bg-gray-200 rounded-full justify-center items-center">
          <Text className="text-3xl">👤</Text>
        </View>

        {/* Info */}
        <View className="flex-1 ml-4">
          <Text className="text-xl font-semibold text-gray-900 mb-1">
            {job.requester}
          </Text>
          <Text className="text-base text-gray-500">
            {job.doc}
          </Text>
        </View>

        {/* Document Icon */}
        <View className="w-12 h-12 bg-gray-100 rounded-xl justify-center items-center">
          <Text className="text-2xl">🪪</Text>
        </View>
      </View>
    </TouchableOpacity>
  );
}