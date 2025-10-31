# Verza

**Verify Once, Use Everywhere**

A cross-platform KYC verification platform that enables users to complete identity verification once and reuse it across multiple services. Built with React Native, Expo Router, and NativeWind.

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Expo SDK](https://img.shields.io/badge/Expo%20SDK-53.0.22-blue.svg)](https://expo.dev/)
[![React Native](https://img.shields.io/badge/React%20Native-0.79.5-green.svg)](https://reactnative.dev/)
[![TypeScript](https://img.shields.io/badge/TypeScript-5.8.3-blue.svg)](https://www.typescriptlang.org/)

## Overview

Verza is a reusable KYC verification platform that streamlines identity verification across multiple platforms. Users complete verification once and can authenticate seamlessly across integrated services.
### Android APK File: https://drive.google.com/file/d/1LJBjNBh489vnNG-tzTEddlUnVSfEkDEs/view?usp=drivesdk

### Key Features

- **One-Time Verification** - Complete KYC once, use everywhere
- **Cross-Platform Support** - iOS, Android, and Web compatibility
- **Modern UI/UX** - Built with NativeWind (Tailwind CSS for React Native)
- **Secure Integration** - Powered by Onfido SDK for enterprise-grade verification
- **Performance Optimized** - Built with Expo Router and React Native best practices
- **Dark Mode Support** - Automatic theme switching
- **Type Safety** - Full TypeScript implementation

## Quick Start

### Prerequisites

- **Node.js** (v18 or higher)
- **npm** or **yarn** package manager
- **Expo CLI** - [Installation Guide](https://docs.expo.dev/get-started/installation/)

### Installation

1. Clone the repository:
   ```bash
   git clone https://github.com/verzaprod/Verza.git
   cd Verza
   ```

2. Install dependencies:
   ```bash
   npm install
   ```

3. Start the development server:
   ```bash
   npx expo start
   ```

4. Run on your platform:
   ```bash
   # iOS (macOS only)
   npm run ios
   
   # Android
   npm run android
   
   # Web
   npm run web
   ```

## Technology Stack

- **React Native** `0.79.5` - Cross-platform mobile development
- **Expo** `53.0.22` - React Native platform and toolchain
- **Expo Router** `5.1.5` - File-based routing system
- **TypeScript** `5.8.3` - Type-safe JavaScript
- **NativeWind** `4.0.1` - Tailwind CSS for React Native
- **Onfido SDK** `15.0.0` - Identity verification
- **React Native Reanimated** `3.17.4` - Smooth animations

## Configuration

### Environment Variables

Create a `.env.local` file:

```env
ONFIDO_API_TOKEN=your_onfido_api_token
EXPO_PUBLIC_API_URL=https://your-api-domain.com
EXPO_PUBLIC_ENVIRONMENT=development
```

### TypeScript Path Mapping

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

## KYC Integration

Verza integrates with Onfido SDK for enterprise-grade identity verification:

### Setup

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

### Features

- Document capture (passport, ID cards, driving licenses)
- Biometric verification with liveness detection
- Real-time processing and fraud detection

## Deployment

### Web Deployment

```bash
npm run deploy
```

### Mobile App Store Deployment

```bash
# Install EAS CLI
npm install -g @expo/eas-cli

# Configure and build
eas build:configure
eas build --platform all

# Submit to app stores
eas submit
```

## Platform Support

| Platform | Status | Version |
|----------|--------|---------|
| iOS | Supported | iOS 13+ |
| Android | Supported | API 21+ |
| Web | Supported | Modern browsers |

## Contributing

1. Fork the repository
2. Create a feature branch: `git checkout -b feature/name`
3. Make your changes
4. Commit: `git commit -m 'Add feature'`
5. Push: `git push origin feature/name`
6. Open a Pull Request

## License

MIT License - see the [LICENSE](LICENSE) file for details.

## Support

- [Documentation](https://docs.expo.dev/)
- [GitHub Issues](https://github.com/verzaprod/Verza/issues)
- [GitHub Discussions](https://github.com/verzaprod/Verza/discussions)

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
git remote add upstream https://github.com/verzaprod/Verza.git

# Install dependencies
npm install

# Start development server
npx expo start
```

### Code Guidelines

- Follow the existing code style
- Write meaningful commit messages
- Add tests for new features
- Update documentation as needed
- Ensure cross-platform compatibility


## 🆘 Support & Help

### Documentation
- 📚 [Expo Documentation](https://docs.expo.dev/)
- 🎨 [NativeWind Documentation](https://www.nativewind.dev/)
- 🔐 [Onfido SDK Documentation](https://documentation.onfido.com/)
- ⚡ [React Native Documentation](https://reactnative.dev/docs/getting-started)

### Community
- 💬 [GitHub Discussions](https://github.com/verzaprod/Verza/discussions)
- 🐛 [Report Issues](https://github.com/verzaprod/Verza/issues)
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

