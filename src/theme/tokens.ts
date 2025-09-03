export const THEME_TOKENS = {
  colors: {
    primaryGreen: '', // Will be replaced with actual hex
    secondaryAccent: '#FF6B35', // Will be replaced
    backgroundLight: '#FFFFFF',
    backgroundDark: '#1A1A1A',
    textPrimaryLight: '#2C2C2C',
    textPrimaryDark: '#FFFFFF',
    textSecondary: '#666666',
    error: '#FF4444',
    success: '#00C851',
  },
  spacing: {
    xs: 4,
    sm: 8,
    md: 16,
    lg: 24,
    xl: 32,
    xxl: 48,
  },
  borderRadius: {
    sm: 8,
    md: 12,
    lg: 16,
    full: 9999,
  },
  shadows: {
    subtle: {
      shadowColor: '#000',
      shadowOffset: { width: 0, height: 2 },
      shadowOpacity: 0.1,
      shadowRadius: 4,
      elevation: 3,
    },
  },
} as const;