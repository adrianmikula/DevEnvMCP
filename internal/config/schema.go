package config

// EcosystemConfig represents the complete ecosystem configuration
type EcosystemConfig struct {
	Ecosystem Ecosystem `yaml:"ecosystem"`
}

// Ecosystem defines an ecosystem (language/tool combination)
type Ecosystem struct {
	Name    string `yaml:"name"`
	ID      string `yaml:"id"`
	Version string `yaml:"version"`
	
	Detection      Detection      `yaml:"detection"`
	Manifest       Manifest       `yaml:"manifest"`
	Cache          Cache          `yaml:"cache"`
	Build          Build          `yaml:"build"`
	Dependencies   Dependencies   `yaml:"dependencies"`
	Verification   Verification   `yaml:"verification"`
	Environment    Environment    `yaml:"environment"`
	Infrastructure Infrastructure `yaml:"infrastructure"`
	Reconciliation Reconciliation `yaml:"reconciliation"`
	VersionConfig  VersionConfig  `yaml:"version"`
	Requirements   Requirements   `yaml:"requirements"`
}

// Detection defines how to detect this ecosystem
type Detection struct {
	ManifestFiles     []string `yaml:"manifest_files"`
	RequiredFiles     []string `yaml:"required_files"`
	OptionalFiles     []string `yaml:"optional_files"`
	DirectoryPatterns []string `yaml:"directory_patterns"`
}

// Manifest defines the manifest file
type Manifest struct {
	PrimaryFile string `yaml:"primary_file"`
	Location    string `yaml:"location"`
	Format      string `yaml:"format"`
}

// Cache defines cache locations
type Cache struct {
	Locations      []string `yaml:"locations"`
	Structure      string   `yaml:"structure"`
	ArtifactPattern string  `yaml:"artifact_pattern"`
}

// Build defines build output
type Build struct {
	OutputDirectories []string `yaml:"output_directories"`
	ArtifactPatterns  []string `yaml:"artifact_patterns"`
	CleanCommand      string   `yaml:"clean_command"`
}

// Dependencies defines dependency management
type Dependencies struct {
	LockFile      string `yaml:"lock_file"`
	LockFileFormat string `yaml:"lock_file_format"`
	ResolveCommand string `yaml:"resolve_command"`
	CheckCommand   string `yaml:"check_command"`
}

// Verification defines verification commands
type Verification struct {
	BuildFreshness BuildFreshness `yaml:"build_freshness"`
	DependencyAudit DependencyAudit `yaml:"dependency_audit"`
}

// BuildFreshness defines build freshness checks
type BuildFreshness struct {
	ManifestTimestampCheck bool              `yaml:"manifest_timestamp_check"`
	CacheTimestampCheck    bool              `yaml:"cache_timestamp_check"`
	BuildOutputCheck       bool              `yaml:"build_output_check"`
	Commands               []VerificationCommand `yaml:"commands"`
}

// DependencyAudit defines dependency audit checks
type DependencyAudit struct {
	Enabled  bool                `yaml:"enabled"`
	Commands []VerificationCommand `yaml:"commands"`
}

// VerificationCommand defines a single verification command
type VerificationCommand struct {
	Name        string `yaml:"name"`
	Type        string `yaml:"type"`
	Source      string `yaml:"source,omitempty"`
	Target      string `yaml:"target,omitempty"`
	TargetPattern string `yaml:"target_pattern,omitempty"`
	Command     string `yaml:"command,omitempty"`
	Description string `yaml:"description"`
}

// Environment defines environment variable handling
type Environment struct {
	VariablePatterns []string `yaml:"variable_patterns"`
	ConfigFiles      []string `yaml:"config_files"`
	RequiredVars     []string `yaml:"required_vars"`
}

// Infrastructure defines infrastructure requirements
type Infrastructure struct {
	Services []Service `yaml:"services"`
}

// Service defines a service requirement
type Service struct {
	Name           string `yaml:"name"`
	Type           string `yaml:"type"`
	CheckCommand   string `yaml:"check_command"`
	VersionExtract string `yaml:"version_extract"`
}

// Reconciliation defines auto-fix commands
type Reconciliation struct {
	Fixes []Fix `yaml:"fixes"`
}

// Fix defines a fix command
type Fix struct {
	IssueType     string `yaml:"issue_type"`
	Command       string `yaml:"command"`
	VerifyCommand string `yaml:"verify_command"`
	Description   string `yaml:"description"`
}

// VersionConfig defines version management configuration
type VersionConfig struct {
	Language          string   `yaml:"language"`
	VersionCommand    string   `yaml:"version_command"`
	VersionPattern    string   `yaml:"version_pattern"`
	RuntimePattern    string   `yaml:"runtime_pattern,omitempty"` // For Java and similar
	VersionManagers   []VersionManager `yaml:"version_managers"`
	RuntimeVariants   []RuntimeVariant `yaml:"runtime_variants,omitempty"` // For Java
}

// VersionManager defines a version management tool
type VersionManager struct {
	Name         string `yaml:"name"`
	CheckCommand string `yaml:"check_command"`
	ListCommand  string `yaml:"list_command"`
	InstallCommand string `yaml:"install_command"` // Template: "install {version}"
	SwitchCommand  string `yaml:"switch_command"`  // Template: "use {version}"
	CurrentCommand string `yaml:"current_command,omitempty"`
}

// RuntimeVariant defines a runtime variant (e.g., Java runtimes)
type RuntimeVariant struct {
	Name        string   `yaml:"name"`
	Provider    string   `yaml:"provider"`
	Pattern     string   `yaml:"pattern"` // Regex to identify this variant
	Compatible  bool     `yaml:"compatible"` // Generally compatible
	Description string   `yaml:"description,omitempty"`
}

// Requirements defines version requirements
type Requirements struct {
	MinVersion       string   `yaml:"min_version,omitempty"`
	MaxVersion       string   `yaml:"max_version,omitempty"`
	PreferredVersions []string `yaml:"preferred_versions,omitempty"`
	PreferredRuntimes []string `yaml:"preferred_runtimes,omitempty"` // For Java
	ExcludedVersions  []string `yaml:"excluded_versions,omitempty"`
	ExcludedRuntimes  []string `yaml:"excluded_runtimes,omitempty"` // For Java
}

