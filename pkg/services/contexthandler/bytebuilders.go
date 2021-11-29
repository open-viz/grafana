package contexthandler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/grafana/grafana/pkg/bus"
	"github.com/grafana/grafana/pkg/models"
	"github.com/grafana/grafana/pkg/setting"
	"github.com/grafana/grafana/pkg/util"
)

// ByteBuildersUser represents logged in b3 user in the system
type ByteBuildersUser struct {
	// the user's id
	ID int64 `json:"id"`
	// the user's username
	UserName string `json:"login"`
	// the user's full name
	FullName string `json:"full_name"`
	// swagger:strfmt email
	Email string `json:"email"`
	// URL to the user's avatar
	AvatarURL string `json:"avatar_url"`
	// User locale
	Language string `json:"language"`
	// Is the user an administrator
	IsAdmin bool `json:"is_admin"`
	// swagger:strfmt date-time
	LastLogin time.Time `json:"last_login,omitempty"`
	// swagger:strfmt date-time
	Created time.Time `json:"created,omitempty"`
	// Is user restricted
	Restricted bool `json:"restricted"`
	// Is user active
	IsActive bool `json:"active"`
	// Is user login prohibited
	ProhibitLogin bool `json:"prohibit_login"`
	// the user's location
	Location string `json:"location"`
	// the user's website
	Website string `json:"website"`
	// the user's description
	Description string `json:"description"`
}

func (h *ContextHandler) initContextWithByteBuilders(ctx *models.ReqContext) bool {
	user, err := loginWithByteBuildersCookie(ctx)
	if err != nil {
		ctx.Logger.Error("failed to get user", "error", err)
		return false
	}

	loginQuery := models.GetUserByLoginQuery{LoginOrEmail: user.UserName}
	if err := bus.Dispatch(&loginQuery); err != nil {
		ctx.JsonApiErr(401, "Basic auth failed", err)
		return true
	}

	gUser := loginQuery.Result

	query := models.GetSignedInUserQuery{UserId: gUser.Id, OrgId: gUser.OrgId}
	if err := bus.Dispatch(&query); err != nil {
		ctx.JsonApiErr(401, "Authentication error", err)
		return true
	}

	ctx.Resp.Header().Set("X-Grafana-Org-Id", strconv.FormatInt(gUser.OrgId, 10))

	ctx.SignedInUser = query.Result
	ctx.IsSignedIn = true
	return true

}

func loginWithByteBuildersCookie(ctx *models.ReqContext) (ByteBuildersUser, error) {
	baseUrl := util.GetBaseUrl(ctx.Req.Host, setting.Env == setting.Prod)

	req, err := http.NewRequest("GET", fmt.Sprintf("%v/api/v1/user", baseUrl), nil)
	if err != nil {
		return ByteBuildersUser{}, err
	}

	for _, cookie := range ctx.Req.Cookies() {
		req.AddCookie(cookie)
	}

	req.Header = ctx.Req.Header
	req.Header.Add("X-Csrf-Token", ctx.GetCookie("_csrf"))
	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return ByteBuildersUser{}, err
	}
	defer resp.Body.Close()

	var user ByteBuildersUser
	if err = json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return ByteBuildersUser{}, err
	}

	return user, nil
}
