package service_type

const (
	RABBITMQ = "RABBITMQ"
)

func GetAll() []string {
	return []string{
		RABBITMQ,
	}
}
