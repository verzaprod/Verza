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
import { useRouter } from "expo-router";
import { useClerk } from "@clerk/clerk-expo";

const verifiedAccounts = [
  { id: "1", name: "Respress", status: "verified" },
  { id: "2", name: "Respress", status: "verified" },
  { id: "3", name: "Respress", status: "verified" },
];

export default function DashboardScreen() {
  const theme = useTheme();
  const insets = useSafeAreaInsets();

  const { signOut } = useClerk();
  const router = useRouter();

  const handleSignOut = async () => {
    try {
      await signOut();
      // Redirect to your desired page
      router.replace("/(auth)/sign-in");
    } catch (err) {
      console.error(JSON.stringify(err, null, 2));
    }
  };

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
              fontWeight: "bold",
              color: theme.colors.textPrimary,
              fontFamily: theme.fonts.welcomeHeading,
              marginBottom: theme.spacing.lg,
            }}
          >
            Verified Accounts
          </Text>
          <AccountsList accounts={verifiedAccounts} />
        </View>

        <View style={{ alignItems: "center", paddingBottom: theme.spacing.xl }}>
          <AddAccountButton />
        </View>
        <TouchableOpacity onPress={handleSignOut}>
          <Text>Sign out</Text>
        </TouchableOpacity>
      </ScrollView>
    </SafeAreaView>
  );
}
