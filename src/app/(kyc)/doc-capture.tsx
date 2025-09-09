import { useEffect, useState } from "react";
import { View, Text, ActivityIndicator, Alert } from "react-native";
import { useRouter, useLocalSearchParams } from "expo-router";
import { Onfido, OnfidoCaptureType, OnfidoDocumentType } from "@onfido/react-native-sdk";
import { useTheme } from "@/theme/ThemeProvider";

export default function DocCapture() {
  const router = useRouter();
  const theme = useTheme();
  const { type } = useLocalSearchParams();
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const startOnfido = async () => {
      try {
        // Fetch your SDK token from your backend
        const response = await fetch("https://your-backend.com/api/onfido-token");
        const { sdkToken } = await response.json();

        setLoading(false);

        const result = await Onfido.start({
          sdkToken,
          flowSteps: {
            welcome: true,
            captureDocument: {
              docType: type === "passport"
                ? OnfidoDocumentType.PASSPORT
                : type === "driver-license"
                ? OnfidoDocumentType.DRIVING_LICENCE
                : OnfidoDocumentType.NATIONAL_IDENTITY_CARD,
            },
            captureFace: {
              type: OnfidoCaptureType.PHOTO,}
          },
        });

        console.log("Onfido result:", result);
        
        router.replace("/(kyc)/selfie-note");
      } catch (error) {
        setLoading(false);
        Alert.alert("Verification Error", "Unable to start verification. Please try again.");
      }
    };

    startOnfido();
  }, [type, router]);

  if (loading) {
    return (
      <View style={{ flex: 1, justifyContent: "center", alignItems: "center", backgroundColor: theme.colors.background }}>
        <ActivityIndicator size="large" color={theme.colors.primaryGreen} />
        <Text style={{ marginTop: 16, color: theme.colors.textSecondary }}>Preparing verification...</Text>
      </View>
    );
  }

  return null;
}