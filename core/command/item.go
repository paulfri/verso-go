package command

import (
	"context"

	"github.com/versolabs/verso/db/query"
)

func (c *Command) MarkUnread(ctx context.Context, readerID string, userID int64) error {
	_, err := c.Queries.UpdateQueueItemReadState(
		ctx,
		query.UpdateQueueItemReadStateParams{
			ReaderID: readerID,
			UserID:   userID,
			Unread:   true,
		},
	)

	return err
}

func (c *Command) MarkRead(ctx context.Context, readerID string, userID int64) error {
	_, err := c.Queries.UpdateQueueItemReadState(
		ctx,
		query.UpdateQueueItemReadStateParams{
			ReaderID: readerID,
			UserID:   userID,
			Unread:   false,
		},
	)

	return err
}

func (c *Command) MarkStarred(ctx context.Context, readerID string, userID int64) error {
	_, err := c.Queries.UpdateQueueItemStarredState(
		ctx,
		query.UpdateQueueItemStarredStateParams{
			ReaderID: readerID,
			UserID:   userID,
			Starred:  true,
		},
	)

	return err
}

func (c *Command) MarkUnstarred(ctx context.Context, readerID string, userID int64) error {
	_, err := c.Queries.UpdateQueueItemStarredState(
		ctx,
		query.UpdateQueueItemStarredStateParams{
			ReaderID: readerID,
			UserID:   userID,
			Starred:  false,
		},
	)

	return err
}
