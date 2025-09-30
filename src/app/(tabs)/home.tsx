import {
  View,
  Text,
  SafeAreaView,
  ScrollView,
  TouchableOpacity,
} from "react-native";
import { useTheme } from "@/theme/ThemeProvider";
import { useSafeAreaInsets } from "react-native-safe-area-context";
import { Icon } from "@/components/ui/Icon";
import { DashboardHeader } from "@/components/home/DashboardHeader";
import { PatternedIDCard } from "@/components/home/PatternedIDCard";
import { AccountsList } from "@/components/home/AccountsList";
import { CircularAddButton } from "@/components/home/CircularAddButton";
import { AddAccountButton } from "@/components/home/AddAccountButton";
import { UserIDCard } from "@/components/home/UserIDCard";
import { AddAccountModal } from "@/components/home/AddAccountModal";
import { useRouter } from "expo-router";
import { useClerk } from "@clerk/clerk-expo";
import { useAuthStore } from "@/store/authStore";
import { useState } from "react";

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
 
  const { signOut } = useClerk();
  const router = useRouter();

  const handleSignOut = async () => {
    try {
      await signOut();
      setAuthenticated(false);
      router.replace("/(auth)/sign-in");
    } catch (err) {
      console.error(JSON.stringify(err, null, 2));
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
          <AccountsList accounts={allAccounts} />
        </View>

        <View style={{ alignItems: "center", paddingBottom: theme.spacing.xl }}>
          <AddAccountButton onPress={handleAddAccount} />
        </View>
        <TouchableOpacity onPress={handleSignOut}>
          <Text>Sign out</Text>
        </TouchableOpacity>
      </ScrollView>

      <AddAccountModal
        visible={modalVisible}
        onClose={() => setModalVisible(false)}
        onSelectIntegration={handleIntegrationSelect}
      />
    </SafeAreaView>
  );
}