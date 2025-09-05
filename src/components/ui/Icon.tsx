import React from 'react';
import { View, ViewStyle } from 'react-native';
import { SvgProps } from 'react-native-svg';

import VerzaLogoIcon from "@/assets/images/icon.svg";
import OnboardingImage1 from "@/assets/images/onboarding1.svg";
import OnboardingImage2 from "@/assets/images/onboarding2.svg";
import OnboardingImage3 from "@/assets/images/onboarding3.svg";
import OnboardingBackIcon from "@/assets/images/chevron-left.svg";
import OnboardingNextIcon from "@/assets/images/chevron-right.svg";
// import SplashIllustration from "@/assets/images/splash-illustration.png";

interface IconProps {
  name: string;
  size?: number;
  color?: string;
  style?: ViewStyle;
  SvgComponent?: React.ComponentType<SvgProps>; 
}

const iconMap: Record<string, React.FC<SvgProps>> = {
  'verza-logo': VerzaLogoIcon,
  'onboarding-1': OnboardingImage1,
  'onboarding-2': OnboardingImage2,
  'onboarding-3': OnboardingImage3,
  'chevron-left': OnboardingBackIcon,
  'chevron-right': OnboardingNextIcon,
  // 'splash-illustration': SplashIllustration,
};

export const Icon: React.FC<IconProps> = ({ 
  name, 
  size = 24, 
  color = '#000', 
  style, 
}) => {
  const SvgComponent = iconMap[name];

  if (SvgComponent) {
    return (
      <View style={[{ width: size, height: size }, style]}>
        <SvgComponent width={size} height={size} color={color} />
      </View>
    );
  }

  return (
    <View 
      style={[
        {
          width: size,
          height: size,
          backgroundColor: color,
          borderRadius: size / 2,
        },
        style,
      ]}
    />
  );
};