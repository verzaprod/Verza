import React from "react";
import { View, Text, FlatList } from "react-native";
import VerificationJobCard from "./VerificationJobCard";

export default function VerificationJobList({ jobs, onJobPress }) {
  const renderHeader = () => (
    <View className="flex-row justify-between items-center mt-5 mb-6">
      <Text className="text-4xl font-bold text-gray-900">
        Due Tasks
      </Text>
      <View className="w-12 h-12 bg-white rounded-xl justify-center items-center shadow-sm">
        <Text className="text-2xl">📋</Text>
      </View>
    </View>
  );

  return (
    <FlatList
      data={jobs}
      keyExtractor={(item) => item.id}
      renderItem={({ item }) => (
        <VerificationJobCard job={item} onPress={() => onJobPress(item)} />
      )}
      ListHeaderComponent={renderHeader}
      contentContainerClassName="pb-24"
      showsVerticalScrollIndicator={false}
    />
  );
}