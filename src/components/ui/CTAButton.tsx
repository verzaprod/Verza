import React from 'react';
import { Text, TouchableOpacity, ActivityIndicator, TouchableOpacityProps } from 'react-native';
import { useTheme } from '@/theme/ThemeProvider';

interface CTAButtonProps extends TouchableOpacityProps {
  title: string;
  loading?: boolean;
  variant?: 'primary' | 'secondary';
}

export const CTAButton: React.FC<CTAButtonProps> = ({
  title,
  loading = false,
  variant = 'primary',
  disabled,
  style,
  ...props
}) => {
  const theme = useTheme();

  return (
    <TouchableOpacity
      className="w-full py-4 rounded-full flex-row items-center justify-center"
      style={[
        {
          backgroundColor: variant === 'primary' ? theme.colors.primaryGreen : 'transparent',
          borderWidth: variant === 'secondary' ? 1 : 0,
          borderColor: theme.colors.primaryGreen,
          opacity: (disabled || loading) ? 0.6 : 1,
          ...theme.shadows.subtle,
        },
        style,
      ]}
      disabled={disabled || loading}
      {...props}
    >
      {loading ? (
        <ActivityIndicator size="small" color="white" />
      ) : (
        <Text 
          className="text-xl font-bold"
          style={{ 
            color: variant === 'primary' ? 'white' : theme.colors.primaryGreen 
          }}
        >
          {title}
        </Text>
      )}
    </TouchableOpacity>
  );
};