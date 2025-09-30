import { useState, useEffect } from "react";
import { Alert, ToastAndroid } from "react-native";
import { useRouter } from "expo-router";
import { apiService } from "@/services/api/apiService";
import { useKYCStore } from "@/store/kycStore";
import { MOCK_DATA } from "@/services/api/mockData";

interface VerifierDetails {
  id: string;
  name: string;
  fee: number;
  currency: string;
  description: string;
  estimatedTime: string;
}

export function useVerifierDetails(verifierId: string) {
  const router = useRouter();
  const [isProcessing, setIsProcessing] = useState(false);
  const [verifierDetails, setVerifierDetails] = useState<VerifierDetails | null>(null);
  const [isLoading, setIsLoading] = useState(true);
  const { setEscrowInfo, setCurrentStep } = useKYCStore();

  useEffect(() => {
    const fetchVerifierDetails = () => {
      try {
        setIsLoading(true);
        
        const verifier = MOCK_DATA.verifiers.find(v => v.id === verifierId);
        
        if (verifier) {
          setVerifierDetails({
            id: verifier.id,
            name: verifier.name,
            fee: 25.0,
            currency: "HBAR",
            description: verifier.description,
            estimatedTime: "2-5 minutes",
          });
        } else {
          setVerifierDetails({
            id: verifierId,
            name: "TechCorp Solutions",
            fee: 25.0,
            currency: "HBAR",
            description: "Enterprise-grade identity verification",
            estimatedTime: "2-5 minutes",
          });
        }
      } catch (error) {
        console.error("Error fetching verifier details:", error);
        Alert.alert("Error", "Failed to load verifier details");
      } finally {
        setIsLoading(false);
      }
    };

    if (verifierId) {
      fetchVerifierDetails();
    }
  }, [verifierId]);

  const handleConfirmPayment = async () => {
    if (!verifierDetails) return;

    setIsProcessing(true);

    try {
      const escrowResponse = await apiService.initiateEscrow({
        verifier_id: verifierId,
        amount: verifierDetails.fee,
        currency: verifierDetails.currency,
        auto_release_hours: 24,
      });

      const escrowData = await escrowResponse.json();

      if (escrowResponse.ok) {
        setEscrowInfo(escrowData.escrow_id, verifierId);
        setCurrentStep("selection");
        ToastAndroid.show("Escrow initiated successfully!", ToastAndroid.SHORT);
        router.push(`/(kyc)/selection-type`);
      } else {
        Alert.alert("Payment Failed", escrowData.error || "Unable to process payment");
      }
    } catch (error) {
      Alert.alert("Error", "Network error occurred. Please try again.");
    } finally {
      setIsProcessing(false);
    }
  };

  return {
    verifierDetails,
    isLoading,
    isProcessing,
    handleConfirmPayment
  };
}