package bot

import "context"

// Controller is interface to control bot
type Controller interface {
	// Run starts bot until context will end
	Run(ctx context.Context)
}
