package userman

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/dghubble/gologin"
	"github.com/dghubble/gologin/facebook"
	"github.com/dghubble/gologin/github"
	"github.com/dghubble/gologin/google"
	"github.com/dghubble/gologin/twitter"
	"github.com/dghubble/oauth1"
	"github.com/gorilla/mux"
	uuid "github.com/satori/go.uuid"
	"github.com/spaceuptech/space-cloud/model"
	"github.com/spaceuptech/space-cloud/utils"
	"golang.org/x/oauth2"

	twitterOAuth1 "github.com/dghubble/oauth1/twitter"
	facebookOAuth2 "golang.org/x/oauth2/facebook"
	githubOAuth2 "golang.org/x/oauth2/github"
	googleOAuth2 "golang.org/x/oauth2/google"
)

// OAuthConfig configures the main ServeMux.
type OAuthConfig struct {
	ClientID     string
	ClientSecret string
}

// NewGoogleOAuth registers the routes for Google OAuth
func (m *Module) NewGoogleOAuth(router *mux.Router) {

	// Don't register routes if module is disabled
	if enabled := m.IsActive("google"); !enabled {
		return
	}

	// Get the OAuth config
	c, _ := m.getOAuth("google")

	oauth2Config := &oauth2.Config{
		ClientID:     c.ID,
		ClientSecret: c.Secret,
		RedirectURL:  c.Host + "/v1/api/" + m.project + "/auth/" + c.DBType + "/oauth/google-callback",
		Endpoint:     googleOAuth2.Endpoint,
		Scopes:       []string{"profile", "email"},
	}

	// state param cookies require HTTPS by default; disable for localhost development
	//stateConfig := gologin.DebugOnlyCookieConfig
	stateConfig := gologin.DefaultCookieConfig
	router.Handle("/oauth/google", google.StateHandler(stateConfig, google.LoginHandler(oauth2Config, nil)))
	router.Handle("/oauth/google-callback", google.StateHandler(stateConfig, google.CallbackHandler(oauth2Config, issueGoogleSession(m), nil)))
}

func issueGoogleSession(m *Module) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		// Get the path parameters
		vars := mux.Vars(r)
		project := vars["project"]
		dbType := vars["dbType"]

		ctx := r.Context()
		googleUser, err := google.UserFromContext(ctx)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		token, err := m.addUser(ctx, project, dbType, googleUser.Email, googleUser.Name)
		if err != nil {
			log.Println("Err: ", err)
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		// Get the OAuth config
		c, _ := m.getOAuth("google")

		http.Redirect(w, r, c.RedirectURL+"?token="+token, http.StatusSeeOther)
	}
	return http.HandlerFunc(fn)
}

// NewFacebookOAuth registers the routes for Facebook OAuth
func (m *Module) NewFacebookOAuth(router *mux.Router) {
	// Don't register routes if module is disabled
	if enabled := m.IsActive("fb"); !enabled {
		return
	}

	// Get the OAuth config
	c, _ := m.getOAuth("fb")

	oauth2Config := &oauth2.Config{
		ClientID:     c.ID,
		ClientSecret: c.Secret,
		RedirectURL:  c.Host + "/v1/api/" + m.project + "/auth/" + c.DBType + "/oauth/fb-callback",
		Endpoint:     facebookOAuth2.Endpoint,
		Scopes:       []string{"email"},
	}

	// state param cookies require HTTPS by default; disable for localhost development
	//stateConfig := gologin.DebugOnlyCookieConfig
	stateConfig := gologin.DefaultCookieConfig
	router.Handle("/oauth/fb", facebook.StateHandler(stateConfig, facebook.LoginHandler(oauth2Config, nil)))
	router.Handle("/oauth/fb-callback", facebook.StateHandler(stateConfig, facebook.CallbackHandler(oauth2Config, issueFacebookSession(m), nil)))
}

func issueFacebookSession(m *Module) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		// Get the path parameters
		vars := mux.Vars(r)
		project := vars["project"]
		dbType := vars["dbType"]

		ctx := r.Context()
		facebookUser, err := facebook.UserFromContext(ctx)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		token, err := m.addUser(ctx, project, dbType, facebookUser.Email, facebookUser.Name)
		if err != nil {
			log.Println("Err: ", err)
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		// Get the OAuth config
		c, _ := m.getOAuth("fb")

		http.Redirect(w, r, c.RedirectURL+"?token="+token, http.StatusSeeOther)
	}
	return http.HandlerFunc(fn)
}

