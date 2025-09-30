import React, { useState, useMemo } from "react";
import { View, Text, SafeAreaView, ScrollView } from "react-native";
import { useTheme } from "@/theme/ThemeProvider";
import { useSafeAreaInsets } from "react-native-safe-area-context";
import { VerifiersHeader } from "@/components/verifiers/VerifiersHeader";
import { VerifiersList } from "@/components/verifiers/VerifiersList";
import { SearchBar } from "@/components/verifiers/SearchBar";
import { MOCK_DATA } from "@/services/api/mockData";

const verifiers = MOCK_DATA.verifiers;

export default function VerifiersScreen() {
  const theme = useTheme();
  const insets = useSafeAreaInsets();
  const [searchQuery, setSearchQuery] = useState("");

  // Filter verifiers based on search query
  const filteredVerifiers = useMemo(() => {
    if (!searchQuery.trim()) {
      return verifiers;
    }

    const query = searchQuery.toLowerCase();
    return verifiers.filter(verifier => 
      verifier.name.toLowerCase().includes(query) ||
      verifier.type.toLowerCase().includes(query) ||
      verifier.description.toLowerCase().includes(query)
    );
  }, [searchQuery]);

  const handleSearch = (query: string) => {
    setSearchQuery(query);
  };

  const clearSearch = () => {
    setSearchQuery("");
  };

  return (
    <SafeAreaView
      style={{
        flex: 1,
        backgroundColor: theme.colors.background,
        paddingTop: insets.top + 16,
      }}
    >
      <ScrollView
        className="flex-1"
        style={{ paddingHorizontal: 20 }}
        showsVerticalScrollIndicator={false}
      >
        <View style={{ marginBottom: theme.spacing.lg }}>
          <Text
            style={{
              fontSize: 24,
              color: theme.colors.textPrimary,
              fontFamily: theme.fonts.welcomeHeading,
              marginBottom: theme.spacing.sm,
            }}
          >
            Identity Verifiers
          </Text>
          <Text
            style={{
              fontSize: 16,
              color: theme.colors.textSecondary,
              marginBottom: theme.spacing.lg,
            }}
          >
            Choose a trusted verifier to authenticate your identity
          </Text>
          <SearchBar 
            value={searchQuery}
            onSearch={handleSearch}
            onClear={clearSearch}
            placeholder="Search verifiers..."
          />
        </View>

        <View style={{ paddingBottom: theme.spacing.xl }}>
          {filteredVerifiers.length > 0 ? (
            <VerifiersList verifiers={filteredVerifiers} />
          ) : (
            <View style={{
              alignItems: 'center',
              justifyContent: 'center',
              paddingVertical: theme.spacing.xl * 2,
            }}>
              <Text style={{
                fontSize: 16,
                color: theme.colors.textSecondary,
                textAlign: 'center',
              }}>
                No verifiers found for "{searchQuery}"
              </Text>
              <Text style={{
                fontSize: 14,
                color: theme.colors.textSecondary,
                textAlign: 'center',
                marginTop: theme.spacing.sm,
              }}>
                Try searching with different keywords
              </Text>
            </View>
          )}
        </View>
      </ScrollView>
    </SafeAreaView>
  );
}
