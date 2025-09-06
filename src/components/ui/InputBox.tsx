import React, { useState } from 'react';
import { TextInput, View, TextInputProps, ViewStyle, StyleProp } from 'react-native';
import Animated, { useSharedValue, useAnimatedStyle, withTiming } from 'react-native-reanimated';
import { useTheme } from '@/theme/ThemeProvider';

interface InputBoxProps extends Omit<TextInputProps, 'style'> {
  variant?: 'rounded' | 'box';
  active?: boolean;
  style?: StyleProp<ViewStyle>;
}

export const InputBox: React.FC<InputBoxProps> = ({
  variant = 'rounded',
  active = false,
  style,
  onFocus,
  onBlur,
  ...props
}) => {
  const theme = useTheme();
  const [isFocused, setIsFocused] = useState(false);
  const borderColor = useSharedValue<string>(theme.colors.textSecondary);

  const animatedStyle = useAnimatedStyle(() => ({
    borderColor: borderColor.value,
  }));

  const handleFocus = (e: any) => {
    setIsFocused(true);
    borderColor.value = withTiming(theme.colors.primaryGreen);
    onFocus?.(e);
  };

  const handleBlur = (e: any) => {
    setIsFocused(false);
    borderColor.value = withTiming(theme.colors.textSecondary);
    onBlur?.(e);
  };

  return (
    <Animated.View
      style={[
        {
          borderWidth: variant === 'box' ? 1 : .5,
          borderRadius: variant === 'rounded' ? theme.borderRadius.full : theme.borderRadius.md,
          backgroundColor: theme.colors.background,
        },
        animatedStyle,
        style,
      ]}
    >
      <TextInput
        className={`px-4 ${variant === 'rounded' ? 'py-4' : 'py-3'} text-lg`}
        style={{ color: theme.colors.textPrimary }}
        placeholderTextColor={theme.colors.textSecondary}
        onFocus={handleFocus}
        onBlur={handleBlur}
        {...props}
      />
    </Animated.View>
  );
};