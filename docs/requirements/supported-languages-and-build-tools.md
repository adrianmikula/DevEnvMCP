# Supported Languages and Build Tools

This document outlines the languages and build tools that the Dev-Env Sentinel MCP will support, organized by release phases.

## Phase 1: Initial Support (MVP)

### Java Ecosystem

The Sentinel will support the following Java build tools and package managers:

#### Build Tools
- **Maven** (`pom.xml`)
  - Cache location: `~/.m2/repository`
  - Build output: `target/`
  - Dependency resolution: Maven Central, custom repositories
  
- **Gradle** (`build.gradle`, `build.gradle.kts`)
  - Cache location: `~/.gradle/caches`
  - Build output: `build/`
  - Dependency resolution: Maven Central, JCenter, custom repositories

#### Package Managers
- **Maven** (via `pom.xml`)
- **Gradle** (via `build.gradle` or `build.gradle.kts`)

#### Key Features to Support
- Verify build freshness by comparing manifest timestamps against cache artifacts
- Detect stale `.class` files in `target/` or `build/` directories
- Check for missing dependencies in local repositories
- Validate Java version compatibility

---

### npm Ecosystem

The Sentinel will support the following Node.js package managers:

#### Package Managers
- **npm** (`package.json`, `package-lock.json`)
  - Cache location: `~/.npm` (platform-dependent)
  - Build output: `node_modules/`, `dist/`, `build/`
  
- **Yarn** (`package.json`, `yarn.lock`)
  - Cache location: `~/.yarn/cache` (Yarn v2+), `~/.yarn` (Yarn v1)
  - Build output: `node_modules/`, `dist/`, `build/`
  
- **pnpm** (`package.json`, `pnpm-lock.yaml`)
  - Cache location: `~/.pnpm-store`
  - Build output: `node_modules/`, `dist/`, `build/`

#### Key Features to Support
- Verify package lock file consistency with `node_modules/`
- Detect stale build artifacts in `dist/` or `build/` directories
- Check for missing or corrupted dependencies
- Validate Node.js version compatibility (via `engines` field in `package.json`)

---

## Phase 2: Future Support (Planned)

### Python Ecosystem

#### Package Managers
- **pip** (`requirements.txt`, `setup.py`, `pyproject.toml`)
- **pipenv** (`Pipfile`, `Pipfile.lock`)
- **poetry** (`pyproject.toml`, `poetry.lock`)
- **conda** (`environment.yml`, `conda-lock.yml`)

#### Build Tools
- **setuptools** (via `setup.py`)
- **poetry** (build system)
- **pip** (build backend)

#### Cache Locations
- `~/.cache/pip`
- `~/.local/share/pip`
- `~/.cache/pypoetry`
- `~/anaconda3/pkgs` or `~/miniconda3/pkgs`

---

### C++ Ecosystem

#### Build Tools
- **CMake** (`CMakeLists.txt`)
- **Make** (`Makefile`)
- **Bazel** (`BUILD`, `WORKSPACE`)
- **Ninja** (via CMake or standalone)

#### Package Managers
- **vcpkg** (`vcpkg.json`)
- **Conan** (`conanfile.txt`, `conanfile.py`)
- **Hunter** (CMake-based)

#### Cache Locations
- `build/`, `cmake-build-*/`
- `~/.conan`
- `vcpkg_installed/`

---

### C# / .NET Ecosystem

#### Build Tools
- **MSBuild** (`.csproj`, `.sln`)
- **dotnet CLI** (`*.csproj`, `*.sln`)

#### Package Managers
- **NuGet** (`packages.config`, `*.csproj` PackageReference)
- **dotnet CLI** (via `*.csproj`)

#### Cache Locations
- `~/.nuget/packages`
- `bin/`, `obj/`
- `packages/` (legacy)

---

### Additional Languages (Future Consideration)

- **Rust** (Cargo, `Cargo.toml`, `Cargo.lock`)
- **Go** (`go.mod`, `go.sum`)
- **Ruby** (Bundler, `Gemfile`, `Gemfile.lock`)
- **PHP** (Composer, `composer.json`, `composer.lock`)
- **Swift** (Swift Package Manager, `Package.swift`)
- **Kotlin** (Gradle, Maven)
- **Scala** (sbt, Maven)

---

## Implementation Notes

### Detection Strategy

The Sentinel will detect supported ecosystems by:
1. **File-based detection**: Looking for characteristic manifest files (e.g., `pom.xml`, `package.json`)
2. **Directory structure**: Recognizing standard build output directories
3. **Cache location discovery**: Checking common cache locations for each tool

### Build Freshness Verification

For each supported ecosystem, the Sentinel will:
- Compare manifest file modification times with cache/artifact timestamps
- Identify stale build artifacts that may cause runtime issues
- Report discrepancies between declared dependencies and cached artifacts

### Environment Parity Checks

The Sentinel will verify:
- Required runtime versions (Java version, Node.js version, Python version, etc.)
- Service dependencies (Docker containers, databases, message queues)
- Environment variables referenced in code vs. actual environment

---

## Version History

- **v1.0** (Phase 1): Java (Maven, Gradle) and npm (npm, Yarn, pnpm)
- **v2.0** (Phase 2): Python, C++, C# support (TBD)

