import { CTAButton } from "@/components/ui/CTAButton";
import { Icon } from "@/components/ui/Icon";
import { useTheme } from "@/theme/ThemeProvider";
import { useRouter } from "expo-router";
import { View, Text, Image, KeyboardAvoidingView, Platform, } from "react-native";

export default function Success() {
  
  const theme = useTheme();
  const router = useRouter();

  const handleStartKYC = () => {
    router.replace("/(kyc)/");
  }

  return (
    <KeyboardAvoidingView style={{ flex: 1 }} behavior={Platform.OS === "ios" ? "padding" : "height"}>
      <View className="flex-1 items-center justify-center px-6"
        style={{ 
          backgroundColor: theme.colors.background,
        }}
      >
        <Icon name="success" size={100}/>

        <Text
          className="text-3xl"
          style={{
            fontFamily: theme.fonts.welcomeHeading,
            color: theme.colors.textPrimary,
            marginVertical: 20,
          }}
        >
          Success!
        </Text>
        <Text
          className="text-center text-lg mb-10 px-4"
          style={{
            color: theme.colors.textSecondary,
            fontFamily: theme.fonts.onboardingTagline,
          }}
        >
          Your wallet has been created successfully. Let’s verify your identity to unlock all features.
        </Text>

        <CTAButton 
          title="Start KYC Verification"
          onPress={handleStartKYC}
        />
      </View>
    </KeyboardAvoidingView>
  )
}
