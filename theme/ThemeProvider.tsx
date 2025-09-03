import React, { createContext, useContext, ReactNode } from 'react';
import { useColorScheme } from 'react-native';
import { THEME_TOKENS } from './tokens';

type Theme = typeof THEME_TOKENS & {
  isDark: boolean;
  colors: typeof THEME_TOKENS.colors & {
    background: string;
    textPrimary: string;
  };
};

const ThemeContext = createContext<Theme | null>(null);

interface ThemeProviderProps {
  children: ReactNode;
}

export const ThemeProvider: React.FC<ThemeProviderProps> = ({ children }) => {
  const colorScheme = useColorScheme();
  const isDark = colorScheme === 'dark';

  const theme: Theme = {
    ...THEME_TOKENS,
    isDark,
    colors: {
      ...THEME_TOKENS.colors,
      background: isDark ? THEME_TOKENS.colors.backgroundDark : THEME_TOKENS.colors.backgroundLight,
      textPrimary: isDark ? THEME_TOKENS.colors.textPrimaryDark : THEME_TOKENS.colors.textPrimaryLight,
    },
  };

  return <ThemeContext.Provider value={theme}>{children}</ThemeContext.Provider>;
};

export const useTheme = (): Theme => {
  const theme = useContext(ThemeContext);
  if (!theme) {
    throw new Error('useTheme must be used within a ThemeProvider');
  }
  return theme;
};