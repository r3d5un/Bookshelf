package orchestrator

import "context"

func (m *Module) helloWorld(ctx context.Context) error {
	m.logger.Info("Hello, World!")

	return nil
}
