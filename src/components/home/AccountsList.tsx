import React from "react";
import { View, Text, TouchableOpacity, Alert } from "react-native";
import { useTheme } from "@/theme/ThemeProvider";
import { Icon } from "@/components/ui/Icon";
import FontAwesome5 from "@expo/vector-icons/FontAwesome5";

interface VerifiedAccount {
  id: string;
  name: string;
  status: string;
}

interface OverlappingAccountsListProps {
  accounts: VerifiedAccount[];
  onRemoveAccount: (accountId: string) => void;
}

export const AccountsList: React.FC<OverlappingAccountsListProps> = ({
  accounts,
  onRemoveAccount,
}) => {
  const theme = useTheme();

  const renderAccountContent = (account: VerifiedAccount) => {
    if (account.status === "pending") {
      return {
        iconName: "clock",
        iconBgColor: "#FF9800",
        displayName: account.name,
        statusText: "Pending",
        statusColor: "#FF9800",
        actionText: "Cancel",
        actionColor: theme.colors.textSecondary,
      };
    }

    // Default verified state
    return {
      iconName: "check",
      iconBgColor: "#4CAF50",
      displayName: account.name,
      statusText: "Verified",
      statusColor: theme.colors.textSecondary,
      actionText: "View Details",
      actionColor: theme.colors.primaryGreen,
    };
  };

  const handleRemoveAccount = (accountId: string, accountName: string) => {
    Alert.alert(
      "Remove Account",
      `Are you sure you want to remove ${accountName}?`,
      [
        {
          text: "Cancel",
          style: "cancel",
        },
        {
          text: "Remove",
          style: "destructive",
          onPress: () => {
            onRemoveAccount(accountId);
          },
        },
      ]
    );
  };

  return (
    <View>
      {accounts.map((account, index) => {
        const isLast = index === accounts.length - 1;
        const zIndex = accounts.length + index;
        const opacity = 1 - index * 0.05;
        // const scale = 1 - index * 0.02;
        const marginTop = index === 0 ? 0 : -42;

        const accountContent = renderAccountContent(account);

        return (
          <TouchableOpacity
            key={account.id}
            style={{
              backgroundColor: theme.colors.background,
              borderRadius: 32,
              padding: 24,
              flexDirection: "row",
              alignItems: "center",
              justifyContent: "space-between",
              zIndex,
              opacity,
              marginTop,
              shadowColor: theme.isDark ? "#fff" : "#000",
              shadowOffset: { width: 0, height: 0 },
              shadowOpacity: 0.15,
              shadowRadius: 8,
              elevation: 8,
            }}
          >
            <View
              style={{
                width: 40,
                height: 40,
                backgroundColor: accountContent.iconBgColor,
                borderRadius: 8,
                alignItems: "center",
                justifyContent: "center",
                marginRight: 16,
                shadowColor: theme.isDark ? "#fff" : "#000",
                shadowOffset: { width: 0, height: 2 },
                shadowOpacity: 0.1,
                shadowRadius: 4,
                elevation: 10,
              }}
            >
              <FontAwesome5
                name={accountContent.iconName}
                size={20}
                color="#fff"
              />
            </View>

            <View className="flex-1 flex-col ite ms-end">
              <View className="flex-1 justify-between flex-row">
                <Text
                  style={{
                    fontSize: 18,
                    fontWeight: "600",
                    color: theme.colors.textPrimary,
                    marginBottom: 4,
                  }}
                >
                  {accountContent.displayName}
                </Text>

                <TouchableOpacity
                  style={{
                    width: 32,
                    height: 32,
                    borderRadius: 16,
                    alignItems: "center",
                    justifyContent: "center",
                  }}
                  onPress={() => handleRemoveAccount(account.id, account.name)}
                  hitSlop={{ top: 8, bottom: 8, left: 8, right: 8 }}
                >
                  <Icon name="remove" size={16} />
                </TouchableOpacity>
              </View>

              {
                <View className="flex-row justify-between items-center mt-4">
                  <Text
                    style={{
                      fontSize: 14,
                      color: accountContent.statusColor,
                      fontStyle: "italic",
                    }}
                  >
                    {accountContent.statusText}
                  </Text>

                  <TouchableOpacity>
                    <Text
                      style={{
                        fontSize: 14,
                        fontWeight: "600",
                        color: theme.colors.primaryGreen,
                      }}
                    >
                      {accountContent.actionText}
                    </Text>
                  </TouchableOpacity>
                </View>
              }
            </View>
          </TouchableOpacity>
        );
      })}
    </View>
  );
};
