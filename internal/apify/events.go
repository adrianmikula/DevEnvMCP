package apify

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// EventType represents a billable event type
type EventType string

const (
	// Free tier events (no charge)
	EventVerifyBuildFreshness    EventType = "verify_build_freshness"
	EventCheckInfrastructure     EventType = "check_infrastructure_parity"
	EventEnvVarAudit             EventType = "env_var_audit"
	EventCheckLicenseStatus      EventType = "check_license_status"
	EventGetProLicense           EventType = "get_pro_license"

	// Premium tier events (billable)
	EventReconcileEnvironment    EventType = "reconcile_environment"    // $0.05
	EventAutoFix                 EventType = "auto_fix"                  // $0.05
	EventAdvancedDiagnostics     EventType = "advanced_diagnostics"      // $0.03
	EventDockerOrchestration     EventType = "docker_orchestration"      // $0.10
	EventCustomConfigs           EventType = "custom_configs"            // $0.02
)

// Event represents a billable event
type Event struct {
	Type        EventType `json:"type"`
	ToolName    string    `json:"tool_name"`
	Price       float64   `json:"price"`
	Timestamp   string    `json:"timestamp"`
	UserID      string    `json:"user_id,omitempty"`
	ProjectRoot string    `json:"project_root,omitempty"`
}

// EventTracker tracks and reports billable events to Apify
type EventTracker struct {
	enabled bool
	apiKey  string
}

// NewEventTracker creates a new event tracker
func NewEventTracker() *EventTracker {
	// Check if running in Apify environment
	apiKey := os.Getenv("APIFY_API_TOKEN")
	enabled := apiKey != "" && os.Getenv("APIFY_ACTOR_ID") != ""

	return &EventTracker{
		enabled: enabled,
		apiKey:  apiKey,
	}
}

// TrackEvent tracks a billable event
func (et *EventTracker) TrackEvent(eventType EventType, toolName string, metadata map[string]string) error {
	if !et.enabled {
		// Not in Apify environment, skip tracking
		return nil
	}

	price := GetEventPrice(eventType)
	if price == 0 {
		// Free event, no need to track for billing
		return nil
	}

	event := Event{
		Type:      eventType,
		ToolName:  toolName,
		Price:     price,
		Timestamp: getCurrentTimestamp(),
	}

	if metadata != nil {
		if userID, ok := metadata["user_id"]; ok {
			event.UserID = userID
		}
		if projectRoot, ok := metadata["project_root"]; ok {
			event.ProjectRoot = projectRoot
		}
	}

	// Log event (Apify will pick this up for billing)
	return et.logEvent(event)
}

// logEvent logs the event for Apify billing
func (et *EventTracker) logEvent(event Event) error {
	// In Apify, events are tracked via:
	// 1. Apify's built-in billing system (automatic for Actor runs)
	// 2. Custom event logging (for detailed tracking)
	
	// Log to stdout/stderr in JSON format for Apify to parse
	eventJSON, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal event: %w", err)
	}

	// Apify reads structured logs for billing
	fmt.Fprintf(os.Stderr, "APIFY_EVENT:%s\n", string(eventJSON))
	
	return nil
}

// GetEventPrice returns the price for an event type
func GetEventPrice(eventType EventType) float64 {
	prices := map[EventType]float64{
		// Free tier - no charge
		EventVerifyBuildFreshness:    0.00,
		EventCheckInfrastructure:     0.00,
		EventEnvVarAudit:             0.00,
		EventCheckLicenseStatus:      0.00,
		EventGetProLicense:           0.00,

		// Premium tier - billable
		EventReconcileEnvironment:    0.05, // Auto-fix is high value
		EventAutoFix:                 0.05,
		EventAdvancedDiagnostics:     0.03, // Diagnostics are medium value
		EventDockerOrchestration:     0.10, // Docker ops are high compute
		EventCustomConfigs:           0.02, // Config operations are low value
	}

	if price, ok := prices[eventType]; ok {
		return price
	}
	return 0.00 // Unknown events are free
}

// IsBillableEvent checks if an event type is billable
func IsBillableEvent(eventType EventType) bool {
	return GetEventPrice(eventType) > 0
}

// GetEventDescription returns a human-readable description
func GetEventDescription(eventType EventType) string {
	descriptions := map[EventType]string{
		EventVerifyBuildFreshness:    "Verify build artifact freshness",
		EventCheckInfrastructure:     "Check infrastructure service parity",
		EventEnvVarAudit:             "Audit environment variables",
		EventCheckLicenseStatus:      "Check license status",
		EventGetProLicense:           "Get Pro license information",
		EventReconcileEnvironment:    "Auto-fix environment issues (Premium)",
		EventAutoFix:                 "Automatic issue resolution (Premium)",
		EventAdvancedDiagnostics:     "Advanced diagnostic analysis (Premium)",
		EventDockerOrchestration:     "Docker container orchestration (Enterprise)",
		EventCustomConfigs:           "Custom configuration management (Enterprise)",
	}

	if desc, ok := descriptions[eventType]; ok {
		return desc
	}
	return string(eventType)
}

// getCurrentTimestamp returns current timestamp in ISO format
func getCurrentTimestamp() string {
	return time.Now().Format(time.RFC3339)
}

