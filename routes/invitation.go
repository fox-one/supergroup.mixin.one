package routes

import (
	"net/http"

	"github.com/MixinNetwork/supergroup.mixin.one/middlewares"
	"github.com/MixinNetwork/supergroup.mixin.one/views"
	"github.com/dimfeld/httptreemux"
)

type invitationsImpl struct{}

func registerInvitations(router *httptreemux.TreeMux) {
	impl := &invitationsImpl{}
	router.GET("/invitation_codes", impl.index)
	router.POST("/invitation_codes", impl.create)
	router.PUT("/invitation_codes/:code", impl.apply)
}

func (impl *invitationsImpl) index(w http.ResponseWriter, r *http.Request, params map[string]string) {
	user := middlewares.CurrentUser(r)
	if invitations, err := user.Invitations(r.Context()); err != nil {
		views.RenderErrorResponse(w, r, err)
	} else {
		views.RenderInvitations(w, r, invitations)
	}
}

func (impl *invitationsImpl) create(w http.ResponseWriter, r *http.Request, params map[string]string) {
	user := middlewares.CurrentUser(r)
	if invitations, err := user.CreateInvitations(r.Context()); err != nil {
		views.RenderErrorResponse(w, r, err)
	} else {
		views.RenderInvitations(w, r, invitations)
	}
}

func (impl *invitationsImpl) apply(w http.ResponseWriter, r *http.Request, params map[string]string) {
	// update invitation code status if code exists and unused
	user := middlewares.CurrentUser(r)
	if invitation, err := user.ApplyInvitation(r.Context(), params["code"]); err != nil {
		views.RenderErrorResponse(w, r, err)
	} else {
		views.RenderInvitation(w, r, invitation)
	}
}
