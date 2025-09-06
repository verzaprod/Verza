import React from 'react';
import { View, ViewStyle, Image } from 'react-native';

interface IconProps {
  name: string;
  size?: number;
  color?: string;
  style?: ViewStyle;
}

const iconSources: Record<string, any> = {
  'verza-logo': require('@/assets/images/icon-1.png'),
  'onboarding-1': require('@/assets/images/onboarding-1.png'),
  'onboarding-2': require('@/assets/images/onboarding-2.png'),
  'onboarding-3': require('@/assets/images/onboarding-3.png'),
  'chevron-left': require('@/assets/images/chevron-left.png'),
  'chevron-right': require('@/assets/images/chevron-right.png'),
  'welcome': require('@/assets/images/welcome-1.png'),
  'splash-illustration': require('@/assets/images/splash-illustration-1.png'),
};

export const Icon: React.FC<IconProps> = ({ 
  name, 
  size = 24, 
  color, 
  style,
}) => {
  const source = iconSources[name];
  const finalWidth = style?.width || size;
  const finalHeight = style?.height || size;

  if (source) {
    return (
      <View style={[{ width: finalWidth, height: finalHeight }, style]}>
        <Image 
          source={source}
          style={{ 
            width: "100%", 
            height: "100%",
          }}
          resizeMode="contain"
        />
      </View>
    );
  }

  // Fallback
  return (
    <View 
      style={[
        {
          width: finalWidth,
          height: finalHeight,
          backgroundColor: '#E5E7EB',
          borderRadius: 4,
        },
        style,
      ]}
    />
  );
};