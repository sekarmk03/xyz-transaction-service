package authorization

type AccessibleRoles map[string]map[string][]uint32

const (
	BasePath       = "xyz-transaction-service"
	TransactionSvc = "TransactionService"
)

var roles = AccessibleRoles{
	"/" + BasePath + "." + TransactionSvc + "/": {
		// "DeletePost":  {1, 2, 8},
	},
}

func GetAccessibleRoles() map[string][]uint32 {
	routes := make(map[string][]uint32)

	for service, methods := range roles {
		for method, methodRoles := range methods {
			route := service + method
			routes[route] = methodRoles
		}
	}

	return routes
}
