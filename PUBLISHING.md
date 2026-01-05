# Publishing Guide

This guide explains how to build and publish the `dev-env-sentinel` package to npm.

## Prerequisites

1. **Node.js 14+** - Required for build scripts
2. **Go 1.13+** - Required to compile binaries
3. **npm account** - With access to publish the package

## Build Process

The build process compiles Go binaries for multiple platforms:

- Windows (amd64)
- Windows (arm64)
- Linux (amd64)
- Linux (arm64)

### Local Build

Build for all platforms:
```bash
npm run build:all
```

Build for current platform only:
```bash
npm run build
```

Build for specific platform:
```bash
npm run build:windows
npm run build:linux
```

### Automated Build on Publish

The `prepublishOnly` script automatically builds all platforms before publishing:
```bash
npm publish
```

This ensures the published package includes binaries for all supported platforms.

## Publishing

### First Time Setup

1. **Login to npm**:
   ```bash
   npm login
   ```

2. **Verify package name availability**:
   ```bash
   npm view dev-env-sentinel
   ```
   If the package doesn't exist, you're good to go!

### Publishing Steps

1. **Update version** in `package.json`:
   ```json
   "version": "0.1.0"
   ```

2. **Test the build**:
   ```bash
   npm run build:all
   ```

3. **Test installation locally**:
   ```bash
   npm pack
   npm install -g dev-env-sentinel-0.1.0.tgz
   ```

4. **Publish to npm**:
   ```bash
   npm publish
   ```

   Or publish as a beta/alpha:
   ```bash
   npm publish --tag beta
   npm publish --tag alpha
   ```

### Post-Publish Verification

1. **Verify package on npm**:
   ```bash
   npm view dev-env-sentinel
   ```

2. **Test installation**:
   ```bash
   npm install -g dev-env-sentinel
   dev-env-sentinel --version
   ```

3. **Test with npx**:
   ```bash
   npx dev-env-sentinel
   ```

## Package Contents

The published package includes:

- `bin/` - Platform-specific binaries and wrapper script
- `language-configs/` - Language-level ecosystem configurations
- `tool-configs/` - Tool-specific configurations
- `scripts/` - Build and postinstall scripts
- `README.md` - Documentation
- `LICENSE` - MIT License

## Version Management

Follow [Semantic Versioning](https://semver.org/):

- **MAJOR** (1.0.0) - Breaking changes
- **MINOR** (0.1.0) - New features, backward compatible
- **PATCH** (0.0.1) - Bug fixes, backward compatible

Update version before publishing:
```bash
npm version patch  # 0.1.0 -> 0.1.1
npm version minor  # 0.1.0 -> 0.2.0
npm version major  # 0.1.0 -> 1.0.0
```

## Troubleshooting

### Build Fails

- Ensure Go is installed and in PATH: `go version`
- Ensure Node.js is installed: `node --version`
- Check that all dependencies are available

### Binary Not Found After Install

- Run `npm run build:all` to ensure binaries are built
- Check that `bin/` directory contains platform-specific binaries
- Verify `postinstall.js` runs correctly

### Platform Not Supported

The package currently supports:
- Windows (amd64, arm64)
- Linux (amd64, arm64)

To add more platforms, update `scripts/build.js` with additional platform configurations.

