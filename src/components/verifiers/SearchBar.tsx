import React from "react";
import { View, TextInput, TouchableOpacity } from "react-native";
import { useTheme } from "@/theme/ThemeProvider";
import { Icon } from "@/components/ui/Icon";
import Feather from "@expo/vector-icons/Feather";

interface SearchBarProps {
  value: string;
  onSearch: (query: string) => void;
  onClear: () => void;
  placeholder?: string;
}

export function SearchBar({ 
  value, 
  onSearch, 
  onClear, 
  placeholder = "Search..." 
}: SearchBarProps) {
  const theme = useTheme();

  return (
    <View
      style={{
        flexDirection: "row",
        alignItems: "center",
        backgroundColor: theme.colors.background,
        borderRadius: theme.borderRadius.lg,
        paddingHorizontal: theme.spacing.md,
        paddingVertical: theme.spacing.sm,
        borderWidth: 1,
        borderColor: theme.colors.boxBorder,
      }}
    >
      <Feather
        name="search"
        size={20}
        color={theme.colors.textSecondary}
        style={{ marginRight: theme.spacing.sm }}
      />
      
      <TextInput
        style={{
          flex: 1,
          fontSize: 16,
          color: theme.colors.textPrimary,
          paddingVertical: theme.spacing.xs,
        }}
        placeholder={placeholder}
        placeholderTextColor={theme.colors.textSecondary}
        value={value}
        onChangeText={onSearch}
        autoCapitalize="none"
        autoCorrect={false}
      />

      {value.length > 0 && (
        <TouchableOpacity
          onPress={onClear}
          style={{
            marginLeft: theme.spacing.sm,
            padding: theme.spacing.xs,
          }}
          hitSlop={{ top: 8, bottom: 8, left: 8, right: 8 }}
        >
          <Feather
            name="x"
            size={18}
            color={theme.colors.textSecondary}
          />
        </TouchableOpacity>
      )}
    </View>
  );
}
