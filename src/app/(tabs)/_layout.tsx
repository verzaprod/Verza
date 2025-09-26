import { Tabs, Stack } from "expo-router";
import { Icon } from "@/components/ui/Icon";
import { useTheme } from "@/theme/ThemeProvider";

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
          tabBarIconStyle: { marginTop: 4 },
          tabBarStyle: {
            height: 60,
          }
        }}
      >
        <Tabs.Screen
          name="home"
          options={{
            tabBarLabel: "Home",
            tabBarIcon: ({ color, size, }) => (
              <Icon name="home" size={size} color={color} />
            ),
            headerShown: false,
          }}
        />
        <Tabs.Screen
          name="profile"
          options={{
            tabBarLabel: "Profile",
            tabBarIcon: ({ color, size }) => (
              <Icon name="profile" size={size} color={color} />
            ),
            headerShown: false,
          }}
        />
      </Tabs>
    </>
  );
}
