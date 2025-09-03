import React, { useState, useRef } from 'react';
import { View, TextInput, NativeSyntheticEvent, TextInputKeyPressEventData } from 'react-native';
import { InputBox } from './InputBox';

interface OTPInputProps {
  length?: number;
  onComplete: (otp: string) => void;
  value?: string;
}

export const OTPInput: React.FC<OTPInputProps> = ({ 
  length = 4, 
  onComplete,
  value = '' 
}) => {
  const [otp, setOtp] = useState<string[]>(Array(length).fill(''));
  const inputs = useRef<(TextInput | null)[]>([]);

  const handleChangeText = (text: string, index: number) => {
    const newOtp = [...otp];
    newOtp[index] = text;
    setOtp(newOtp);

    if (text && index < length - 1) {
      inputs.current[index + 1]?.focus();
    }

    if (newOtp.every(digit => digit !== '')) {
      onComplete(newOtp.join(''));
    }
  };

  const handleKeyPress = (e: NativeSyntheticEvent<TextInputKeyPressEventData>, index: number) => {
    if (e.nativeEvent.key === 'Backspace' && !otp[index] && index > 0) {
      inputs.current[index - 1]?.focus();
    }
  };

  return (
    <View className="flex-row justify-between px-4">
      {Array(length).fill(0).map((_, index: number) => (
        <InputBox
          key={index}
          // ref={(ref: TextInput | null) => (inputs.current[index] = ref)}
          variant="box"
          className="w-16 h-16 text-center text-xl font-bold"
          maxLength={1}
          keyboardType="numeric"
          value={otp[index]}
          onChangeText={(text: string) => handleChangeText(text, index)}
          onKeyPress={(e: NativeSyntheticEvent<TextInputKeyPressEventData>) => handleKeyPress(e, index)}
          autoFocus={index === 0}
        />
      ))}
    </View>
  );
};