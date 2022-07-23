package api

import (
	"github.com/go-chi/chi/v5"
	"net/http"
)

// TransactionTypeQueryArg Possible values: "deposit" and "withdrawal"
const TransactionTypeQueryArg = "transaction_type"
const TransactionTypeDeposit = "deposit"
const TransactionTypeWithdrawal = "withdrawal"

func NewAPIRouter(h Handlers, authMiddleware Middleware) http.Handler {
	r := chi.NewRouter()

	r.Route("/api/v1", func(r chi.Router) {
		r.Route("/atm-transaction", func(r chi.Router) {
			// requires a TransactionTypeQueryArg
			// Response: CodeResponse
			// Throws: 401, TransactionTypeNotProvided
			r.Get("/gen-code", authMiddleware(h.GenCode).ServeHTTP)

			// Request: CodeRequest; requires a TransactionTypeQueryArg
			// Response: VerifiedCodeResponse
			// Throws: TransactionTypeNotProvided, InvalidCode, InvalidTransactionType
			r.Post("/verify-code", h.VerifyCode)

			// Request: BanknoteCheckRequest
			// Response: 200 if accepted, client error (or 500) if rejected
			// Throws: InvalidCode (means session-expired)
			r.Post("/check-banknote", h.CheckBanknote)

			// Request: FinalizeTransactionRequest
			// Response: 200 if accepted, client error (or 500) if rejected
			// Throws: InvalidATMSecret, InsufficientFunds
			r.Post("/finalize-transaction", h.FinalizeTransaction)
		})
		// Response: UserInfoResponse 
		// Throws: 401
		r.Get("/user-info", authMiddleware(h.GetUserInfo).ServeHTTP)
	})
	

	return r
}

type Middleware = func(http.Handler) http.Handler
