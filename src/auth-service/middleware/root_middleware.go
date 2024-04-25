package middleware

type RootMiddleware struct {
	TransactionMiddleware *TransactionMiddleware
	AuthMiddleware        *AuthMiddleware
}

func NewRootMiddleware(
	transactionMiddleware *TransactionMiddleware,
	authMiddleware *AuthMiddleware,
) *RootMiddleware {
	rootMiddleware := &RootMiddleware{
		TransactionMiddleware: transactionMiddleware,
		AuthMiddleware:        authMiddleware,
	}
	return rootMiddleware
}
