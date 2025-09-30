import React from 'react';
import { View, TouchableOpacity } from 'react-native';
import { Icon } from '@/components/ui/Icon';
import { useTheme } from '@/theme/ThemeProvider';
import Feather from '@expo/vector-icons/Feather';

interface StarRatingProps {
  rating: number;
  onRatingChange: (rating: number) => void;
  size?: number;
  readonly?: boolean;
}

export function StarRating({ 
  rating, 
  onRatingChange, 
  size = 32,
  readonly = false 
}: StarRatingProps) {
  const theme = useTheme();

  const handleStarPress = (starRating: number) => {
    if (!readonly) {
      onRatingChange(starRating);
    }
  };

  return (
    <View style={{ flexDirection: 'row', gap: 8 }}>
      {[1, 2, 3, 4, 5].map((star) => (
        <TouchableOpacity
          key={star}
          onPress={() => handleStarPress(star)}
          disabled={readonly}
          style={{
            opacity: readonly ? 1 : 0.8,
          }}
        >
          <Feather
            name="star"
            size={size}
            color={star <= rating ? '#FFD700' : theme.colors.boxBorder}
            fill={star <= rating ? '#FFD700' : 'transparent'}
          />
        </TouchableOpacity>
      ))}
    </View>
  );
}
