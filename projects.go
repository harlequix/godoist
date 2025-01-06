package godoist

type Project interface{
	GetTasks() []*Task
	GetChildren() []*Project
}