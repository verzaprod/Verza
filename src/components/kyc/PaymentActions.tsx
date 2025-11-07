import React, { useState } from "react";
import { View, Text, TouchableOpacity, ScrollView } from "react-native";
import { useTheme } from "@/theme/ThemeProvider";
import { Button } from "@/components/ui/Button";

type TokenType = "HBAR" | "ADA" | "NIGHT";

interface PaymentActionsProps {
  verifierDetails: { fee: number; currency: string };
  isProcessing: boolean;
  onConfirmPayment: (selectedToken: TokenType) => void;
  onCancel: () => void;
}

interface Token {
  symbol: TokenType;
  name: string;
  icon: string;
  color: string;
}

const AVAILABLE_TOKENS: Token[] = [
  {
    symbol: "HBAR",
    name: "Hedera",
    icon: "Ⓗ",
    color: "#4F46E5",
  },
  {
    symbol: "ADA",
    name: "Cardano",
    icon: "₳",
    color: "#0033AD",
  },
  {
    symbol: "NIGHT",
    name: "Night Token",
    icon: "🌙",
    color: "#7C3AED",
  },
];

export function PaymentActions({
  verifierDetails,
  isProcessing,
  onConfirmPayment,
  onCancel,
}: PaymentActionsProps) {
  const theme = useTheme();
  const [selectedToken, setSelectedToken] = useState<TokenType>("HBAR");

  return (
    <View style={{ paddingBottom: 40 }}>
      {/* Token Selection Section */}
      <View
        style={{
          backgroundColor: theme.colors.background,
          borderRadius: theme.borderRadius.lg,
          padding: 20,
          marginBottom: 24,
        }}
      >
        <Text
          style={{
            fontSize: 16,
            fontWeight: "600",
            color: theme.colors.textPrimary,
            marginBottom: 16,
          }}
        >
          Select Payment Token
        </Text>

        <ScrollView
          horizontal
          showsHorizontalScrollIndicator={false}
          contentContainerStyle={{ gap: 12 }}
        >
          {AVAILABLE_TOKENS.map((token) => (
            <TouchableOpacity
              key={token.symbol}
              onPress={() => setSelectedToken(token.symbol)}
              style={{
                paddingHorizontal: 20,
                paddingVertical: 16,
                borderRadius: theme.borderRadius.md,
                borderWidth: 2,
                borderColor:
                  selectedToken === token.symbol
                    ? token.color
                    : theme.colors.boxBorder,
                backgroundColor:
                  selectedToken === token.symbol
                    ? `${token.color}15`
                    : theme.colors.background,
                minWidth: 120,
                alignItems: "center",
              }}
            >
              <Text style={{ fontSize: 28, marginBottom: 8 }}>
                {token.icon}
              </Text>
              <Text
                style={{
                  fontSize: 16,
                  fontWeight: "600",
                  color:
                    selectedToken === token.symbol
                      ? token.color
                      : theme.colors.textPrimary,
                  marginBottom: 4,
                }}
              >
                {token.symbol}
              </Text>
              <Text
                style={{
                  fontSize: 12,
                  color: theme.colors.textSecondary,
                }}
              >
                {token.name}
              </Text>
            </TouchableOpacity>
          ))}
        </ScrollView>
       
      </View>

      {/* Payment Buttons */}
      <Button
        text={
          isProcessing
            ? "Processing..."
            : `Pay ${verifierDetails.fee} ${selectedToken}`
        }
        onPress={() => onConfirmPayment(selectedToken)}
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