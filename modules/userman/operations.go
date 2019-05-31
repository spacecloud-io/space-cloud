package userman

import (
	"context"
	"log"
	"net/http"

	"golang.org/x/crypto/bcrypt"

	uuid "github.com/satori/go.uuid"

	"github.com/spaceuptech/space-cloud/model"
	"github.com/spaceuptech/space-cloud/utils"
)

func (m *Module) Profile(ctx context.Context, token, dbType, project, id string) map[string]interface{} {
	result := map[string]interface{}{"status": nil, "error": nil, "result": nil}
	if !m.IsEnabled() {
		result["status"] = http.StatusNotFound
		result["error"] = "This feature isn't enabled"
		return result
	}
	
	authObj, err := m.auth.IsAuthenticated(token, dbType, "users", utils.Read)
	if err != nil {
		result["status"] = http.StatusUnauthorized
		result["error"] = err.Error()
		return result
	}

	// Create the find object
	find := map[string]interface{}{}

	switch utils.DBType(dbType) {
	case utils.Mongo:
		find["_id"] = id
	default:
		find["id"] = id
	}

	// Create an args object
	args := map[string]interface{}{
		"args":    map[string]interface{}{"find": find, "op": utils.One, "auth": authObj},
		"project": project, // Don't forget to do this for every request
	}

	// Check if user is authorized to make this request
	err = m.auth.IsAuthorized(project, dbType, "users", utils.Read, args)
	if err != nil {
		result["status"] = http.StatusForbidden
		result["error"] = err.Error()
		return result
	}
	
	req := &model.ReadRequest{Find: find, Operation: utils.One}
	res, err := m.crud.Read(ctx, dbType, project, "users", req)
	if err != nil {
		result["status"] = http.StatusInternalServerError
		result["error"] = err.Error()
		return result
	}

	// Delete password from user object
	delete(res.(map[string]interface{}), "pass")
	
	result["status"] = http.StatusOK
	result["result"] = res
	return result
}

func (m *Module) Profiles(ctx context.Context, token, dbType, project string) map[string]interface{} {
	result := map[string]interface{}{"status": nil, "error": nil, "result": nil}
	if !m.IsEnabled() {
		result["status"] = http.StatusNotFound
		result["error"] = "This feature isn't enabled"
		return result
	}
	
	authObj, err := m.auth.IsAuthenticated(token, dbType, "users", utils.Read)
	if err != nil {
		result["status"] = http.StatusUnauthorized
		result["error"] = err.Error()
		return result
	}

	// Create the find object
	find := map[string]interface{}{}

	// Create an args object
	args := map[string]interface{}{
		"args":    map[string]interface{}{"find": find, "op": utils.All, "auth": authObj},
		"project": project, // Don't forget to do this for every request
	}

	// Check if user is authorized to make this request
	err = m.auth.IsAuthorized(project, dbType, "users", utils.Read, args)
	if err != nil {
		result["status"] = http.StatusForbidden
		result["error"] = err.Error()
		return result
	}
	
	req := &model.ReadRequest{Find: find, Operation: utils.All}
	res, err := m.crud.Read(ctx, dbType, project, "users", req)
	if err != nil {
		result["status"] = http.StatusInternalServerError
		result["error"] = err.Error()
		return result
	}

	// Delete password from user object
	if usersArray, ok := res.([]interface{}); ok {
		for _, user := range usersArray {
			userObj := user.(map[string]interface{})
			delete(userObj, "pass")
		}
	}
	
	result["status"] = http.StatusOK
	result["result"] = res
	return result
}

func (m *Module) EmailSignIn(ctx context.Context, dbType, project, email, password string) map[string]interface{} {
	result := map[string]interface{}{"status": nil, "error": nil, "result": nil}
	// Allow this feature only if the email sign in function is enabled
	if !m.IsActive("email") {
		result["status"] = http.StatusNotFound
		result["error"] = "Email sign in feature is not enabled"
		return result
	}

	// Create read request
	readReq := &model.ReadRequest{Find: map[string]interface{}{"email": email}, Operation: utils.One}

	user, err := m.crud.Read(ctx, dbType, project, "users", readReq)
	if err != nil {
		result["status"] = http.StatusNotFound
		result["error"] = "User not found"
		return result
	}

	userObj := user.(map[string]interface{})

	//Compares if the given password is correct
	err = bcrypt.CompareHashAndPassword([]byte(userObj["pass"].(string)), []byte(password))
	if err != nil {
		result["status"] = http.StatusUnauthorized
		result["error"] = "Given credentials are not correct"
		return result
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
		result["status"] = http.StatusInternalServerError
		result["error"] = "Failed to create a JWT token"
		return result
	}
	result["status"] = http.StatusOK
	result["result"] = map[string]interface{}{"user": user, "token": token}
	return result
}

