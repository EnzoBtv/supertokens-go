package supertokens

import (
	"context"
	"net/http"

	"github.com/supertokens/supertokens-go/supertokens/core"
	"github.com/supertokens/supertokens-go/supertokens/errors"
)

// Middleware for verifying and refreshing session. ExtraParams are: bool, func(error, http.ResponseWriter)
func Middleware(theirHandler http.HandlerFunc, extraParams ...interface{}) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "OPTIONS" || r.Method == "TRACE" {
			theirHandler.ServeHTTP(w, r)
			return
		}
		var path = r.URL.Path
		handshakeInfo, handshakeInfoError := core.GetHandshakeInfoInstance()
		if handshakeInfoError != nil {
			if len(extraParams) != 2 {
				HandleErrorAndRespond(handshakeInfoError, w)
			} else {
				extraParams[1].(func(err error, w http.ResponseWriter))(handshakeInfoError, w)
			}
			return
		}
		if (handshakeInfo.RefreshTokenPath == path ||
			(handshakeInfo.RefreshTokenPath+"/") == path ||
			handshakeInfo.RefreshTokenPath == (path+"/")) &&
			r.Method == "POST" {
			session, sessionError := RefreshSession(w, r)
			if sessionError != nil {
				if len(extraParams) != 2 {
					HandleErrorAndRespond(sessionError, w)
				} else {
					extraParams[1].(func(err error, w http.ResponseWriter))(sessionError, w)
				}
				return
			}
			ctx := context.WithValue(r.Context(), SessionContext, session)
			theirHandler.ServeHTTP(w, r.WithContext(ctx))
		} else {
			var actualDoAntiCsrfCheck = r.Method != "GET"
			if len(extraParams) != 0 && extraParams[0] != nil {
				actualDoAntiCsrfCheck = extraParams[0].(bool)
			}
			session, sessionError := GetSession(w, r, actualDoAntiCsrfCheck)
			if sessionError != nil {
				if len(extraParams) != 2 {
					HandleErrorAndRespond(sessionError, w)
				} else {
					extraParams[1].(func(err error, w http.ResponseWriter))(sessionError, w)
				}
				return
			}
			ctx := context.WithValue(r.Context(), SessionContext, session)
			theirHandler.ServeHTTP(w, r.WithContext(ctx))
		}
	})
}

// HandleErrorAndRespond if error handlers are provided, then uses those, else does default error handling depending on the type of error
func HandleErrorAndRespond(err error, w http.ResponseWriter) {
	errorHandlers := core.GetErrorHandlersInstance()
	if errors.IsUnauthorizedError(err) {
		errorHandlers.OnUnauthorizedErrorHandler(err, w)
	} else if errors.IsTryRefreshTokenError(err) {
		errorHandlers.OnTryRefreshTokenErrorHandler(err, w)
	} else if errors.IsTokenTheftDetectedError(err) {
		actualError := err.(errors.TokenTheftDetectedError)
		errorHandlers.OnTokenTheftDetectedErrorHandler(actualError.SessionHandle, actualError.UserID, w)
	} else {
		errorHandlers.OnGeneralErrorHandler(err, w)
	}
}
