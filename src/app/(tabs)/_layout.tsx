import { Tabs, Stack } from "expo-router";
import { Icon } from "@/components/ui/Icon";
import { useTheme } from "@/theme/ThemeProvider";
import Feather from "@expo/vector-icons/Feather";
import FontAwesome5 from "@expo/vector-icons/FontAwesome5";
import Entypo from "@expo/vector-icons/Entypo";

export default function TabsLayout() {
  const theme = useTheme();

  return (
    <>
      <Tabs
        screenOptions={{
          headerShown: false,
          animation: "shift",
          // tabBarShowLabel: false,
          tabBarActiveTintColor: theme.colors.primaryGreen,
          tabBarStyle: {
            height: 60,
            paddingTop: 8,
            backgroundColor: theme.colors.background,
          },
        }}
      >
        <Tabs.Screen
          name="home"
          options={{
            tabBarLabel: "Home",
            tabBarIcon: ({ color, size }) => (
              <Feather name="home" size={size} color={color} />
            ),
          }}
        />
        <Tabs.Screen
          name="verifiers"
          options={{
            tabBarLabel: "Verifiers",
            tabBarIcon: ({ color, size }) => (
              <Entypo name="shop" size={size} color={color} />
            ),
          }}
        />
        <Tabs.Screen
          name="profile"
          options={{
            tabBarLabel: "Profile",
            tabBarIcon: ({ color, size }) => (
              <FontAwesome5 name="user-check" size={size} color={color} />
            ),
          }}
        />
      </Tabs>
    </>
  );
}
