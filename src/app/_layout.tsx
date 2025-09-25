import "../global.css";
import React, { useEffect } from "react";
import { Stack } from "expo-router/stack";
import { SafeAreaProvider } from "react-native-safe-area-context";
import { ThemeProvider } from "@/theme/ThemeProvider";
import { StatusBar } from "expo-status-bar";
import { useColorScheme } from "react-native";
import * as SplashScreen from "expo-splash-screen";
import { useFonts } from "expo-font";
import { ClerkProvider, useAuth } from "@clerk/clerk-expo";
import { tokenCache } from "@clerk/clerk-expo/token-cache";

SplashScreen.preventAutoHideAsync();

export default function Layout() {
  const colorScheme = useColorScheme();

  const [loaded] = useFonts({
    SansationLight: require("@/assets/fonts/Sansation-Light.ttf"),
    SFPro: require("@/assets/fonts/sf-pro-med.ttf"),
    UrbanistBold: require("@/assets/fonts/Urbanist-Bold.ttf"),
    UrbanistExtraBold: require("@/assets/fonts/Urbanist-ExtraBold.ttf"),
  });

  useEffect(() => {
    if (loaded) {
      SplashScreen.hide();
    }
  }, [loaded]);

  if (!loaded) {
    return null;
  }

  return (
    <ClerkProvider tokenCache={tokenCache}>
      <SafeAreaProvider>
        <ThemeProvider>
          <Stack screenOptions={{ headerShown: false }}>
            <Stack.Screen name="index" options={{ headerShown: false }} />
            <Stack.Screen
              name="(onboarding)"
              options={{ headerShown: false }}
            />
            <Stack.Screen name="(auth)" options={{ headerShown: false }} />
            <Stack.Screen name="(tabs)" options={{ headerShown: false }} />
            <Stack.Screen name="(kyc)" options={{ headerShown: false }} />
          </Stack>
          <StatusBar style={colorScheme === "dark" ? "light" : "dark"} />
        </ThemeProvider>
      </SafeAreaProvider>
    </ClerkProvider>
  );
}
