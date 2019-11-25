package userman

import (
	"context"
	"errors"
	"log"
	"net/http"

	"golang.org/x/crypto/bcrypt"

	uuid "github.com/satori/go.uuid"

	"github.com/spaceuptech/space-cloud/model"
	"github.com/spaceuptech/space-cloud/utils"
)

// Profile fetches the profile of the user
func (m *Module) Profile(ctx context.Context, token, dbType, project, id string) (int, map[string]interface{}, error) {
	if !m.IsEnabled() {
		return http.StatusNotFound, nil, errors.New("This feature isn't enabled")
	}

	// Create the find object
	find := map[string]interface{}{}
	switch utils.DBType(dbType) {
	case utils.Mongo:
		find["_id"] = id
	default:
		find["id"] = id
	}

	// Create a read request
	req := &model.ReadRequest{Find: find, Operation: utils.One}

	// Check if the user is authenticated
	status, err := m.auth.IsReadOpAuthorised(ctx, project, dbType, "users", token, req)
	if err != nil {
		return status, nil, err
	}

	// Perform database read operation
	res, err := m.crud.Read(ctx, dbType, project, "users", req)
	if err != nil {
		return http.StatusInternalServerError, nil, err
	}

	// Delete password from user object
	delete(res.(map[string]interface{}), "pass")

	return http.StatusOK, res.(map[string]interface{}), nil
}

// Profiles fetches all the user profiles
func (m *Module) Profiles(ctx context.Context, token, dbType, project string) (int, map[string]interface{}, error) {
	if !m.IsEnabled() {
		return http.StatusNotFound, nil, errors.New("This feature isn't enabled")
	}
	// Create the find object
	find := map[string]interface{}{}

	// Creata a read request
	req := &model.ReadRequest{Find: find, Operation: utils.All}

	// Check if the user is authenticated
	status, err := m.auth.IsReadOpAuthorised(ctx, project, dbType, "users", token, req)
	if err != nil {
		return status, nil, err
	}

	res, err := m.crud.Read(ctx, dbType, project, "users", req)
	if err != nil {
		return http.StatusInternalServerError, nil, err
	}

	// Delete password from user object
	if usersArray, ok := res.([]interface{}); ok {
		for _, user := range usersArray {
			userObj := user.(map[string]interface{})
			delete(userObj, "pass")
		}
	}

	return http.StatusOK, map[string]interface{}{"users": res}, nil
}

// EmailSignIn signins the user and returns a JWT token
func (m *Module) EmailSignIn(ctx context.Context, dbType, project, email, password string) (int, map[string]interface{}, error) {
	// Allow this feature only if the email sign in function is enabled
	if !m.IsActive("email") {
		return http.StatusNotFound, nil, errors.New("Email sign in feature is not enabled")
	}

	// Create read request
	readReq := &model.ReadRequest{Find: map[string]interface{}{"email": email}, Operation: utils.One}

	user, err := m.crud.Read(ctx, dbType, project, "users", readReq)
	if err != nil {
		return http.StatusNotFound, nil, errors.New("User not found")
	}

	userObj := user.(map[string]interface{})

	//Compares if the given password is correct
	err = bcrypt.CompareHashAndPassword([]byte(userObj["pass"].(string)), []byte(password))
	if err != nil {
		return http.StatusUnauthorized, nil, errors.New("Given credentials are not correct")
	}

	// Delete password from user
	delete(userObj, "pass")

	req := map[string]interface{}{}
	req["email"] = email

	// Create a token
	if dbType == string(utils.Mongo) {
		req["id"] = userObj["_id"]
	} else {
		req["id"] = userObj["id"]
	}
	req["role"] = userObj["role"]

	token, err := m.auth.CreateToken(req)
	if err != nil {
		return http.StatusInternalServerError, nil, errors.New("Failed to create a JWT token")
	}
	return http.StatusOK, map[string]interface{}{"user": user, "token": token}, nil
}

