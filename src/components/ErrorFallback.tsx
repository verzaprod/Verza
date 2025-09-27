import {
  View, Text, TouchableOpacity,
} from "react-native";
import { SafeAreaView, } from "react-native-safe-area-context";

export const ErrorFallback: React.FC<{ onReset: () => void; error?: Error }> = ({ onReset, error }) => {
  return (
    <SafeAreaView style={{ flex: 1, backgroundColor: '#f5f5f5' }}>
      <View style={{ flex: 1, justifyContent: 'center', alignItems: 'center', padding: 20 }}>
        <View
          style={{
            width: 80,
            height: 80,
            backgroundColor: '#ef4444',
            borderRadius: 40,
            alignItems: 'center',
            justifyContent: 'center',
            marginBottom: 24,
          }}
        >
          <Text style={{ color: 'white', fontSize: 40 }}>⚠️</Text>
        </View>

        <Text
          style={{
            fontSize: 24,
            fontWeight: 'bold',
            color: '#1f2937',
            textAlign: 'center',
            marginBottom: 12,
          }}
        >
          Oops! Something went wrong
        </Text>

        <Text
          style={{
            fontSize: 16,
            color: '#6b7280',
            textAlign: 'center',
            lineHeight: 24,
            marginBottom: 32,
          }}
        >
          We encountered an unexpected error. Don't worry, your data is safe.
        </Text>

        {__DEV__ && error && (
          <View
            style={{
              backgroundColor: '#fef2f2',
              borderRadius: 8,
              padding: 16,
              marginBottom: 24,
              maxHeight: 150,
            }}
          >
            <Text
              style={{
                fontSize: 12,
                color: '#dc2626',
                fontFamily: 'monospace',
              }}
            >
              {error.message}
            </Text>
          </View>
        )}

        <TouchableOpacity
          style={{
            backgroundColor: '#10b981',
            paddingVertical: 16,
            paddingHorizontal: 32,
            borderRadius: 8,
            marginBottom: 12,
          }}
          onPress={onReset}
        >
          <Text
            style={{
              color: 'white',
              fontSize: 16,
              fontWeight: '600',
            }}
          >
            Try Again
          </Text>
        </TouchableOpacity>

        <TouchableOpacity
          style={{
            paddingVertical: 12,
            paddingHorizontal: 20,
          }}
          onPress={() => {
            // Optional: Navigate to support or contact screen
            console.log('Report issue clicked')
          }}
        >
          <Text
            style={{
              color: '#6b7280',
              fontSize: 14,
            }}
          >
            Report this issue
          </Text>
        </TouchableOpacity>
      </View>
    </SafeAreaView>
  )
}
