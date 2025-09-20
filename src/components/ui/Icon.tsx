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
  'copy': require('@/assets/images/copy.png'),
  'save': require('@/assets/images/save.png'),
  'id-card': require('@/assets/images/id-card.png'),
  'driver-license': require('@/assets/images/driver-license.png'),
  'passport': require('@/assets/images/passport.png'),
  'success': require('@/assets/images/success.png'),
  'face': require('@/assets/images/face.png'),
  'right-pattern': require('@/assets/images/right-pattern.png'),
  'left-pattern': require('@/assets/images/left-pattern.png'),
  'note-list': require('@/assets/images/note-list.png'),
  'shield-check': require('@/assets/images/shield-check.png'),
  'shield': require('@/assets/images/shield.png'),
  'avatar': require('@/assets/images/avatar.png'),
  'cancel': require('@/assets/images/cancel.png'),
  'remove': require('@/assets/images/remove.png'),
  'wifi': require('@/assets/images/wifi.png'),
  'plus': require('@/assets/images/plus.png'),
  'notification': require('@/assets/images/notification.png'),
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