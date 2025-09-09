import { CTAButton } from "@/components/ui/CTAButton";
import { Icon } from "@/components/ui/Icon";
import { useTheme } from "@/theme/ThemeProvider";
import { useRouter } from "expo-router";
import { View, Text, KeyboardAvoidingView, Platform } from "react-native";

export default function Success({
  redirectType,
  tagline,
  title,
  buttonText,
}: {
  redirectType: "kyc" | "auth";
  tagline: string;
  title: string;
  buttonText: string;
}) {
  const theme = useTheme();
  const router = useRouter();

  const handleNavigation = () => {
    if (redirectType === "kyc") {
      router.push("/(tabs)/home");
    } else {
      router.replace("/(kyc)/selection-type");
    }
  };

  return (
    <KeyboardAvoidingView
      style={{ flex: 1 }}
      behavior={Platform.OS === "ios" ? "padding" : "height"}
    >
      <View
        className="flex-1 items-center justify-center px-6"
        style={{
          backgroundColor: theme.colors.background,
        }}
      >
        <Icon name="success" size={100} />

        <Text
          className="text-3xl"
          style={{
            fontFamily: theme.fonts.welcomeHeading,
            color: theme.colors.textPrimary,
            marginVertical: 20,
          }}
        >
          {title}
        </Text>
        <Text
          className="text-center text-lg mb-10 px-4"
          style={{
            color: theme.colors.textSecondary,
            fontFamily: theme.fonts.onboardingTagline,
          }}
        >
          {tagline}
        </Text>

        <CTAButton title={buttonText} onPress={handleNavigation} />
      </View>
    </KeyboardAvoidingView>
  );
}
