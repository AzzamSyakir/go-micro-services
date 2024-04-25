package middleware

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"net/http"
	"net/http/httputil"
	"social-media/internal/config"
	"social-media/internal/model"
	"social-media/internal/model/response"

	"github.com/cockroachdb/cockroach-go/v2/crdb"
)

type TransactionMiddleware struct {
	DatabaseConfig *config.DatabaseConfig
}

func NewTransactionMiddleware(
	databaseConfig *config.DatabaseConfig,
) *TransactionMiddleware {
	transactionMiddleware := &TransactionMiddleware{
		DatabaseConfig: databaseConfig,
	}
	return transactionMiddleware
}

func (transactionMiddleware *TransactionMiddleware) GetMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, reader *http.Request) {
		maxRetries := 3
		retryCount := 0
		var result *response.Response[any]

		readerDump, readerDumpErr := httputil.DumpRequest(reader, true)
		if readerDumpErr != nil {
			result = &response.Response[any]{
				Code:    http.StatusInternalServerError,
				Message: "TransactionMiddleware GetMiddleware dump request error.",
				Data:    nil,
			}
			response.NewResponse(writer, result)
			return
		}

		executeErr := crdb.Execute(func() (err error) {
			begin, err := transactionMiddleware.DatabaseConfig.CockroachdbDatabase.Connection.Begin()
			if err != nil {
				result = &response.Response[any]{
					Code:    http.StatusInternalServerError,
					Message: "TransactionMiddleware GetMiddleware begin error.",
					Data:    nil,
				}
				return err
			}

			newReader, err := http.ReadRequest(bufio.NewReader(bytes.NewBuffer(readerDump)))
			if err != nil {
				result = &response.Response[any]{
					Code:    http.StatusInternalServerError,
					Message: "TransactionMiddleware GetMiddleware read request error.",
					Data:    nil,
				}
				return err
			}

			transaction := &model.Transaction{
				Tx:    begin,
				TxErr: nil,
			}
			ctx := context.WithValue(reader.Context(), "transaction", transaction)
			newReaderWithContext := newReader.WithContext(ctx)
			next.ServeHTTP(writer, newReaderWithContext)

			err = transaction.TxErr
			if err != nil {
				result = &response.Response[any]{
					Code:    http.StatusInternalServerError,
					Message: "TransactionMiddleware GetMiddleware transaction error, " + err.Error(),
					Data:    nil,
				}
				return err
			}

			err = begin.Commit()
			if err != nil {
				result = &response.Response[any]{
					Code:    http.StatusInternalServerError,
					Message: "TransactionMiddleware GetMiddleware commit error.",
					Data:    nil,
				}
				return err
			}

			retryCount += 1

			if retryCount > maxRetries {
				result = &response.Response[any]{
					Code:    http.StatusInternalServerError,
					Message: "TransactionMiddleware GetMiddleware max retries reached.",
					Data:    nil,
				}
				err = fmt.Errorf("transactionMiddleware GetMiddleware max retries reached")
				return err
			}

			result = &response.Response[any]{
				Code:    http.StatusOK,
				Message: "TransactionMiddleware GetMiddleware succeed.",
				Data:    nil,
			}
			return err
		})

		if executeErr != nil {
			result = &response.Response[any]{
				Code:    http.StatusInternalServerError,
				Message: "TransactionMiddleware GetMiddleware execute error.",
				Data:    nil,
			}
			response.NewResponse(writer, result)
		}
	})
}
