import React, { useState } from 'react';
import { View, Text, TextInput, TouchableOpacity, Alert } from 'react-native';
import { useTheme } from '@/theme/ThemeProvider';
import { StarRating } from '@/components/ui/StarRating';
import { useRatingStore } from '@/store/ratingStore';

interface VerifierRatingProps {
  verifierId: string;
  verifierName: string;
  onRatingSubmitted?: () => void;
}

export function VerifierRating({ 
  verifierId, 
  verifierName, 
  onRatingSubmitted 
}: VerifierRatingProps) {
  const theme = useTheme();
  const {
    currentRating,
    currentComment,
    isSubmitting,
    setCurrentRating,
    setCurrentComment,
    submitRating,
    getRatingForVerifier
  } = useRatingStore();

  const [hasSubmitted, setHasSubmitted] = useState(false);
  const existingRating = getRatingForVerifier(verifierId);

  const handleSubmitRating = async () => {
    if (currentRating === 0) {
      Alert.alert('Rating Required', 'Please select a star rating before submitting.');
      return;
    }

    try {
      await submitRating(verifierId);
      setHasSubmitted(true);
      Alert.alert('Thank You!', 'Your rating has been submitted successfully.');
      onRatingSubmitted?.();
    } catch (error) {
      Alert.alert('Error', 'Failed to submit rating. Please try again.');
    }
  };

  if (existingRating || hasSubmitted) {
    return (
      <View style={{
        backgroundColor: theme.colors.background,
        borderRadius: theme.borderRadius.lg,
        padding: theme.spacing.lg,
        marginBottom: theme.spacing.lg,
      }}>
        <Text style={{
          fontSize: 16,
          fontWeight: '600',
          color: theme.colors.textPrimary,
          marginBottom: theme.spacing.sm,
        }}>
          Your Rating for {verifierName}
        </Text>
        
        <StarRating 
          rating={existingRating?.rating || currentRating} 
          onRatingChange={() => {}} 
          readonly={true}
          size={24}
        />
        
        {(existingRating?.comment || currentComment) && (
          <Text style={{
            fontSize: 14,
            color: theme.colors.textSecondary,
            marginTop: theme.spacing.sm,
            fontStyle: 'italic',
          }}>
            "{existingRating?.comment || currentComment}"
          </Text>
        )}
      </View>
    );
  }

  return (
    <View style={{
      backgroundColor: theme.colors.background,
      borderRadius: theme.borderRadius.lg,
      padding: theme.spacing.lg,
      marginBottom: theme.spacing.lg,
    }}>
      <Text style={{
        fontSize: 16,
        fontWeight: '600',
        color: theme.colors.textPrimary,
        marginBottom: theme.spacing.sm,
      }}>
        Rate Your Experience
      </Text>
      
      <Text style={{
        fontSize: 14,
        color: theme.colors.textSecondary,
        marginBottom: theme.spacing.lg,
      }}>
        How was your verification with {verifierName}?
      </Text>

      <View style={{ alignItems: 'center', marginBottom: theme.spacing.lg }}>
        <StarRating 
          rating={currentRating} 
          onRatingChange={setCurrentRating}
        />
      </View>

      <TextInput
        style={{
          backgroundColor: theme.colors.background,
          borderWidth: 1,
          borderColor: theme.colors.boxBorder,
          borderRadius: theme.borderRadius.md,
          padding: theme.spacing.md,
          fontSize: 14,
          color: theme.colors.textPrimary,
          textAlignVertical: 'top',
          marginBottom: theme.spacing.lg,
        }}
        placeholder="Share your feedback (optional)"
        placeholderTextColor={theme.colors.textSecondary}
        value={currentComment}
        onChangeText={setCurrentComment}
        multiline={true}
        numberOfLines={3}
        maxLength={200}
      />

      <TouchableOpacity
        style={{
          backgroundColor: currentRating > 0 ? theme.colors.primaryGreen : theme.colors.boxBorder,
          borderRadius: theme.borderRadius.md,
          padding: theme.spacing.md,
          alignItems: 'center',
        }}
        onPress={handleSubmitRating}
        disabled={currentRating === 0 || isSubmitting}
      >
        <Text style={{
          color: currentRating > 0 ? 'white' : theme.colors.textSecondary,
          fontSize: 16,
          fontWeight: '600',
        }}>
          {isSubmitting ? 'Submitting...' : 'Submit Rating'}
        </Text>
      </TouchableOpacity>
    </View>
  );
}
