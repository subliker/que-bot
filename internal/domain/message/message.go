package message

type Manager interface {
	// Text returns localization text
	Text(textName string) string
}

type manager struct {
}

func NewManager() Manager {
	return &manager{}
}

func (m *manager) Text(textName string) string {
	// TODO
	return textName
}
