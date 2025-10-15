import React from "react";
import { View, Text, Modal, TouchableOpacity, Image } from "react-native";

export default function DocumentReviewModal({ visible, job, onClose, onApprove, onReject }) {
  if (!job) return null;

  return (
    <Modal
      visible={visible}
      transparent
      animationType="fade"
      onRequestClose={onClose}
    >
      <View className="flex-1 bg-black/50 justify-center p-5">
        <View className="bg-gray-50 rounded-3xl p-6">
          {/* Document Preview */}
          <View className="mb-8">
            <View className="bg-green-500 rounded-3xl p-5 aspect-video justify-center items-center shadow-xl">
              <Image
                source={{ uri: job.documentImage || "https://via.placeholder.com/400x250" }}
                className="w-full h-full rounded-2xl"
                resizeMode="contain"
              />
            </View>
          </View>

          {/* Action Buttons */}
          <View className="flex-row justify-center gap-8 mb-6">
            {/* Reject Button */}
            <TouchableOpacity
              className="w-18 h-18 bg-red-500 rounded-full justify-center items-center shadow-lg active:opacity-80"
              onPress={() => onReject(job)}
              activeOpacity={0.8}
            >
              <Text className="text-4xl text-white font-bold">✕</Text>
            </TouchableOpacity>

            {/* Approve Button */}
            <TouchableOpacity
              className="w-18 h-18 bg-green-500 rounded-full justify-center items-center shadow-lg active:opacity-80"
              onPress={() => onApprove(job)}
              activeOpacity={0.8}
            >
              <Text className="text-4xl text-white font-bold">✓</Text>
            </TouchableOpacity>
          </View>

          {/* Requester Info (Dimmed) */}
          <View className="opacity-60">
            <View className="flex-row items-center bg-white/60 rounded-2xl p-4">
              <View className="w-11 h-11 bg-gray-200 rounded-full justify-center items-center">
                <Text className="text-2xl">👤</Text>
              </View>
              
              <View className="flex-1 ml-3">
                <Text className="text-lg font-semibold text-gray-900">
                  {job.requester}
                </Text>
                <Text className="text-sm text-gray-600">
                  {job.doc}
                </Text>
              </View>
              
              <View className="w-9 h-9 justify-center items-center">
                <Text className="text-xl">🪪</Text>
              </View>
            </View>
          </View>
        </View>
      </View>
    </Modal>
  );
}