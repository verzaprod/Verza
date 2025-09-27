import React from 'react'
import { View, Text, TouchableOpacity, Image } from 'react-native'
import { useTheme } from '@/theme/ThemeProvider'
import { Icon } from '@/components/ui/Icon'
import Feather from "@expo/vector-icons/Feather";

interface CameraCaptureProps {
  title: string
  subtitle: string
  image: string | null
  onCapture: () => void
  isCircular?: boolean
}

export const CameraCapture: React.FC<CameraCaptureProps> = ({
  title,
  subtitle,
  image,
  onCapture,
  isCircular = false,
}) => {
  const theme = useTheme()

  return (
    <View
      style={{
        backgroundColor: theme.colors.backgroundLight,
        borderRadius: theme.borderRadius.lg,
        padding: 20,
      }}
    >
      <Text
        style={{
          fontSize: 18,
          fontWeight: '600',
          color: theme.colors.textPrimary,
          marginBottom: 4,
        }}
      >
        {title}
      </Text>
      <Text
        style={{
          fontSize: 14,
          color: theme.colors.textSecondary,
          marginBottom: 16,
        }}
      >
        {subtitle}
      </Text>

      <TouchableOpacity
        style={{
          height: 200,
          borderRadius: isCircular ? 100 : theme.borderRadius.md,
          backgroundColor: image ? 'transparent' : theme.colors.background,
          borderWidth: 2,
          borderColor: image ? theme.colors.primaryGreen : theme.colors.textSecondary,
          borderStyle: image ? 'solid' : 'dashed',
          alignItems: 'center',
          justifyContent: 'center',
          overflow: 'hidden',
        }}
        onPress={onCapture}
      >
        {image ? (
          <Image
            source={{ uri: image }}
            style={{
              width: '100%',
              height: '100%',
              borderRadius: isCircular ? 100 : theme.borderRadius.md,
            }}
            resizeMode="cover"
          />
        ) : (
          <View style={{ alignItems: 'center' }}>
            <View
              style={{
                width: 60,
                height: 60,
                borderRadius: 30,
                backgroundColor: theme.colors.primaryGreen + '20',
                alignItems: 'center',
                justifyContent: 'center',
                marginBottom: 12,
              }}
            >
              <Feather name="camera" size={30} />
            </View>
            <Text
              style={{
                fontSize: 16,
                fontWeight: '500',
                color: theme.colors.textPrimary,
              }}
            >
              Tap to capture
            </Text>
          </View>
        )}
      </TouchableOpacity>

      {image && (
        <TouchableOpacity
          style={{
            marginTop: 12,
            alignSelf: 'center',
            paddingVertical: 8,
            paddingHorizontal: 16,
          }}
          onPress={onCapture}
        >
          <Text
            style={{
              color: theme.colors.primaryGreen,
              fontSize: 14,
              fontWeight: '500',
            }}
          >
            Retake Photo
          </Text>
        </TouchableOpacity>
      )}
    </View>
  )
}
