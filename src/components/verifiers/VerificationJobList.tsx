import React from "react";
import { View, Text, FlatList, TouchableOpacity } from "react-native";
import { Ionicons } from "@expo/vector-icons";
import { useTheme } from "@/theme/ThemeProvider";
import VerificationJobCard from "./VerificationJobCard";

export default function VerificationJobList({ jobs, onJobPress }) {
  const { isDark } = useTheme();

  const renderHeader = () => (
    <View className="flex-row justify-between items-center mb-6 px-5">
      <Text 
        className={`text-[36px] font-bold ${
          isDark ? 'text-white' : 'text-gray-900'
        }`}
      >
        Due Tasks
      </Text>
      <TouchableOpacity 
        className={`w-14 h-14 rounded-2xl items-center justify-center ${
          isDark ? 'bg-[#1C1C1E]' : 'bg-white'
        }`}
        style={{
          shadowColor: isDark ? '#000' : '#000',
          shadowOffset: { width: 0, height: 2 },
          shadowOpacity: isDark ? 0.3 : 0.08,
          shadowRadius: 8,
          elevation: 3,
        }}
      >
        <Ionicons 
          name="clipboard-outline" 
          size={24} 
          color={isDark ? '#9CA3AF' : '#6B7280'} 
        />
      </TouchableOpacity>
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
      contentContainerClassName="pb-32"
      className="flex-1"
      showsVerticalScrollIndicator={false}
    />
  );
}