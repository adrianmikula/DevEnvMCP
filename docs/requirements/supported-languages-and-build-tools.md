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

## Phase 2: Additional Language Support (Implemented)

### React Ecosystem

#### Framework
- **React** (via `package.json` with React dependencies)
  - Detection: `package.json` with React dependencies, `src/` directory
  - Build output: `dist/`, `build/`, `.next/`, `out/`
  - Environment variables: `REACT_APP_*` prefix
  - Cache location: npm/yarn/pnpm caches

#### Key Features
- Verify React build freshness
- Check source files vs. build artifacts
- Environment variable detection for React apps

---

### Vite Ecosystem

#### Build Tool
- **Vite** (`vite.config.js`, `vite.config.ts`)
  - Detection: `vite.config.*` files, `package.json` with Vite
  - Build output: `dist/`, `build/`
  - Cache location: npm/yarn/pnpm caches, `node_modules/.vite`
  - Environment variables: `VITE_*` prefix, `import.meta.env.*`

#### Key Features
- Verify Vite config vs. build output
- Check source files vs. Vite build artifacts
- Fast refresh and HMR support detection

---

### Python Ecosystem

#### Package Managers
- **pip** (`requirements.txt`, `setup.py`, `pyproject.toml`)
- **pipenv** (`Pipfile`, `Pipfile.lock`)
- **poetry** (`pyproject.toml`, `poetry.lock`)

#### Build Tools
- **setuptools** (via `setup.py`)
- **poetry** (build system)
- **pip** (build backend)

#### Cache Locations
- `~/.cache/pip`
- `~/.local/share/pip`
- `~/.cache/poetry`
- `~/anaconda3/pkgs` or `~/miniconda3/pkgs`

#### Key Features
- Verify Python source vs. compiled bytecode (`__pycache__`)
- Check requirements.txt vs. pip cache
- Detect stale build artifacts

---

### Poetry Ecosystem

#### Package Manager
- **Poetry** (`pyproject.toml`, `poetry.lock`)
  - Detection: `pyproject.toml` with Poetry configuration
  - Build output: `dist/`
  - Cache location: `~/.cache/pypoetry`
  - Lock file: `poetry.lock`

#### Key Features
- Verify `pyproject.toml` vs. `poetry.lock` freshness
- Check Poetry cache vs. dependencies
- Poetry-specific build verification

---

### Conda Ecosystem

#### Package Manager
- **Conda** (`environment.yml`, `conda-lock.yml`)
  - Detection: `environment.yml` or `conda-lock.yml`
  - Cache location: `~/.conda/pkgs`, `~/anaconda3/pkgs`
  - Environment management

#### Key Features
- Verify environment.yml vs. conda cache
- Check conda environment consistency
- Conda package build verification

---

### Docker Ecosystem

#### Containerization
- **Docker** (`Dockerfile`, `docker-compose.yml`)
  - Detection: `Dockerfile`, `docker-compose.*` files
  - Cache location: `~/.docker`, `/var/lib/docker`
  - Environment variables: `ENV`, `ARG` in Dockerfile

#### Key Features
- Verify Dockerfile vs. built images
- Check docker-compose.yml vs. running containers
- Docker service health checks

---

### PostgreSQL Ecosystem

#### Database
- **PostgreSQL** (via Docker or native)
  - Detection: `docker-compose.yml` with Postgres, `*.sql` files, `migrations/`
  - Environment variables: `POSTGRES_*`, `DATABASE_*`, `DB_*`
  - Required vars: `POSTGRES_HOST`, `POSTGRES_DB`, `POSTGRES_USER`, `POSTGRES_PASSWORD`

#### Key Features
- Verify migration files vs. database schema
- Check database connection configuration
- Detect missing migrations

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

#### Key Features
- Verify .csproj vs. compiled DLLs
- Check C# source files vs. build output
- NuGet package cache verification

---

### Webpack Ecosystem

#### Build Tool
- **Webpack** (`webpack.config.js`, `webpack.config.ts`)
  - Detection: `webpack.config.*` files
  - Build output: `dist/`, `build/`
  - Cache location: `node_modules/.cache/webpack`

#### Key Features
- Verify webpack config vs. build output
- Check source files vs. webpack bundles
- Webpack cache verification

---

### Rollup Ecosystem

#### Build Tool
- **Rollup** (`rollup.config.js`, `rollup.config.ts`)
  - Detection: `rollup.config.*` files
  - Build output: `dist/`, `build/`
  - Cache location: `node_modules/.cache/rollup`

#### Key Features
- Verify rollup config vs. build output
- Check source files vs. rollup bundles

---

### Sass/SCSS Ecosystem

#### Preprocessor
- **Sass/SCSS** (`*.scss`, `*.sass`)
  - Detection: `*.scss` or `*.sass` files
  - Build output: `dist/css/`, `build/css/`, `public/css/`
  - Compiles to CSS

#### Key Features
- Verify SCSS/Sass source vs. compiled CSS
- Check for stale CSS files

---

### Spring Framework Ecosystem

#### Framework
- **Spring Framework** (Spring Boot, Spring MVC)
  - Detection: `application.properties`, `application.yml`, Spring dependencies in `pom.xml`/`build.gradle`
  - Build output: `target/*.jar`, `target/*.war`
  - Configuration: `application.properties`, `application.yml`

#### Key Features
- Verify Spring config files vs. compiled classes
- Check Java source vs. build artifacts
- Spring-specific environment variable detection

---

### Apache Tomcat Ecosystem

#### Application Server
- **Apache Tomcat** (Servlet container)
  - Detection: `web.xml`, `context.xml`, `server.xml`, WAR files
  - Build output: `target/*.war`
  - Configuration: `web.xml`, `context.xml`, `server.xml`

#### Key Features
- Verify web.xml vs. WAR file freshness
- Check Java source vs. WAR artifacts
- Tomcat-specific environment variables

---

### JBoss/WildFly Ecosystem

#### Application Server
- **JBoss/WildFly** (Enterprise application server)
  - Detection: `jboss-web.xml`, `jboss-deployment-structure.xml`, `standalone.xml`
  - Build output: `target/*.war`, `target/*.ear`
  - Configuration: `jboss-web.xml`, `standalone.xml`, `domain.xml`

#### Key Features
- Verify JBoss config files vs. WAR/EAR files
- Check Java source vs. deployment artifacts
- JBoss/WildFly-specific environment variables

---

## Phase 3: Future Support (Planned)

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
- **v2.0** (Phase 2): React, Vite, Python, Docker, PostgreSQL, C# (.NET) support
- **v2.1** (Phase 2 Extended): Poetry, Conda, Webpack, Rollup, Sass, Spring, Tomcat, JBoss support
- **v3.0** (Phase 3): C++, Rust, Go, Ruby, PHP, Swift, Kotlin, Scala (Planned)