func (m *Module) EmailSignUp(ctx context.Context, dbType, project, email, name, password, role string) map[string]interface{} {
	result := map[string]interface{}{"status": nil, "error": nil, "result": nil}
	// Allow this feature only if the email sign in function is enabled
	if !m.IsActive("email") {
		result["status"] = http.StatusNotFound
		result["error"] = "Email sign in feature is not enabled"
		return result
	}

	//Hash the password that's in the request
	var err error
	password, err = hashPassword(password)
	if err != nil {
		log.Println("Err: ", err)
		result["status"] = http.StatusInternalServerError
		result["error"] = "Failed to hash password"
		return result
	}

	// Create read request
	readReq := &model.ReadRequest{Find: map[string]interface{}{"email": email}, Operation: utils.One}
	_, err = m.crud.Read(ctx, dbType, project, "users", readReq)
	if err == nil {
		log.Println("Err: ", err)
		result["status"] = http.StatusConflict
		result["error"] = "User with provided email already exists"
		return result
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
		result["status"] = http.StatusInternalServerError
		result["error"] = "Failed to create user account"
		return result
	}

	delete(req, "pass")

	// Create a new token Object
	tokenObj := map[string]interface{}{
		"email": email,
		"role":  role,
	}
	tokenObj["id"] = id.String()
	
	token, err := m.auth.CreateToken(tokenObj)
	if err != nil {
		result["status"] = http.StatusInternalServerError
		result["error"] = "Failed to create a JWT token"
		return result
	}
	result["status"] = http.StatusOK
	result["result"] = map[string]interface{}{"user": req, "token": token}
	return result
}

func (m *Module) EmailEditProfile(ctx context.Context, token, dbType, project, id, email, name, password string) map[string]interface{} {
	result := map[string]interface{}{"status": nil, "error": nil, "result": nil}
	// Allow this feature only if the email sign in function is enabled
	if !m.IsActive("email") {
		result["status"] = http.StatusNotFound
		result["error"] = "Email sign in feature is not enabled"
		return result
	}

	authObj, err := m.auth.IsAuthenticated(token, dbType, "users", utils.Update)
	if err != nil {
		result["status"] = http.StatusUnauthorized
		result["error"] = err.Error()
		return result
	}

	req := model.UpdateRequest{}
	temp := map[string]interface{}{}
	var id_string string
	if dbType == string(utils.Mongo) {
		id_string = "_id"
	} else {
		id_string = "id"
	}
	temp[id_string] = id
	req.Find = temp

	temp = map[string]interface{}{}
	temp_set := map[string]interface{}{}
	if email != "" {
		temp_set["email"] = email
	}
	if name != "" {
		temp_set["name"] = name
	}
	if password != "" {
		var err1 error
		password, err1 = hashPassword(password)
		if err1 != nil {
			log.Println("Err: ", err1)
			result["status"] = http.StatusInternalServerError
			result["error"] = "Failed to hash password"
			return result
		}
		temp_set["pass"] = password
	}
	temp["$set"] = temp_set
	req.Update = temp
	req.Operation = utils.One

	// Create an args object
	args := map[string]interface{}{
		"args":    map[string]interface{}{"find": req.Find, "op": req.Operation, "auth": authObj},
		"project": project, // Don't forget to do this for every request
	}

	// Check if user is authorized to make this request
	err = m.auth.IsAuthorized(project, dbType, "users", utils.Update, args)
	if err != nil {
		result["status"] = http.StatusForbidden
		result["error"] = err.Error()
		return result
	}

	err = m.crud.Update(ctx, dbType, project, "users", &req)
	if err != nil {
		result["status"] = http.StatusInternalServerError
		result["error"] = err.Error()
		return result
	}

	readReq := &model.ReadRequest{Find: map[string]interface{}{id_string: id}, Operation: utils.One}
	user, err1 := m.crud.Read(ctx, dbType, project, "users", readReq)
	if err1 != nil {
		result["status"] = http.StatusNotFound
		result["error"] = "User not found"
		return result
	}

	userObj := user.(map[string]interface{})

	// Delete password from user
	delete(userObj, "pass")

	req1 := map[string]interface{}{}
	req1["email"] = userObj["email"]
	req1["id"] = userObj[id_string]
	req1["role"] = userObj["role"]

	token1, err := m.auth.CreateToken(req1)
	if err != nil {
		result["status"] = http.StatusInternalServerError
		result["error"] = "Failed to create a JWT token"
		return result
	}
	result["status"] = http.StatusOK
	result["result"] = map[string]interface{}{"user": user, "token": token1}
	return result
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
