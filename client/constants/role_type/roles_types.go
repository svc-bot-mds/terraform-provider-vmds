package role_type

import "fmt"

const (
	MDS      = "MDS"
	RABBITMQ = "RABBITMQ"
	POSTGRES = "POSTGRES"
	MYSQL    = "MYSQL"
	REDIS    = "REDIS"
)

func ValidateRoleType(stateType string) error {
	switch stateType {
	case MYSQL, RABBITMQ, POSTGRES, REDIS:
		return nil
	default:
		return fmt.Errorf("invalid type: supported types are [%s, %s, %s, %s]",
			MYSQL,
			RABBITMQ,
			POSTGRES,
			REDIS)
	}
}