// EmailSignUp signs up a user and return a JWT token
func (m *Module) EmailSignUp(ctx context.Context, dbType, project, email, name, password, role string) (int, map[string]interface{}, error) {
	// Allow this feature only if the email sign in function is enabled
	if !m.IsActive("email") {
		return http.StatusNotFound, nil, errors.New("Email sign in feature is not enabled")
	}

	//Hash the password that's in the request
	var err error
	password, err = hashPassword(password)
	if err != nil {
		log.Println("Err: ", err)
		return http.StatusInternalServerError, nil, errors.New("Failed to hash password")
	}

	// Create read request
	readReq := &model.ReadRequest{Find: map[string]interface{}{"email": email}, Operation: utils.One}
	_, err = m.crud.Read(ctx, dbType, project, "users", readReq)
	if err == nil {
		return http.StatusConflict, nil, errors.New("User with provided email already exists")
	}

	req := map[string]interface{}{}
	req["email"] = email
	req["pass"] = password
	req["name"] = name
	req["role"] = role
	// Create a create request
	id := uuid.NewV1()
	if dbType == string(utils.Mongo) {
		req["_id"] = id.String()
	} else {
		req["id"] = id.String()
	}
	createReq := &model.CreateRequest{Operation: utils.One, Document: req}
	err = m.crud.Create(ctx, dbType, project, "users", createReq)
	if err != nil {
		log.Println("Err: ", err)
		return http.StatusInternalServerError, nil, errors.New("Failed to create user account")
	}

	delete(req, "pass")

	// Create a new token Object
	tokenObj := map[string]interface{}{
		"email": email,
		"role":  role,
		"id":    id.String()}

	token, err := m.auth.CreateToken(tokenObj)
	if err != nil {
		return http.StatusInternalServerError, nil, errors.New("Failed to create a JWT token")
	}
	return http.StatusOK, map[string]interface{}{"user": req, "token": token}, nil
}

// EmailEditProfile allows the user to edit a profile
func (m *Module) EmailEditProfile(ctx context.Context, token, dbType, project, id, email, name, password string) (int, map[string]interface{}, error) {
	// Allow this feature only if the email sign in function is enabled
	if !m.IsEnabled() {
		return http.StatusNotFound, nil, errors.New("Email sign in feature is not enabled")
	}

	req := &model.UpdateRequest{}
	find := map[string]interface{}{}
	var idString string
	if dbType == string(utils.Mongo) {
		idString = "_id"
	} else {
		idString = "id"
	}
	find[idString] = id
	req.Find = find

	update := map[string]interface{}{}
	set := map[string]interface{}{}
	if email != "" {
		set["email"] = email
	}
	if name != "" {
		set["name"] = name
	}
	if password != "" {
		var err1 error
		password, err1 = hashPassword(password)
		if err1 != nil {
			log.Println("Err: ", err1)
			return http.StatusInternalServerError, nil, errors.New("Failed to hash password")
		}
		set["pass"] = password
	}
	update["$set"] = set
	req.Update = update
	req.Operation = utils.One

	status, err := m.auth.IsUpdateOpAuthorised(ctx, project, dbType, "users", token, req)
	if err != nil {
		return status, nil, err
	}

	err = m.crud.Update(ctx, dbType, project, "users", req)
	if err != nil {
		return http.StatusInternalServerError, nil, err
	}

	readReq := &model.ReadRequest{Find: map[string]interface{}{idString: id}, Operation: utils.One}
	user, err1 := m.crud.Read(ctx, dbType, project, "users", readReq)
	if err1 != nil {
		return http.StatusNotFound, nil, errors.New("User not found")
	}

	userObj := user.(map[string]interface{})

	// Delete password from user
	delete(userObj, "pass")

	req1 := map[string]interface{}{}
	req1["email"] = userObj["email"]
	req1["id"] = userObj[idString]
	req1["role"] = userObj["role"]

	token1, err := m.auth.CreateToken(req1)
	if err != nil {
		return http.StatusInternalServerError, nil, errors.New("Failed to create a JWT token")
	}
	return http.StatusOK, map[string]interface{}{"user": user, "token": token1}, nil
}

func hashPassword(pwd string) (string, error) {
	//Generates a new hash from the given password
	hash, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.MinCost)
	if err != nil {
		return "", err
	}

	//Checks if the hash is correct for the given password
	err = bcrypt.CompareHashAndPassword(hash, []byte(pwd))
	if err != nil {
		return "", err
	}
	return string(hash), nil
}
