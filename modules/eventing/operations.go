package eventing

// QueueEvent adds a new event to the queue
// func (m *Module) QueueEvent(ctx context.Context, event *model.AddEventRequest) error {
// 	m.lock.RLock()
// 	defer m.lock.RUnlock()

// 	if !m.config.Enabled {
// 		return errors.New("Eventing module is disabled")
// 	}

// 	req := m.generateQueueEventRequest(event)

// 	if err := m.crud.Create(ctx, m.config.DBType, m.project, m.config.Col, req); err != nil {
// 		return err
// 	}

// 	return nil
// }
