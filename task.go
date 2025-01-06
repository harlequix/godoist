package godoist

type Task interface{
	GetChildren() *[]Task
	AddLabel(label string)
	RemoveLabel(label string) error
	Update(key string, value interface{}) error
}