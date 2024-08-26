package orchestrator

import "context"

func (m *Module) helloWorld(ctx context.Context) error {
	m.logger.InfoContext(ctx, "Hello, World!")

	return nil
}
