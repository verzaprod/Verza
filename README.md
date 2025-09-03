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
Verza/
├── 📁 src/
│   ├── 📁 app/                    # Expo Router Pages
│   │   ├── 📄 _layout.tsx         # Root layout component
│   │   ├── 📄 index.tsx           # Home page
│   │   └── 📄 +not-found.tsx      # 404 error page
│   └── 📄 global.css              # Global Tailwind styles
├── 📁 .expo/                      # Expo build artifacts (auto-generated)
├── ⚙️ app.json                    # Expo app configuration
├── ⚙️ babel.config.js             # Babel transpiler config
├── ⚙️ metro.config.js             # Metro bundler config
├── ⚙️ tailwind.config.js          # Tailwind CSS configuration
├── ⚙️ tsconfig.json               # TypeScript configuration
├── 📄 package.json                # Dependencies and scripts
├── 📄 global.d.ts                 # Global TypeScript definitions
├── 📄 nativewind-env.d.ts         # NativeWind type definitions
└── 📄 README.md                   # Project documentation
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

### Responsive Design

```tsx
<View className="w-full md:w-1/2 lg:w-1/3 xl:w-1/4">
  <Text className="text-sm md:text-base lg:text-lg">
    Responsive text
  </Text>
</View>
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


## Prompt:

make the prompt saying that it should follow robust setup and follow scaleable standards. based on file structure and the composability and modularity of functions and components. 
I want you to plan for the insets or notch for Androids and iPhones.
and then make sure to use an icon for everything you want to use. 
these icons I'll provide from the Figma design and change as the need arise. I will need you to use a library that can handle svgs well. Also, note that for the dark and light mode, there are subtle differences in color. On light mode which is the predominant mode, the text is dark while on light mode, text is white. The CTA Buttons (like Continue, Verify Code, Create Pin, I've Saved my Passphase) should be reusable and composable and have the Primary Green as the background color. The CTA buttons all have a little shadow applied to them to create the impression of being off the page.

The themes are Primary Green, Secondary Accent, Background Light, Background dark, Text Primary Light, Text Primary dark, text secondary, error, success. I will provide the respective color codes they represent. 

I want to follow a DRY principle. And separation of concerns. ANd  best practices that is expected of a Senior Mobile Developer. The backend calls should be modular and arranged in a way that they can be extended or changed as the need arises. 

for most of the texts on the app on the light mode, they use the Text Primary Light as their color unless otherwise stated. 

The splashscreen Shows Verza icon centered and then Animates into a logo that shifts left, with the “Verza” text appearing on the right. At the bottom right corner, an image animates into the frame and stays at that bottom corner. 

For the onboarding, there's a back pressable icon at the top left with a colored primary green background. At the right side, there's a Skip pressable unstyled  that just takes the user to the Signup screen. Note that all element is centered except those back and skip pressables. Note that the Skip and Back are not CTA buttons and should not be styled as such. 
Below this, is the first image provided as a SVG and handled by that library I talked about the other time. Below the image, there's a text saying "Onboard Smarter", below this is another text saying "Skip the long forms. With Verza, onboarding takes minutes, not hours."
At the bottom of each onboarding screen is a prominent next button. But note that it isn't a text. It is a circle with the theme primary green having a white right angle bracket icon. This circle has a circular progress bar surrounding it. On the first page, it's barely half-way and then increases for each of the onboarding screens. 

for 2nd onboarding, there's the same back button and Skip. Below is the text "Identity, Simplified". below is the text: "One secureID, verified once, used everywhere - safely and instantly."
Below is the SVG image and then following the image is the Next button with the progress bar. I don't know what to call it. 

for the 3rd onboarding, there's the back and skip buttons. Below is the image, followed by the text "Seamless Access". below this is the text, "Move freely across platforms with one account. Verza keeps it effortless." Finally below it is the Next button representation as I described before with the progress bar completing a full circle. 

For the Register page, the user sees a predominantly centered text that says "Welcome to Verza". Below this is a text that says "Enter your email or phone number to get started" with the color: Text Primary dark. the two texts are left aligned. below this is the SVG image i'll provide. Finally below is: a input field stretching end to end that has a placeholder of "Enter email or phone number". The border radius is rounded-full. The button below it has the background color of the Primary Green theme set. It has a Continue text. The same rounded and end to end styling applies to the button too. Note that there should be a loading indicator. 

For the Verify your email screen, there's a long Back arrow that when clicked navigates back and standing alone except switching colors based on light or dark modes. Below it is the Heading "Verify Your Email" and below it is a text styled Text Secondary for both light and dark modes: "We have sent the verification code to your email address." the two texts are left aligned.
below this content is the OTP input field that has 4 boxes with thin borders on all sides and whose active state is the primary theme color. they are arranged to fit the screen. And then below this is the Verify Code button. Below that is the text "Didn't receive code? Resend". The Resend text is in the Primary Green theme color. There should be a loading indicator for loading states.

Similarly for the Create PIn screen: there's a Header saying "Create Your PIN" styled as Text Primary Light for light mode and Text Primary Dark for dark mode. Also, the text styled TExt Secondary below, is "Create a 4-digit pin to secure your wallet". It should have 4 boxes as the Verify OTP screen and then styled similarly. At the bottom should be a Create PIN button. 


After this is a Backup passphrase screen, it has a Back button, below this is the "Backup Passphrase" heading and below in Text Secondary is the "Write down these 12 words in order."
And also "You'll need them to recover your wallet." on another line with the same styling. the two texts are left aligned.
Below these texts is a 3x4 grid. The grid is a combination of small cards. Each card has its content centered. and for each card, you have a number at the top and then a word below (e.g. 1. abandon, 2. ability, 3. family, etc.) and then having the same background color of Text Secondary but with reduced transparency to 10%. These words are returned from the backend and so there should be appropriate handling of that. After this grid at the bottom are two buttons with icons for each of them which I'll provide. the buttons are Copy and Save where they Copy saves the passphrases to the user's clipbaord, and then the Save downloads it into the user's device. At the bottom is the "I've Saved My Passphrase" button. 

Now for the Confirm Passphrase screen, we have a heading "Confirm Passphrase", then the tagline is "" in Text Secondary. the two texts are left aligned. There is a grid of 12 pressable words and the user selects 3 right ones in the right order, the active color for selection is Secondary Accent. Then at the bottom is the Continue button. 

Next is the Success page, showing a checkmark image I'll provide, also, centered on the page is a Success message and then finally there's a Start Verification button.


__________________________________________________-
________________________________________--

_____________________________________________________


PROJECT: Verza (mobile) — React Native with Expo Router

GOAL
Implement the mobile UI flow and components for Verza following robust, scalable, and composable standards:
- File-based navigation (Expo Router)
- State management: Zustand
- Styling: NativeWind (Tailwind)
- SVG handling: react-native-svg (and transformer support)
- Secure storage: Expo SecureStore (PIN, minimal secrets)
- Animation: react-native-reanimated
- Keep code modular, testable, and DRY. Separate view logic, presentation, and API networking. Make components composable and theme-driven.


Icons and SVG illustrations will be provided by Figma later — use a simple generic Icon component that accepts an SVG source prop and name for replacement later.

REQUIREMENTS
1. Theming
  - Provide theme tokens: PrimaryGreen, SecondaryAccent, BackgroundLight, BackgroundDark, TextPrimaryLight, TextPrimaryDark, TextSecondary, Error, Success (exact hex codes will be provided later). Implement a theme provider/hooks that choose correct tokens for light & dark modes.
  - Default text colors: in Light mode most texts use TextPrimaryLight (dark colored). In Dark mode most use TextPrimaryDark (light colored).
  - CTA Buttons: reusable/ composable component with PrimaryGreen background, subtle shadow, and loading state. Use same CTA for Continue, Verify, Create PIN, I've saved my passphrase.

2. Navigation & Flow
  - Use Expo Router file-based system. Pages needed: Splash, Onboarding (3 slides), Register (email/phone), Verify Email (OTP), Create PIN (four boxes), Backup Passphrase (3x4 grid), Confirm Passphrase (pick 3 in order), Success, and Home (KYC entry).
  - Respect safe area / notches on iOS and Android (use SafeAreaView + insets via react-native-safe-area-context).
  - Back & Skip: Back icon top-left with PrimaryGreen round background. Skip top-right plain text (unstyled). Both not CTAs.

3. Splash & Onboarding UI specifics
  - Splash: centered Verza icon animates left into a left-aligned logo + Verza text; a small SVG animates into bottom-right and stays. Provide re-usable AnimatedSplash component using reanimated.
  - Onboarding slides: each has back & skip, a center SVG, Title and subtitle, and a circular Next button (PrimaryGreen circle with white right-angle bracket icon). The next button has an outer circular progress bar which goes 1/3, 2/3, full on pages 1–3. The next button is visually prominent — implement with react-native-svg Circle progress.

4. Inputs & buttons
  - Register: left-aligned heading ("Welcome to Verza") + subtext; full-width rounded input (rounded-full) placeholder "Enter email or phone number"; Primary CTA below (Continue) with loading state.
  - Verify Email: back arrow on its own; heading & text left aligned; OTP input of 4 boxes (thin border, active border = primary green); Verify Code button; "Didn't receive code? Resend" in PrimaryGreen.
  - Create PIN: header & secondary text left aligned; 4-box PIN input visually same as OTP; Create PIN CTA.
  - Backup Passphrase: back button; heading + two lines descriptive text left aligned; 3x4 grid of cards — each card: number and word (words returned from backend) with background = TextSecondary at 10% opacity; Copy & Save buttons with icons; "I've Saved My Passphrase" CTA.
  - Confirm Passphrase: heading + tagline left aligned; words grid (shuffled) as pressable chips; selected chips show SecondaryAccent; bottom Continue button (disabled until correct 3 words selected in right order).
  - Success: centered check image, success text, Start Verification CTA.

5. Accessibility & UX
  - Inputs should be keyboard-friendly, focusable, and support screen readers.
  - All interactive elements must use SafeArea and have adequate hit area.
  - Provide fallback for missing SVGs / icons (placeholder).

6. Implementation quality
  - Components must be small and composable: Icon, CTAButton, PinBoxes/OTPBoxes, CardGrid, AnimatedSplash, CircularNextButton.
  - Networking calls must be abstracted to a modular API client (e.g., src/api/onfido or src/api/auth). The prompt assumes actual API details live in project context; do not hard-code URLs. Provide sample interface signatures for the APIs.

OUTPUT
- Provide code for the components and pages with clear file structure.
- Use comments where integration with real backend or assets is required.
- Keep the implementation flexible: do not hard-wire colors, images, or API URLs. Read color tokens from a theme file.

NOTE
- Use surrounding project context if available (existing stores, theme files, API clients); if not available, create small modular placeholders that can be replaced.

it should follow robust setup and follow scaleable standards. based on file structure and the composability and modularity of functions and components. 
I want you to plan for the insets or notch for Androids and iPhones.
and then make sure to use an icon for everything you want to use. 
these icons I'll provide from the Figma design and change as the need arise. I will need you to use a library that can handle svgs well. Also, note that for the dark and light mode, there are subtle differences in color. On light mode which is the predominant mode, the text is dark while on light mode, text is white. The CTA Buttons (like Continue, Verify Code, Create Pin, I've Saved my Passphase) should be reusable and composable and have the Primary Green as the background color. The CTA buttons all have a little shadow applied to them to create the impression of being off the page.

The themes are Primary Green, Secondary Accent, Background Light, Background dark, Text Primary Light, Text Primary dark, text secondary, error, success. I will provide the respective color codes they represent. 

I want to follow a DRY principle. And separation of concerns. ANd  best practices that is expected of a Senior Mobile Developer. The backend calls should be modular and arranged in a way that they can be extended or changed as the need arises. 

for most of the texts on the app on the light mode, they use the Text Primary Light as their color unless otherwise stated. 

The splashscreen Shows Verza icon centered and then Animates into a logo that shifts left, with the “Verza” text appearing on the right. At the bottom right corner, an image animates into the frame and stays at that bottom corner. 

For the onboarding, there's a back pressable icon at the top left with a colored primary green background. At the right side, there's a Skip pressable unstyled  that just takes the user to the Signup screen. Note that all element is centered except those back and skip pressables. Note that the Skip and Back are not CTA buttons and should not be styled as such. 
Below this, is the first image provided as a SVG and handled by that library I talked about the other time. Below the image, there's a text saying "Onboard Smarter", below this is another text saying "Skip the long forms. With Verza, onboarding takes minutes, not hours."
At the bottom of each onboarding screen is a prominent next button. But note that it isn't a text. It is a circle with the theme primary green having a white right angle bracket icon. This circle has a circular progress bar surrounding it. On the first page, it's barely half-way and then increases for each of the onboarding screens. 

for 2nd onboarding, there's the same back button and Skip. Below is the text "Identity, Simplified". below is the text: "One secureID, verified once, used everywhere - safely and instantly."
Below is the SVG image and then following the image is the Next button with the progress bar. I don't know what to call it. 

for the 3rd onboarding, there's the back and skip buttons. Below is the image, followed by the text "Seamless Access". below this is the text, "Move freely across platforms with one account. Verza keeps it effortless." Finally below it is the Next button representation as I described before with the progress bar completing a full circle. 

For the Register page, the user sees a predominantly centered text that says "Welcome to Verza". Below this is a text that says "Enter your email or phone number to get started" with the color: Text Primary dark. the two texts are left aligned. below this is the SVG image i'll provide. Finally below is: a input field stretching end to end that has a placeholder of "Enter email or phone number". The border radius is rounded-full. The button below it has the background color of the Primary Green theme set. It has a Continue text. The same rounded and end to end styling applies to the button too. Note that there should be a loading indicator. 

For the Verify your email screen, there's a long Back arrow that when clicked navigates back and standing alone except switching colors based on light or dark modes. Below it is the Heading "Verify Your Email" and below it is a text styled Text Secondary for both light and dark modes: "We have sent the verification code to your email address." the two texts are left aligned.
below this content is the OTP input field that has 4 boxes with thin borders on all sides and whose active state is the primary theme color. they are arranged to fit the screen. And then below this is the Verify Code button. Below that is the text "Didn't receive code? Resend". The Resend text is in the Primary Green theme color. There should be a loading indicator for loading states.

Similarly for the Create PIn screen: there's a Header saying "Create Your PIN" styled as Text Primary Light for light mode and Text Primary Dark for dark mode. Also, the text styled TExt Secondary below, is "Create a 4-digit pin to secure your wallet". It should have 4 boxes as the Verify OTP screen and then styled similarly. At the bottom should be a Create PIN button. 


After this is a Backup passphrase screen, it has a Back button, below this is the "Backup Passphrase" heading and below in Text Secondary is the "Write down these 12 words in order."
And also "You'll need them to recover your wallet." on another line with the same styling. the two texts are left aligned.
Below these texts is a 3x4 grid. The grid is a combination of small cards. Each card has its content centered. and for each card, you have a number at the top and then a word below (e.g. 1. abandon, 2. ability, 3. family, etc.) and then having the same background color of Text Secondary but with reduced transparency to 10%. These words are returned from the backend and so there should be appropriate handling of that. After this grid at the bottom are two buttons with icons for each of them which I'll provide. the buttons are Copy and Save where they Copy saves the passphrases to the user's clipbaord, and then the Save downloads it into the user's device. At the bottom is the "I've Saved My Passphrase" button. 

Now for the Confirm Passphrase screen, we have a heading "Confirm Passphrase", then the tagline is "" in Text Secondary. the two texts are left aligned. There is a grid of 12 pressable words and the user selects 3 right ones in the right order, the active color for selection is Secondary Accent. Then at the bottom is the Continue button. 

Next is the Success page, showing a checkmark image I'll provide, also, centered on the page is a Success message and then finally there's a Start Verification button.

also, make sure that the Keyboard doesn't block the view for the input fields when it is active. and animations should be the smoothest I've ever seen. as much as possible, stick to the design given. and make the maximum number of lines for each file to be 150 and anything other than that is a code smell that needs real refactoring. 



// "@types/react-native": "^0.72.8",