import React, { useState } from "react";
import { View, Text, TouchableOpacity, ScrollView, Image } from "react-native";
import { useTheme } from "@/theme/ThemeProvider";
import { Button } from "@/components/ui/Button";

type TokenType = "HBAR" | "ADA" | "NIGHT";

interface PaymentActionsProps {
  verifierDetails: { 
    fee: number; 
    currency: string;
    fees?: {
      HBAR: number;
      ADA: number;
      NIGHT: number;
    };
  };
  isProcessing: boolean;
  selectedToken: TokenType;
  onTokenChange: (token: TokenType) => void;
  onConfirmPayment: (selectedToken: TokenType) => void;
  onCancel: () => void;
}

interface Token {
  symbol: TokenType;
  name: string;
  iconSource: any;
  color: string;
}

const AVAILABLE_TOKENS: Token[] = [
  {
    symbol: "HBAR",
    name: "Hedera",
    iconSource: require("@/assets/images/shield.png"), // You can replace with actual HBAR logo
    color: "#4F46E5",
  },
  {
    symbol: "ADA",
    name: "Cardano",
    iconSource: require("@/assets/images/shield-check.png"), // You can replace with actual ADA logo
    color: "#0033AD",
  },
  {
    symbol: "NIGHT",
    name: "Night Token",
    iconSource: require("@/assets/images/wifi.png"), // You can replace with actual NIGHT logo
    color: "#7C3AED",
  },
];

export function PaymentActions({
  verifierDetails,
  isProcessing,
  selectedToken,
  onTokenChange,
  onConfirmPayment,
  onCancel,
}: PaymentActionsProps) {
  const theme = useTheme();

  // Get the fee for the selected token - now comes directly from verifierDetails
  const currentFee = verifierDetails.fee;
  const currentCurrency = verifierDetails.currency;

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
          {AVAILABLE_TOKENS.map((token) => {
            const tokenFee = verifierDetails.fees?.[token.symbol] || verifierDetails.fee;
            
            return (
              <TouchableOpacity
                key={token.symbol}
                onPress={() => onTokenChange(token.symbol)}
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
                <View
                  style={{
                    width: 40,
                    height: 40,
                    marginBottom: 8,
                    alignItems: "center",
                    justifyContent: "center",
                  }}
                >
                  <Image
                    source={token.iconSource}
                    style={{
                      width: 40,
                      height: 40,
                      tintColor:
                        selectedToken === token.symbol
                          ? token.color
                          : theme.colors.textSecondary,
                    }}
                    resizeMode="contain"
                  />
                </View>
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
                    marginBottom: 4,
                  }}
                >
                  {token.name}
                </Text>
                <Text
                  style={{
                    fontSize: 14,
                    fontWeight: "700",
                    color:
                      selectedToken === token.symbol
                        ? token.color
                        : theme.colors.textPrimary,
                  }}
                >
                  {tokenFee} {token.symbol}
                </Text>
              </TouchableOpacity>
            );
          })}
        </ScrollView>
       
      </View>

      {/* Payment Buttons */}
      <Button
        text={
          isProcessing
            ? "Processing..."
            : `Pay ${currentFee} ${currentCurrency}`
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