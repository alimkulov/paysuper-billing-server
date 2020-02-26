package repository

import (
	"context"
	"github.com/paysuper/paysuper-billing-server/internal/pkg"
	"github.com/paysuper/paysuper-proto/go/billingpb"
	"time"
)

// AccountingEntryRepositoryInterface is abstraction layer for working with accounting entry and representation in database.
type AccountingEntryRepositoryInterface interface {
	// Insert adds multiple accounting entries to the collection.
	MultipleInsert(context.Context, []*billingpb.AccountingEntry) error

	// GetCorrectionsForRoyaltyReport returns the account entries by merchant id, currency and dates.
	GetCorrectionsForRoyaltyReport(context.Context, string, string, time.Time, time.Time) ([]*billingpb.AccountingEntry, error)

	// GetRollingReservesForRoyaltyReport returns the account entries by merchant id, currency and dates with merchant_royalty_correction type.
	GetRollingReservesForRoyaltyReport(context.Context, string, string, []string, time.Time, time.Time) ([]*billingpb.AccountingEntry, error)

	// GetRollingReserveForBalance returns the account entries by merchant id, currency and types.
	GetRollingReserveForBalance(context.Context, string, string, []string, time.Time) ([]*pkg.ReserveQueryResItem, error)

	// GetById returns the account entry by unique identifier.
	GetById(context.Context, string) (*billingpb.AccountingEntry, error)

	// GetByTypeWithTaxes returns the account entries by type with none empty amount.
	GetByTypeWithTaxes(context.Context, []string) ([]*billingpb.AccountingEntry, error)

	// GetByObjectSource returns the account entries by object, source type and id.
	GetByObjectSource(context.Context, string, string, string) (*billingpb.AccountingEntry, error)

	// ApplyObjectSource fills the passed object with data from the search result (by object, source type and id).
	ApplyObjectSource(context.Context, string, string, string, string, *billingpb.AccountingEntry) (*billingpb.AccountingEntry, error)

	// GetDistinctBySourceId returns distinct source identities.
	GetDistinctBySourceId(context.Context) ([]string, error)

	// FindBySource returns the account entries by source type and id.
	FindBySource(context.Context, string, string) ([]*billingpb.AccountingEntry, error)

	// FindByTypeCountryDates returns the account entries by type, country and dates.
	FindByTypeCountryDates(context.Context, string, []string, time.Time, time.Time) ([]*billingpb.AccountingEntry, error)

	// BulkWrite writing account entries.
	BulkWrite(context.Context, []*billingpb.AccountingEntry) error
}
