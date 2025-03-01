package store

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/bytebase/bytebase/api"
	"go.uber.org/zap"
)

var (
	_ api.IssueSubscriberService = (*IssueSubscriberService)(nil)
)

// IssueSubscriberService represents a service for managing issueSubscriber.
type IssueSubscriberService struct {
	l  *zap.Logger
	db *DB
}

// NewIssueSubscriberService returns a new instance of IssueSubscriberService.
func NewIssueSubscriberService(logger *zap.Logger, db *DB) *IssueSubscriberService {
	return &IssueSubscriberService{l: logger, db: db}
}

// CreateIssueSubscriber creates a new issueSubscriber.
func (s *IssueSubscriberService) CreateIssueSubscriber(ctx context.Context, create *api.IssueSubscriberCreate) (*api.IssueSubscriber, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, FormatError(err)
	}
	defer tx.PTx.Rollback()

	issueSubscriber, err := createIssueSubscriber(ctx, tx.PTx, create)
	if err != nil {
		return nil, err
	}

	if err := tx.PTx.Commit(); err != nil {
		return nil, FormatError(err)
	}

	return issueSubscriber, nil
}

// FindIssueSubscriberList retrieves a list of issueSubscribers based on find.
func (s *IssueSubscriberService) FindIssueSubscriberList(ctx context.Context, find *api.IssueSubscriberFind) ([]*api.IssueSubscriber, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, FormatError(err)
	}
	defer tx.PTx.Rollback()

	list, err := findIssueSubscriberList(ctx, tx.PTx, find)
	if err != nil {
		return nil, err
	}

	return list, nil
}

// DeleteIssueSubscriber deletes an existing issueSubscriber by ID.
func (s *IssueSubscriberService) DeleteIssueSubscriber(ctx context.Context, delete *api.IssueSubscriberDelete) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return FormatError(err)
	}
	defer tx.PTx.Rollback()

	if err := deleteIssueSubscriber(ctx, tx.PTx, delete); err != nil {
		return FormatError(err)
	}

	if err := tx.PTx.Commit(); err != nil {
		return FormatError(err)
	}

	return nil
}

// createIssueSubscriber creates a new issueSubscriber.
func createIssueSubscriber(ctx context.Context, tx *sql.Tx, create *api.IssueSubscriberCreate) (*api.IssueSubscriber, error) {
	// Insert row into database.
	row, err := tx.QueryContext(ctx, `
		INSERT INTO issue_subscriber (
			issue_id,
			subscriber_id
		)
		VALUES ($1, $2)
		RETURNING issue_id, subscriber_id
	`,
		create.IssueID,
		create.SubscriberID,
	)

	if err != nil {
		return nil, FormatError(err)
	}
	defer row.Close()

	row.Next()
	var issueSubscriber api.IssueSubscriber
	if err := row.Scan(
		&issueSubscriber.IssueID,
		&issueSubscriber.SubscriberID,
	); err != nil {
		return nil, FormatError(err)
	}

	return &issueSubscriber, nil
}

func findIssueSubscriberList(ctx context.Context, tx *sql.Tx, find *api.IssueSubscriberFind) ([]*api.IssueSubscriber, error) {
	// Build WHERE clause.
	where, args := []string{"1 = 1"}, []interface{}{}
	if v := find.IssueID; v != nil {
		where, args = append(where, fmt.Sprintf("issue_id = $%d", len(args)+1)), append(args, *v)
	}
	if v := find.SubscriberID; v != nil {
		where, args = append(where, fmt.Sprintf("subscriber_id = $%d", len(args)+1)), append(args, *v)
	}

	rows, err := tx.QueryContext(ctx, `
		SELECT
			issue_id,
			subscriber_id
		FROM issue_subscriber
		WHERE `+strings.Join(where, " AND "),
		args...,
	)
	if err != nil {
		return nil, FormatError(err)
	}
	defer rows.Close()

	// Iterate over result set and deserialize rows into issueSubscriberList.
	var issueSubscriberList []*api.IssueSubscriber
	for rows.Next() {
		var issueSubscriber api.IssueSubscriber
		if err := rows.Scan(
			&issueSubscriber.IssueID,
			&issueSubscriber.SubscriberID,
		); err != nil {
			return nil, FormatError(err)
		}

		issueSubscriberList = append(issueSubscriberList, &issueSubscriber)
	}
	if err := rows.Err(); err != nil {
		return nil, FormatError(err)
	}

	return issueSubscriberList, nil
}

// deleteIssueSubscriber permanently deletes a issueSubscriber by ID.
func deleteIssueSubscriber(ctx context.Context, tx *sql.Tx, delete *api.IssueSubscriberDelete) error {
	// Remove row from database.
	if _, err := tx.ExecContext(ctx, `DELETE FROM issue_subscriber WHERE issue_id = $1 AND subscriber_id = $2`, delete.IssueID, delete.SubscriberID); err != nil {
		return FormatError(err)
	}
	return nil
}
