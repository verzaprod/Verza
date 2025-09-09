<div align="center">
  <h1>🔐 Verza</h1>
  <p><strong>Verify Once, Use Everywhere</strong></p>
  <p>A comprehensive reusable KYC (Know Your Customer) verification platform built with React Native and Expo Router</p>
  
  [![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
  [![Expo SDK](https://img.shields.io/badge/Expo%20SDK-53.0.22-blue.svg)](https://expo.dev/)
  [![React Native](https://img.shields.io/badge/React%20Native-0.79.5-green.svg)](https://reactnative.dev/)
  [![TypeScript](https://img.shields.io/badge/TypeScript-5.8.3-blue.svg)](https://www.typescriptlang.org/)
</div>

---

## 🌟 Overview

**Verza** is a revolutionary KYC verification platform that allows users to complete their identity verification once and reuse it across multiple platforms and services. Built with modern technologies including React Native, Expo Router, and NativeWind, Verza provides a seamless cross-platform experience for both mobile and web applications.

### ✨ Key Features

- 🔒 **One-Time Verification** - Complete KYC once, use everywhere
- 📱 **Cross-Platform Support** - iOS, Android, and Web compatibility
- 🎨 **Modern UI/UX** - Built with NativeWind (Tailwind CSS for React Native)
- 🔐 **Secure Integration** - Powered by Onfido SDK for reliable identity verification
- ⚡ **Lightning Fast** - Optimized performance with Expo Router
- 🌓 **Dark Mode Support** - Automatic light/dark theme switching
- 📐 **Responsive Design** - Optimized for all screen sizes
- 🔧 **TypeScript Ready** - Full type safety and enhanced developer experience

---

## 🚀 Quick Start

### Prerequisites

Before getting started, ensure you have the following installed:

- **Node.js** (v18 or higher) - [Download](https://nodejs.org/)
- **npm** or **yarn** - Package manager
- **Expo CLI** - [Installation Guide](https://docs.expo.dev/get-started/installation/)
- **Git** - Version control

For mobile development:
- **iOS Simulator** (macOS only) - Xcode required
- **Android Studio** - For Android emulator

### Installation

1. **Clone the repository**
   ```bash
   git clone https://github.com/mighty-odewumi/Verza.git
   cd Verza
   ```

2. **Install dependencies**
   ```bash
   npm install
   # or
   yarn install
   ```

3. **Start the development server**
   ```bash
   npm start
   # or
   yarn start
   ```

4. **Run on your preferred platform**
   ```bash
   # iOS (macOS only)
   npm run ios
   
   # Android
   npm run android
   
   # Web
   npm run web
   ```

---

## 📋 Available Scripts

| Command | Description | Usage |
|---------|-------------|-------|
| `npm start` | Start Expo development server | Development |
| `npm run android` | Launch on Android device/emulator | Mobile Testing |
| `npm run ios` | Launch on iOS device/simulator | Mobile Testing |
| `npm run web` | Launch in web browser | Web Testing |
| `npm run deploy` | Export and deploy to web | Production |

---

## 🏗️ Project Architecture

```
src/
├── app/                           # Expo Router pages
│   ├── _layout.tsx               # Main layout with providers
│   ├── index.tsx                 # Entry point (redirects to splash)
│   ├── splash.tsx                # Splash screen
│   ├── onboarding/
│   │   ├── _layout.tsx           # Onboarding layout
│   │   ├── index.tsx             # First onboarding slide
│   │   ├── identity.tsx          # Second onboarding slide  
│   │   └── access.tsx            # Third onboarding slide
│   ├── auth/
│   │   ├── _layout.tsx           # Auth layout
│   │   ├── register.tsx          # Registration screen
│   │   ├── verify-email.tsx      # Email verification
│   │   ├── create-pin.tsx        # PIN creation
│   │   ├── backup-passphrase.tsx # Passphrase backup
│   │   ├── confirm-passphrase.tsx# Passphrase confirmation
│   │   └── success.tsx           # Success screen
│   └── home/
│       └── index.tsx             # Home/KYC entry screen
│
├── components/                    # Reusable components
│   ├── ui/                       # Basic UI components
│   │   ├── Icon.tsx              # Generic icon component
│   │   ├── CTAButton.tsx         # Primary action button
│   │   ├── InputBox.tsx          # Base input component
│   │   ├── OTPInput.tsx          # OTP/PIN input boxes
│   │   ├── CircularProgress.tsx  # Circular progress for onboarding
│   │   └── BackButton.tsx        # Reusable back button
│   ├── layout/                   # Layout components
│   │   ├── SafeLayout.tsx        # Safe area wrapper
│   │   └── KeyboardAwareLayout.tsx # Keyboard handling wrapper
│   ├── onboarding/               # Onboarding-specific components
│   │   ├── OnboardingSlide.tsx   # Base slide component
│   │   └── CircularNextButton.tsx # Next button with progress
│   ├── auth/                     # Auth-specific components
│   │   ├── PassphraseGrid.tsx    # 3x4 passphrase grid
│   │   ├── WordChip.tsx          # Selectable word chips
│   │   └── VerificationHeader.tsx # Header for verification screens
│   └── AnimatedSplash.tsx        # Splash screen animation
│
├── theme/                        # Theme system
│   ├── tokens.ts                 # Design tokens
│   ├── ThemeProvider.tsx         # Theme context provider
│   └── index.ts                  # Theme exports
│
├── store/                        # State management
│   ├── authStore.ts              # Authentication state
│   ├── onboardingStore.ts        # Onboarding progress
│   └── index.ts                  # Store exports
│
├── api/                          # API layer
│   ├── client.ts                 # Base API client
│   ├── auth.ts                   # Auth API calls
│   ├── wallet.ts                 # Wallet API calls
│   └── types.ts                  # API response types
│
├── types/                        # TypeScript types
│   ├── auth.ts                   # Auth-related types
│   ├── navigation.ts             # Navigation types
│   └── index.ts                  # Type exports
│
├── utils/                        # Utility functions
│   ├── validation.ts             # Input validation
│   ├── storage.ts                # Secure storage helpers
│   ├── clipboard.ts              # Clipboard operations
│   └── index.ts                  # Utility exports
│
└── global.css                    # NativeWind global styles
```
---

## 🛠️ Technology Stack

### Core Framework
- **[React Native](https://reactnative.dev/)** `0.79.5` - Cross-platform mobile development
- **[Expo](https://expo.dev/)** `53.0.22` - React Native platform and toolchain
- **[Expo Router](https://docs.expo.dev/router/)** `5.1.5` - File-based routing system
- **[TypeScript](https://www.typescriptlang.org/)** `5.8.3` - Type-safe JavaScript

### Styling & UI
- **[NativeWind](https://www.nativewind.dev/)** `4.0.1` - Tailwind CSS for React Native
- **[Tailwind CSS](https://tailwindcss.com/)** `3.4.0` - Utility-first CSS framework
- **[React Native Reanimated](https://docs.swmansion.com/react-native-reanimated/)** `3.17.4` - Smooth animations

### KYC Integration
- **[Onfido React Native SDK](https://github.com/onfido/react-native-sdk)** `15.0.0` - Identity verification

### Navigation & Safety
- **[React Native Safe Area Context](https://github.com/th3rdwave/react-native-safe-area-context)** `5.4.0` - Safe area handling
- **[React Native Screens](https://github.com/software-mansion/react-native-screens)** `4.11.1` - Native screen management

---

## ⚙️ Configuration

### Expo Configuration

The `app.json` file contains essential app configuration:

```json
{
  "name": "Verza",
  "slug": "Verza",
  "scheme": "verza",
  "userInterfaceStyle": "automatic",
  "orientation": "default",
  "web": {
    "output": "static"
  }
}
```

### TypeScript Configuration

Path mapping is configured for cleaner imports:

```json
{
  "compilerOptions": {
    "baseUrl": ".",
    "paths": {
      "@/*": ["./src/*"]
    }
  }
}
```

### Tailwind CSS Setup

Custom configuration for React Native compatibility:

```javascript
module.exports = {
  content: ["./src/**/*.{js,jsx,ts,tsx}"],
  presets: [require("nativewind/preset")],
  theme: {
    extend: {},
  },
  future: {
    hoverOnlyWhenSupported: true,
  },
  plugins: [],
};
```

---

## 🎨 Styling Guide

Verza uses **NativeWind** to bring Tailwind CSS to React Native with full feature parity.

### Basic Usage

```tsx
import { View, Text } from 'react-native';

export default function Component() {
  return (
    <View className="flex-1 bg-white p-4">
      <Text className="text-2xl font-bold text-center text-gray-900 dark:text-white">
        Welcome to Verza
      </Text>
    </View>
  );
}
```


### Platform-Specific Styles

```tsx
<View className="p-4 web:shadow-lg ios:shadow-sm android:elevation-2">
  <Text className="text-base web:text-lg native:text-sm">
    Platform-aware styling
  </Text>
</View>
```

### Dark Mode Support

```tsx
<View className="bg-white dark:bg-gray-900">
  <Text className="text-gray-900 dark:text-white">
    Automatic dark mode
  </Text>
</View>
```

---

## 🔐 KYC Integration

Verza integrates with **Onfido SDK** for robust identity verification:

### Setup Process

1. **Obtain Onfido API Credentials**
   - Sign up for an Onfido account
   - Generate API tokens from the dashboard
   - Configure webhook endpoints

2. **Environment Configuration**
   ```env
   ONFIDO_API_TOKEN=your_onfido_token_here
   ONFIDO_WEBHOOK_SECRET=your_webhook_secret
   EXPO_PUBLIC_API_URL=your_backend_api_url
   ```

3. **Verification Workflow**
   ```tsx
   import { Onfido } from '@onfido/react-native-sdk';
   
   const startVerification = async () => {
     try {
       const result = await Onfido.start({
         sdkToken: 'your_sdk_token',
         flowSteps: {
           welcome: true,
           documentCapture: true,
           faceCapture: true,
         },
       });
       // Handle successful verification
     } catch (error) {
       // Handle verification error
     }
   };
   ```

### Verification Features

- **Document Capture** - Passport, ID cards, driving licenses
- **Biometric Verification** - Facial recognition and liveness detection
- **Address Verification** - Proof of address documents
- **Real-time Processing** - Instant verification results
- **Fraud Detection** - Advanced security measures

---

## 🚀 Deployment

### Web Deployment

Deploy your web app using the built-in export functionality:

```bash
npm run deploy
```

This command:
1. Exports the project for web (`expo export -p web`)
2. Deploys using EAS CLI (`eas-cli deploy`)

### Mobile App Store Deployment

Build and deploy mobile apps using Expo Application Services:

```bash
# Install EAS CLI
npm install -g @expo/eas-cli

# Configure your project
eas build:configure

# Build for iOS and Android
eas build --platform all

# Submit to app stores
eas submit
```

### Environment-Specific Builds

```bash
# Development build
eas build --profile development

# Preview build
eas build --profile preview

# Production build
eas build --profile production
```

---

## 🧪 Development Guidelines

### Code Style

- Use **TypeScript** for all new files
- Follow **React Native best practices**
- Implement **responsive design patterns**
- Use **semantic HTML elements** for web compatibility
- Maintain **accessibility standards**

### Component Structure

```tsx
// MyComponent.tsx
import React from 'react';
import { View, Text } from 'react-native';

interface MyComponentProps {
  title: string;
  subtitle?: string;
}

export default function MyComponent({ title, subtitle }: MyComponentProps) {
  return (
    <View className="p-4 bg-white dark:bg-gray-900">
      <Text className="text-xl font-bold text-gray-900 dark:text-white">
        {title}
      </Text>
      {subtitle && (
        <Text className="text-gray-600 dark:text-gray-300 mt-2">
          {subtitle}
        </Text>
      )}
    </View>
  );
}
```

### File Organization

- **Components**: Reusable UI components in `/src/components/`
- **Screens**: Page-level components using Expo Router in `/src/app/`
- **Utils**: Helper functions in `/src/utils/`
- **Types**: TypeScript definitions in `/src/types/`
- **Hooks**: Custom React hooks in `/src/hooks/`
- **Services**: API and external service integrations in `/src/services/`

---

## 🔧 Environment Variables

Create a `.env.local` file in the project root:

```env
# Onfido Configuration
ONFIDO_API_TOKEN=your_onfido_api_token
ONFIDO_WEBHOOK_SECRET=your_webhook_secret

# API Configuration
EXPO_PUBLIC_API_URL=https://your-api-domain.com
EXPO_PUBLIC_WEB_URL=https://your-web-domain.com

# Development
EXPO_PUBLIC_ENVIRONMENT=development
```

---

## 📱 Platform Support

| Platform | Status | Version |
|----------|--------|---------|
| 📱 iOS | ✅ Supported | iOS 13+ |
| 🤖 Android | ✅ Supported | API 21+ |
| 🌐 Web | ✅ Supported | Modern browsers |
| 💻 macOS | ⏳ Planned | - |
| 🖥️ Windows | ⏳ Planned | - |

---

## 🤝 Contributing

We welcome contributions from the community! Here's how you can help:

### Getting Started

1. **Fork the repository**
2. **Create a feature branch**
   ```bash
   git checkout -b feature/amazing-new-feature
   ```
3. **Make your changes**
4. **Run tests and linting**
   ```bash
   npm run test
   npm run lint
   ```
5. **Commit your changes**
   ```bash
   git commit -m 'Add amazing new feature'
   ```
6. **Push to your branch**
   ```bash
   git push origin feature/amazing-new-feature
   ```
7. **Open a Pull Request**

### Development Setup

```bash
# Clone your fork
git clone https://github.com/YOUR_USERNAME/Verza.git
cd Verza

# Add upstream remote
git remote add upstream https://github.com/mighty-odewumi/Verza.git

# Install dependencies
npm install

# Start development server
npm start
```

### Code Guidelines

- Follow the existing code style
- Write meaningful commit messages
- Add tests for new features
- Update documentation as needed
- Ensure cross-platform compatibility

---

## 📄 License

This project is licensed under the **MIT License** - see the [LICENSE](LICENSE) file for details.

```
MIT License

Copyright (c) 2025 Verza Contributors
```


## 🆘 Support & Help

### Documentation
- 📚 [Expo Documentation](https://docs.expo.dev/)
- 🎨 [NativeWind Documentation](https://www.nativewind.dev/)
- 🔐 [Onfido SDK Documentation](https://documentation.onfido.com/)
- ⚡ [React Native Documentation](https://reactnative.dev/docs/getting-started)

### Community
- 💬 [GitHub Discussions](https://github.com/mighty-odewumi/Verza/discussions)
- 🐛 [Report Issues](https://github.com/mighty-odewumi/Verza/issues)
- 📧 [Email Support](mailto:support@verza.app)

### Quick Links
- 🚀 [Getting Started Guide](https://docs.expo.dev/get-started/installation/)
- 📱 [Expo Router](https://docs.expo.dev/router/introduction/)
- 🎨 [Tailwind CSS](https://tailwindcss.com/docs)
- 🔧 [TypeScript Handbook](https://www.typescriptlang.org/docs/)

---

## 🔗 Related Resources

| Resource | Description | Link |
|----------|-------------|------|
| Expo Router | File-based routing for React Native | [Documentation](https://docs.expo.dev/router/) |
| NativeWind | Tailwind CSS for React Native | [Documentation](https://www.nativewind.dev/) |
| Onfido SDK | Identity verification platform | [GitHub](https://github.com/onfido/react-native-sdk) |
| React Native | Cross-platform mobile development | [Website](https://reactnative.dev/) |
| Tailwind CSS | Utility-first CSS framework | [Documentation](https://tailwindcss.com/) |
| TypeScript | Typed JavaScript at scale | [Website](https://www.typescriptlang.org/) |

---

<div align="center">
  <p><strong>Built with ❤️ by the Verza Team</strong></p>
  <p>Made possible by <a href="https://expo.dev">Expo</a>, <a href="https://reactnative.dev">React Native</a>, and <a href="https://www.nativewind.dev">NativeWind</a></p>
  
  <!-- ⭐ **Star this repository if you find it helpful!** ⭐ -->
</div>

