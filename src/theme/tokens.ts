export const THEME_TOKENS = {
  colors: {
    primaryGreen: '#16A34A', 
    secondaryAccent: '#22C55E',
    backgroundLight: '#F9FAFB',
    backgroundDark: '#0A0F0D',
    textPrimaryLight: '#111827',
    textPrimaryDark: '#F9FAFB',
    textSecondary: '#6B7280',
    error: '#EF4444',
    success: '#10B981',
  },
  fonts: {
    onboardingHeading: "Urbanist-Bold",
    onboardingTagline: "SFPro",
    welcomeHeading: "Urbanist-Bold",
    welcomeTagline: "Urbanist-ExtraBold",
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