package routes

import (
	"net/http"

	"github.com/MixinNetwork/supergroup.mixin.one/middlewares"
	"github.com/MixinNetwork/supergroup.mixin.one/views"
	"github.com/MixinNetwork/supergroup.mixin.one/models"
	"github.com/dimfeld/httptreemux"
)

type invitationsImpl struct{}

func registerInvitations(router *httptreemux.TreeMux) {
	impl := &invitationsImpl{}
	router.GET("/invitations", impl.index)
	router.POST("/invitations", impl.create)
	router.PUT("/invitations/:code", impl.apply)
}

func (impl *invitationsImpl) index(w http.ResponseWriter, r *http.Request, params map[string]string) {
	user := middlewares.CurrentUser(r)
	var err error
	var invitations []*models.Invitation
	keys, ok := r.URL.Query()["history"]
	if ok && keys[0] == "true" {
		invitations, err = user.InvitationsHistory(r.Context())
	} else {
		invitations, err = user.Invitations(r.Context())
	}
	if err != nil {
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
	user := middlewares.CurrentUser(r)
	if invitation, err := user.ApplyInvitation(r.Context(), params["code"]); err != nil {
		views.RenderErrorResponse(w, r, err)
	} else {
		views.RenderInvitation(w, r, invitation)
	}
}