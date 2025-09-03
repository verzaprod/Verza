"use client"

import type React from "react"
import { useRef } from "react"
import { View, TextInput, type TextInputProps, TouchableOpacity, Text } from "react-native"
import { useTheme } from "@/theme/ThemeProvider"

interface InputBoxesProps extends Omit<TextInputProps, "value" | "onChangeText"> {
  value: string
  onChangeText: (text: string) => void
  length: number
  type?: "otp" | "pin"
}

export const InputBoxes: React.FC<InputBoxesProps> = ({ value, onChangeText, length, type = "otp", ...props }) => {
  const theme = useTheme()
  const inputRef = useRef<TextInput>(null)

  const handlePress = () => {
    inputRef.current?.focus()
  }

  return (
    <View style={{ position: "relative" }}>
      <TextInput
        ref={inputRef}
        value={value}
        onChangeText={onChangeText}
        maxLength={length}
        keyboardType="numeric"
        secureTextEntry={type === "pin"}
        style={{
          position: "absolute",
          opacity: 0,
          width: "100%",
          height: 60,
        }}
        {...props}
      />
      <View
        style={{
          flexDirection: "row",
          justifyContent: "space-between",
          gap: 12,
        }}
      >
        {Array.from({ length }).map((_, index) => (
          <TouchableOpacity
            key={index}
            onPress={handlePress}
            style={{
              flex: 1,
              height: 60,
              borderWidth: 1,
              borderColor: value.length > index ? theme.colors.primaryGreen : "#E2E8F0",
              borderRadius: 8,
              alignItems: "center",
              justifyContent: "center",
              backgroundColor: theme.isDark ? theme.colors.backgroundDark : theme.colors.backgroundLight,
            }}
          >
            <Text
              style={{
                fontSize: 24,
                fontWeight: "600",
                color: theme.isDark ? theme.colors.textPrimaryDark : theme.colors.textPrimaryLight,
              }}
            >
              {type === "pin" && value[index] ? "•" : value[index] || ""}
            </Text>
          </TouchableOpacity>
        ))}
      </View>
    </View>
  )
}
