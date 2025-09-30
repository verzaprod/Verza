import React from "react";
import { View } from "react-native";
import { Button } from "@/components/ui/Button";

interface PaymentActionsProps {
  verifierDetails: { fee: number; currency: string };
  isProcessing: boolean;
  onConfirmPayment: () => void;
  onCancel: () => void;
}

export function PaymentActions({
  verifierDetails,
  isProcessing,
  onConfirmPayment,
  onCancel,
}: PaymentActionsProps) {
  return (
    <View style={{ paddingBottom: 40 }}>
      <Button
        text={
          isProcessing
            ? "Processing..."
            : `Pay ${verifierDetails.fee} ${verifierDetails.currency}`
        }
        onPress={onConfirmPayment}
        disabled={isProcessing}
        style={{ marginBottom: 16 }}
      />

      <Button
        text="Cancel"
        variant="secondary"
        onPress={onCancel}
        disabled={isProcessing}
      />
    </View>
  );
}