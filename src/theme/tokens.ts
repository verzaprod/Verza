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
    boxBorder: "#DDDDDD",
    boxText: "#757575",
  },
  fonts: {
    onboardingHeading: "UrbanistBold",
    onboardingTagline: "SFPro",
    welcomeHeading: "UrbanistBold",
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
      shadowColor: '#000000',
      shadowOffset: { width: 10, height: 28 },
      shadowOpacity: 0,
      shadowRadius: 10,
      elevation: 7,
    },
  },
} as const;