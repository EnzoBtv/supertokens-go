package testing

import (
	"testing"

	"github.com/supertokens/supertokens-go/supertokens"
	"github.com/supertokens/supertokens-go/supertokens/core"
	"github.com/supertokens/supertokens-go/supertokens/errors"
)

func TestTokenTheftDetection(t *testing.T) {
	beforeEach()
	startST("localhost", "8080")
	supertokens.Config("localhost:8080")

	response, err := core.CreateNewSession("", map[string]interface{}{}, map[string]interface{}{})
	if err != nil {
		t.Error(err)
	}

	response2, err := core.RefreshSession(response.RefreshToken.Token)
	if err != nil {
		t.Error(err)
	}
	_, err = core.GetSession(response2.AccessToken.Token, response2.AntiCsrfToken, true)
	if err != nil {
		t.Error(err)
	}
	_, err = core.RefreshSession(response.RefreshToken.Token)
	if err == nil {
		t.Error("should not have come here")
	} else if !errors.IsTokenTheftDetectedError(err) {
		t.Error("failed")
	}
}

func TestBasicSessionUse(t *testing.T) {
	beforeEach()
	startST("localhost", "8080")
	supertokens.Config("localhost:8080")

	response, err := core.CreateNewSession("", map[string]interface{}{}, map[string]interface{}{})
	if err != nil {
		t.Error(err)
	}

	if response.AccessToken.Token == "" {
		t.Error("accessToken is empty")
	}
	if response.RefreshToken.Token == "" {
		t.Error("refreshToken is empty")
	}
	if response.IDRefreshToken.Token == "" {
		t.Error("idrefreshToken is empty")
	}
	if response.Handle == "" {
		t.Error("handle is empty")
	}
	if response.AntiCsrfToken == nil {
		t.Error("antiCsrfToken is nil")
	}

	_, err = core.GetSession(response.AccessToken.Token, response.AntiCsrfToken, true)
	if err != nil {
		t.Error(err)
	}

	p := core.GetProcessStateInstance()

	if p.GetLastEventByName(core.CallingServiceInVerify) != nil {
		t.Error("processState contains CallingServiceInVerify")
	}

	response2, err := core.RefreshSession(response.RefreshToken.Token)
	if err != nil {
		t.Error(err)
	}

	if response2.AccessToken.Token == "" {
		t.Error("accessToken is empty")
	}
	if response2.RefreshToken.Token == "" {
		t.Error("refreshToken is empty")
	}
	if response2.IDRefreshToken.Token == "" {
		t.Error("idrefreshToken is empty")
	}
	if response2.Handle == "" {
		t.Error("handle is empty")
	}
	if response2.AntiCsrfToken == nil {
		t.Error("antiCsrfToken is nil")
	}

	response3, err := core.GetSession(response2.AccessToken.Token, response2.AntiCsrfToken, true)
	if err != nil {
		t.Error(err)
	}

	if p.GetLastEventByName(core.CallingServiceInVerify) == nil {
		t.Error("processState does not contain CallingServiceInVerify")
	}
	if response3.Handle == "" {
		t.Error("handle is empty")
	}
	if response3.AccessToken.Token == "" {
		t.Error("accessToken is empty")
	}
	if response3.AntiCsrfToken != nil {
		t.Error("antiCsrfToken is not nil")
	}
	if response3.RefreshToken != nil {
		t.Error("refreshToken is not empty")
	}
	if response3.IDRefreshToken != nil {
		t.Error("idrefreshToken is not empty")
	}
	core.ResetProcessState()

	response4, err := core.GetSession(response3.AccessToken.Token, response2.AntiCsrfToken, true)
	if err != nil {
		t.Error(err)
	}
	if core.GetProcessStateInstance().GetLastEventByName(core.CallingServiceInVerify) != nil {
		t.Error("processState contains CallingServiceInVerify")
	}
	if response4.Handle == "" {
		t.Error("handle is empty")
	}
	if response4.AccessToken != nil {
		t.Error("accessToken is not empty")
	}
	if response4.AntiCsrfToken != nil {
		t.Error("antiCsrfToken is not nil")
	}
	if response4.RefreshToken != nil {
		t.Error("refreshToken is not empty")
	}
	if response4.IDRefreshToken != nil {
		t.Error("idrefreshToken is not empty")
	}

	revokeResponse, err := core.RevokeSession(response4.Handle)
	if err != nil {
		t.Error(err)
	}
	if !revokeResponse {
		t.Error("session was not revoked")
	}
}

