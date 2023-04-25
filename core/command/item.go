package command

import (
	"context"

	"github.com/versolabs/verso/db/query"
)

func (c *Command) MarkUnread(ctx context.Context, itemID int64, userID int64) error {
	_, err := c.Queries.UpdateQueueItemReadState(
		ctx,
		query.UpdateQueueItemReadStateParams{
			RSSItemID: itemID,
			UserID:    userID,
			Unread:    true,
		},
	)

	return err
}

func (c *Command) MarkRead(ctx context.Context, itemID int64, userID int64) error {
	_, err := c.Queries.UpdateQueueItemReadState(
		ctx,
		query.UpdateQueueItemReadStateParams{
			RSSItemID: itemID,
			UserID:    userID,
			Unread:    false,
		},
	)

	return err
}

func (c *Command) MarkStarred(ctx context.Context, itemID int64, userID int64) error {
	_, err := c.Queries.UpdateQueueItemStarredState(
		ctx,
		query.UpdateQueueItemStarredStateParams{
			RSSItemID: itemID,
			UserID:    userID,
			Starred:   true,
		},
	)

	return err
}

func (c *Command) MarkUnstarred(ctx context.Context, itemID int64, userID int64) error {
	_, err := c.Queries.UpdateQueueItemStarredState(
		ctx,
		query.UpdateQueueItemStarredStateParams{
			RSSItemID: itemID,
			UserID:    userID,
			Starred:   false,
		},
	)

	return err
}
