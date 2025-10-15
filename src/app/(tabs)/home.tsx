import {
  View,
  Text,
  SafeAreaView,
  ScrollView,
  TouchableOpacity,
} from "react-native";
import { useTheme } from "@/theme/ThemeProvider";
import { useSafeAreaInsets } from "react-native-safe-area-context";
import { DashboardHeader } from "@/components/home/DashboardHeader";
import { AccountsList } from "@/components/home/AccountsList";
import { AddAccountButton } from "@/components/home/AddAccountButton";
import { UserIDCard } from "@/components/home/UserIDCard";
import { AddAccountModal } from "@/components/home/AddAccountModal";
import { useRouter } from "expo-router";
import { useClerk } from "@clerk/clerk-expo";
import { useAuthStore } from "@/store/authStore";
import { useState } from "react";
import { AccountDetailsModal } from "@/components/home/AccountDetailsModal.tsx";

const verifiedAccounts = [
  { id: "1", name: "Rexan", status: "verified" },
  { id: "2", name: "Respress", status: "verified" },
  { id: "3", name: "Peking", status: "verified" },
];

export default function DashboardScreen() {
  const theme = useTheme();
  const insets = useSafeAreaInsets();
  const { setAuthenticated } = useAuthStore();
  const [modalVisible, setModalVisible] = useState(false);
  const [pendingAccounts, setPendingAccounts] = useState([]);
 
  const [detailsModalVisible, setDetailsModalVisible] = useState(false);
  const [selectedAccount, setSelectedAccount] = useState(null);

  const { signOut } = useClerk();
  const router = useRouter();

  const handleViewDetails = (account) => {
    setSelectedAccount(account);
    setDetailsModalVisible(true);
  };

  const handleDisconnectAccount = (accountId) => {
    // Remove from verified accounts if it's a verified account
    const account = allAccounts.find(acc => acc.id === accountId);
    if (account?.status === 'verified') {
      // You might want to add state for verified accounts to modify them
      console.log('Disconnecting verified account:', accountId);
    } else {
      // Remove from pending accounts
      setPendingAccounts(prev => prev.filter(account => account.id !== accountId));
    }
  };

  const handleAddAccount = () => {
    setModalVisible(true);
  };

  const handleIntegrationSelect = (integration) => {
    const newPendingAccount = {
      id: Date.now().toString(),
      name: integration.name,
      status: "pending",
    };
    setPendingAccounts(prev => [...prev, newPendingAccount]);
    setModalVisible(false);
  };

  const allAccounts = [...verifiedAccounts, ...pendingAccounts];

  const handleRemoveAccount = (accountId) => {
    setPendingAccounts(prev => prev.filter(account => account.id !== accountId));
  }

  return (
    <SafeAreaView
      style={{
        flex: 1,
        backgroundColor: theme.colors.background,
        paddingTop: insets.top,
      }}
    >
      <ScrollView
        className="flex-1"
        style={{ paddingHorizontal: 16 }}
        showsVerticalScrollIndicator={false}
      >
        <DashboardHeader />

        <View style={{ marginBottom: theme.spacing.xl }}>
          <UserIDCard />
        </View>

        <View style={{ marginBottom: theme.spacing.xl }}>
          <Text
            style={{
              fontSize: 24,
              color: theme.colors.textPrimary,
              fontFamily: theme.fonts.welcomeHeading,
              marginBottom: theme.spacing.lg,
            }}
          >
            {pendingAccounts.length > 0 ? "Accounts" : "Verified Accounts"}
          </Text>
          <AccountsList accounts={allAccounts} onRemoveAccount={handleRemoveAccount} onViewDetails={handleViewDetails}/>
        </View>

        <View style={{ alignItems: "center", paddingBottom: theme.spacing.xl }}>
          <AddAccountButton onPress={handleAddAccount} />
        </View>
      </ScrollView>

      <AddAccountModal
        visible={modalVisible}
        onClose={() => setModalVisible(false)}
        onSelectIntegration={handleIntegrationSelect}
      />

      <AccountDetailsModal
        visible={detailsModalVisible}
        account={selectedAccount}
        onClose={() => {
          setDetailsModalVisible(false);
          setSelectedAccount(null);
        }}
        onDisconnect={handleDisconnectAccount}
      />
    </SafeAreaView>
  );
}
