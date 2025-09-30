import React from "react";
import { View } from "react-native";
import { Button } from "@/components/ui/Button";

interface ResultsActionsProps {
  status: "verified" | "rejected" | "pending_review";
  onNavigate: (route: string) => void;
}

export function ResultsActions({ status, onNavigate }: ResultsActionsProps) {
  const getActionConfig = () => {
    switch (status) {
      case "verified":
        return {
          primaryText: "Go to Profile",
          primaryRoute: "/(tabs)/profile",
        };
      case "rejected":
        return {
          primaryText: "Try Again",
          primaryRoute: "/(tabs)/verifiers",
        };
      default:
        return {
          primaryText: "Check Status",
          primaryRoute: "/(tabs)/home",
        };
    }
  };

  const config = getActionConfig();

  return (
    <View>
      <Button
        text={config.primaryText}
        onPress={() => onNavigate(config.primaryRoute)}
        style={{ marginBottom: 16 }}
      />

      <Button
        text="Back to Home"
        variant="secondary"
        onPress={() => onNavigate("/(tabs)/home")}
      />
    </View>
  );
}
