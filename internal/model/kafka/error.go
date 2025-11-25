package kafka

type Error struct {
	Err error
}

type Status struct {
	Status string
}

type Success struct {
	Msg string
}