// NewGithubOAuth registers the routes for Github OAuth
func (m *Module) NewGithubOAuth(router *mux.Router) {
	// Don't register routes if module is disabled
	if enabled := m.IsActive("github"); !enabled {
		return
	}

	// Get the OAuth config
	c, _ := m.getOAuth("github")

	oauth2Config := &oauth2.Config{
		ClientID:     c.ID,
		ClientSecret: c.Secret,
		RedirectURL:  c.Host + "/v1/api/" + m.project + "/auth/" + c.DBType + "/oauth/github-callback",
		Endpoint:     githubOAuth2.Endpoint,
	}
	// state param cookies require HTTPS by default; disable for localhost development
	stateConfig := gologin.DefaultCookieConfig
	router.Handle("/oauth/github", github.StateHandler(stateConfig, github.LoginHandler(oauth2Config, nil)))
	router.Handle("/oauth/github-callback", github.StateHandler(stateConfig, github.CallbackHandler(oauth2Config, issueGithubSession(m), nil)))
}

func issueGithubSession(m *Module) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		// Get the path parameters
		vars := mux.Vars(r)
		project := vars["project"]
		dbType := vars["dbType"]

		ctx := r.Context()
		githubUser, err := github.UserFromContext(ctx)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		token, err := m.addUser(ctx, project, dbType, githubUser.GetEmail(), githubUser.GetName())
		if err != nil {
			log.Println("Err: ", err)
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		// Get the OAuth config
		c, _ := m.getOAuth("github")

		http.Redirect(w, r, c.RedirectURL+"?token="+token, http.StatusSeeOther)
	}
	return http.HandlerFunc(fn)
}

// NewTwitterOAuth registers the routes for Twitter OAuth
func (m *Module) NewTwitterOAuth(router *mux.Router) {
	// Don't register routes if module is disabled
	if enabled := m.IsActive("twitter"); !enabled {
		return
	}

	// Get the OAuth config
	c, _ := m.getOAuth("twitter")

	oauth1Config := &oauth1.Config{
		ConsumerKey:    c.ID,
		ConsumerSecret: c.Secret,
		CallbackURL:    c.Host + "/v1/api/" + m.project + "/auth/" + c.DBType + "/oauth/twitter-callback",
		Endpoint:       twitterOAuth1.AuthorizeEndpoint,
	}
	router.Handle("/oauth/twitter", twitter.LoginHandler(oauth1Config, nil))
	router.Handle("/oauth/twitter-callback", twitter.CallbackHandler(oauth1Config, issueTwitterSession(m), nil))

}

func issueTwitterSession(m *Module) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		// Get the path parameters
		vars := mux.Vars(r)
		project := vars["project"]
		dbType := vars["dbType"]

		ctx := r.Context()
		twitterUser, err := twitter.UserFromContext(ctx)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		token, err := m.addUser(ctx, project, dbType, twitterUser.Email, twitterUser.Name)
		if err != nil {
			log.Println("Err: ", err)
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		// Get the OAuth config
		c, _ := m.getOAuth("twitter")

		http.Redirect(w, r, c.RedirectURL+"?token="+token, http.StatusSeeOther)
	}
	return http.HandlerFunc(fn)
}

func (m *Module) addUser(ctx context.Context, project, dbType, email, name string) (string, error) {
	var userObj map[string]interface{}
	id := uuid.NewV1().String()

	// Create a create request
	readReq := &model.ReadRequest{Find: map[string]interface{}{"email": email}, Operation: utils.One}
	user, err := m.crud.Read(ctx, dbType, project, "users", readReq)
	if err == nil {
		userObj = user.(map[string]interface{})
		if dbType == string(utils.Mongo) {
			id = userObj["_id"].(string)
		} else {
			id = userObj["id"].(string)
		}
	} else {
		// Create new user object
		userObj = map[string]interface{}{"email": email, "name": name, "role": m.getDefaultRole()}
		if dbType == string(utils.Mongo) {
			userObj["_id"] = id
		} else {
			userObj["id"] = id
		}

		// Create the user in the db
		createReq := &model.CreateRequest{Operation: utils.One, Document: user}
		err = m.crud.Create(ctx, dbType, project, "users", createReq)
		if err != nil {
			return "", err
		}
	}

	// Create a new token Object
	tokenObj := map[string]interface{}{
		"email": email,
		"role":  userObj["role"],
		"id":    id,
	}

	token, err := m.auth.CreateToken(tokenObj)
	if err != nil {
		return "", err
	}

	return token, nil
}
