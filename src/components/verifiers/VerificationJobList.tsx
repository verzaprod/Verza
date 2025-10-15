import React from "react";
import { View, Text, FlatList, TouchableOpacity } from "react-native";
import { Ionicons } from "@expo/vector-icons";
import { useTheme } from "@/theme/ThemeProvider";
import VerificationJobCard from "./VerificationJobCard";

export default function VerificationJobList({ jobs, onJobPress }) {
  const theme = useTheme();

  const renderHeader = () => (
    <View className="flex-row justify-between items-center mb-6 px-5">
      <Text 
        style={{
          fontSize: 36,
          fontWeight: "700",
          color: theme.colors.textPrimary,
          fontFamily: theme.fonts.welcomeHeading,
        }}
      >
        Due Tasks
      </Text>
      <TouchableOpacity 
        style={{
          width: 56,
          height: 56,
          borderRadius: theme.borderRadius.lg,
          backgroundColor: theme.colors.background,
          alignItems: "center",
          justifyContent: "center",
          shadowColor: theme.isDark ? "#fff" : "#000",
          shadowOffset: { width: 0, height: 2 },
          shadowOpacity: theme.isDark ? 0.3 : 0.08,
          shadowRadius: 8,
          elevation: 3,
        }}
      >
        <Ionicons 
          name="clipboard-outline" 
          size={24} 
          color={theme.colors.textSecondary} 
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
      contentContainerStyle={{ paddingBottom: 120 }}
      className="flex-1"
      showsVerticalScrollIndicator={false}
    />
  );
}