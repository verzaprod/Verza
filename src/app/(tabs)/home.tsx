import { View, Text, SafeAreaView, ScrollView, TouchableOpacity } from "react-native"
import { useTheme } from "@/theme/ThemeProvider"
import { Icon } from "@/components/ui/Icon"
import { UserIDCard } from "@/components/home/UserIdCard"
import { VerifiedAccountsList } from "@/components/home/VerifiedAccountsList"
import { AddAccountButton } from "@/components/home/AddAccountButton"
import { useSafeAreaInsets } from "react-native-safe-area-context"

const verifiedAccounts = [
  { id: '1', name: 'Respress', status: 'verified' },
  { id: '2', name: 'Respress', status: 'verified' },
  { id: '3', name: 'Respress', status: 'verified' },
]

export default function HomeScreen() {
  const theme = useTheme()
  const insets = useSafeAreaInsets();

  return (
    <SafeAreaView 
      style={{ 
        flex: 1, 
        backgroundColor: theme.colors.background,
        paddingTop: insets.top,

      }}

    >
      <ScrollView className="flex-1 px-6" showsVerticalScrollIndicator={false}>
        <View className="flex-row justify-between items-center py-4">
          <TouchableOpacity>
            <Icon name="avatar" size={56} />
          </TouchableOpacity>
          <TouchableOpacity>
            <Icon name="notifications" size={24} />
          </TouchableOpacity>
        </View>

        <View className="mb-8">
          <UserIDCard />
        </View>

        <View className="mb-8">
          <Text
            className="text-2xl font-bold mb-6"
            style={{
              color: theme.colors.textPrimary,
              fontFamily: theme.fonts.welcomeHeading,
            }}
          >
            Verified Accounts
          </Text>
          <VerifiedAccountsList accounts={verifiedAccounts} />
        </View>

        <View className="items-center pb-8">
          <AddAccountButton />
        </View>
      </ScrollView>
    </SafeAreaView>
  )
}
