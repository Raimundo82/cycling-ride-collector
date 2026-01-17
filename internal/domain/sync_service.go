package domain

// SyncService defines the business logic interface for syncing activities
type SyncService interface {
	SyncWeeklyActivities() error
}