func TestSessionVerifyWithAntiCSRF(t *testing.T) {
	beforeEach()
	startST("localhost", "8080")
	supertokens.Config("localhost:8080")

	response, err := core.CreateNewSession("", map[string]interface{}{}, map[string]interface{}{})
	if err != nil {
		t.Error(err)
	}

	_, err = core.GetSession(response.AccessToken.Token, response.AntiCsrfToken, true)
	if err != nil {
		t.Error(err)
	}

	_, err = core.GetSession(response.AccessToken.Token, response.AntiCsrfToken, false)
	if err != nil {
		t.Error(err)
	}
}

func TestSessionVerifyWithoutAntiCSRF(t *testing.T) {
	beforeEach()
	startST("localhost", "8080")
	supertokens.Config("localhost:8080")

	response, err := core.CreateNewSession("", map[string]interface{}{}, map[string]interface{}{})
	if err != nil {
		t.Error(err)
	}
	_, err = core.GetSession(response.AccessToken.Token, response.AntiCsrfToken, false)
	if err != nil {
		t.Error(err)
	}
	_, err = core.GetSession(response.AccessToken.Token, nil, true)
	if err == nil {
		t.Error("should not come here")
	} else if !errors.IsTryRefreshTokenError(err) {
		t.Error(err)
	}
}

func TestRevokingOfSessions(t *testing.T) {
	beforeEach()
	startST("localhost", "8080")
	supertokens.Config("localhost:8080")

	_, err := core.RevokeAllSessionsForUser("someUniqueID")
	if err != nil {
		t.Error(err)
	}

	response, err := core.CreateNewSession("someUniqueID", map[string]interface{}{}, map[string]interface{}{})
	if err != nil {
		t.Error(err)
	}

	revokeResponse, err := core.RevokeSession(response.Handle)
	if err != nil {
		t.Error(err)
	}
	if !revokeResponse {
		t.Error("could not revoke session")
	}

	re3, err := core.GetAllSessionHandlesForUser("someUniqueID")
	if err != nil {
		t.Error(err)
	}
	if len(re3) != 0 {
		t.Error("session handles were revoked")
	}

	_, err = core.CreateNewSession("someUniqueID", map[string]interface{}{}, map[string]interface{}{})
	if err != nil {
		t.Error(err)
	}

	_, err = core.CreateNewSession("someUniqueID", map[string]interface{}{}, map[string]interface{}{})
	if err != nil {
		t.Error(err)
	}

	revokedSessions, err := core.RevokeAllSessionsForUser("someUniqueID")
	if err != nil {
		t.Error(err)
	}
	if len(revokedSessions) != 2 {
		t.Error("incorrect number of sessions revoked")
	}

	re3, err = core.GetAllSessionHandlesForUser("someUniqueID")
	if err != nil {
		t.Error(err)
	}
	if len(re3) != 0 {
		t.Error("session handles were revoked")
	}

	revokeResponse, err = core.RevokeSession("")
	if err != nil {
		t.Error(err)
	}
	if revokeResponse {
		t.Error("revoke session which should not exist")
	}

	revokedSessions, err = core.RevokeAllSessionsForUser("random")
	if err != nil {
		t.Error(err)
	}
	if len(revokedSessions) != 0 {
		t.Error("session revoked when it should not have")
	}
}

func TestNoAntiCSRFRequiredIfDisabledFromCore(t *testing.T) {
	beforeEach()
	setKeyValueInConfig("enable_anti_csrf", "false")
	startST("localhost", "8080")
	supertokens.Config("localhost:8080")

	response, err := core.CreateNewSession("", map[string]interface{}{}, map[string]interface{}{})
	if err != nil {
		t.Error(err)
	}
	_, err = core.GetSession(response.AccessToken.Token, response.AntiCsrfToken, true)
	if err != nil {
		t.Error(err)
	}

	_, err = core.GetSession(response.AccessToken.Token, response.AntiCsrfToken, false)
	if err != nil {
		t.Error(err)
	}
}
