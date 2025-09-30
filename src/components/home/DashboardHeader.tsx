import React from "react";
import { View, TouchableOpacity, ToastAndroid } from "react-native";
import { Icon } from "@/components/ui/Icon";
import { useRouter } from "expo-router";

export const DashboardHeader: React.FC = () => {
  const router = useRouter();

  const handleNotificationsClick = () => {
    ToastAndroid.show("No new notifications", ToastAndroid.SHORT);
  };

  return (
    <View
      className="flex-row justify-between items-center"
      style={{ paddingVertical: 16 }}
    >
      <TouchableOpacity
        style={{
          width: 48,
          height: 48,
          borderRadius: 24,
          backgroundColor: "#16A34A",
          alignItems: "center",
          justifyContent: "center",
        }}
        onPress={() => router.push("/profile")}
      >
        <Icon name="avatar" size={32} />
      </TouchableOpacity>

      <TouchableOpacity
        style={{
          position: "relative",
          width: 24,
          height: 24,
          borderRadius: 12,
          alignItems: "center",
          justifyContent: "center",
        }}
        onPress={handleNotificationsClick}
      >
        <Icon name="notification" size={24} />
      </TouchableOpacity>
    </View>
  );
};
