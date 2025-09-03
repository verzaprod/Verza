import React from 'react';
import { View, ViewStyle } from 'react-native';
import { SvgProps } from 'react-native-svg';

interface IconProps {
  name: string;
  size?: number;
  color?: string;
  style?: ViewStyle;
  SvgComponent?: React.ComponentType<SvgProps>; // Will be replaced with actual SVGs
}

export const Icon: React.FC<IconProps> = ({ 
  name, 
  size = 24, 
  color = '#000', 
  style,
  SvgComponent 
}) => {
  // Placeholder until actual SVGs are provided
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