package auth

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/spaceuptech/helpers"

	"github.com/spaceuptech/space-cloud/gateway/config"
	"github.com/spaceuptech/space-cloud/gateway/model"
	"github.com/spaceuptech/space-cloud/gateway/utils"
)

func (m *Module) matchRule(ctx context.Context, project string, rule *config.Rule, args, auth map[string]interface{}) (*model.PostProcess, error) {
	if project != m.project {
		return nil, formatError(ctx, rule, errors.New("invalid project details provided"))
	}

	if rule.Rule == "allow" || rule.Rule == "authenticated" {
		return nil, nil
	}

	if idTemp, p := auth["id"]; p {
		if id, ok := idTemp.(string); ok && id == utils.InternalUserID {
			return nil, nil
		}
	}

	switch rule.Rule {
	case "deny":
		return nil, formatError(ctx, rule, errors.New("the operation being performed is denied"))

	case "match":
		return nil, match(ctx, rule, args)

	case "and":
		return m.matchAnd(ctx, project, rule, args, auth)

	case "or":
		return m.matchOr(ctx, project, rule, args, auth)

	case "webhook":
		return nil, m.matchFunc(ctx, rule, m.makeHTTPRequest, args)

	case "query":
		return m.matchQuery(ctx, project, rule, m.crud, args, auth)

	case "force":
		return m.matchForce(ctx, project, rule, args, auth)

	case "remove":
		return m.matchRemove(ctx, project, rule, args, auth)

	case "encrypt":
		return m.matchEncrypt(ctx, project, rule, args, auth)

	case "decrypt":
		return m.matchDecrypt(ctx, project, rule, args, auth)

	case "hash":
		return m.matchHash(ctx, project, rule, args, auth)

	default:
		return nil, formatError(ctx, rule, fmt.Errorf("invalid rule type (%s) provided", rule.Rule))
	}
}

func (m *Module) matchFunc(ctx context.Context, rule *config.Rule, MakeHTTPRequest utils.TypeMakeHTTPRequest, args map[string]interface{}) error {
	obj := args["args"].(map[string]interface{})
	token := obj["token"].(string)

	scToken, err := m.GetSCAccessToken(ctx)
	if err != nil {
		return formatError(ctx, rule, err)
	}

	var result interface{}
	if err := MakeHTTPRequest(ctx, http.MethodPost, rule.URL, token, scToken, obj, &result); err != nil {
		return formatError(ctx, rule, err)
	}

	if rule.Store == "" {
		rule.Store = "args.result"
	}

	return formatError(ctx, rule, utils.StoreValue(ctx, rule.Store, result, args))
}

func (m *Module) matchQuery(ctx context.Context, project string, rule *config.Rule, crud model.CrudAuthInterface, args, auth map[string]interface{}) (*model.PostProcess, error) {
	// Adjust the find object to load any variables referenced from state
	find := utils.Adjust(ctx, rule.Find, args).(map[string]interface{})

	// Create a new read request
	req := &model.ReadRequest{Find: find, Operation: utils.All}

	// Execute the read request
	attr := map[string]string{"project": project, "db": rule.DB, "col": rule.Col}
	data, err := crud.Read(ctx, rule.DB, rule.Col, req, model.RequestParams{Claims: auth, Resource: "db-read", Op: "access", Attributes: attr})
	if err != nil {
		return nil, formatError(ctx, rule, err)
	}

	if rule.Store == "" {
		rule.Store = "args.result"
	}
	if err := utils.StoreValue(ctx, rule.Store, data, args); err != nil {
		return nil, formatError(ctx, rule, err)
	}

	postProcess, err := m.matchRule(ctx, project, rule.Clause, args, nil)
	return postProcess, formatError(ctx, rule, err)
}

func (m *Module) matchAnd(ctx context.Context, projectID string, rule *config.Rule, args, auth map[string]interface{}) (*model.PostProcess, error) {
	completeAction := &model.PostProcess{}
	for _, r := range rule.Clauses {
		postProcess, err := m.matchRule(ctx, projectID, r, args, auth)
		// if err is not nil then return error without checking the other clauses.
		if err != nil {
			return &model.PostProcess{}, formatError(ctx, rule, err)
		}
		if postProcess != nil {
			completeAction.PostProcessAction = append(completeAction.PostProcessAction, postProcess.PostProcessAction...)
		}
	}
	return completeAction, nil
}

func (m *Module) matchOr(ctx context.Context, projectID string, rule *config.Rule, args, auth map[string]interface{}) (*model.PostProcess, error) {
	// append all parameters returned by all clauses! and then return mainStruct
	var finalError error
	for _, r := range rule.Clauses {
		postProcess, err := m.matchRule(ctx, projectID, r, args, auth)
		if err == nil {
			// if condition is satisfied -> exit the function
			return postProcess, nil
		}
		finalError = err
	}
	// if condition is not satisfied -> return empty model.PostProcess and error
	return nil, formatError(ctx, rule, finalError)
}

func match(ctx context.Context, rule *config.Rule, args map[string]interface{}) error {
	switch rule.Type {
	case "string":
		return formatError(ctx, rule, matchString(ctx, rule, args))

	case "number":
		return formatError(ctx, rule, matchNumber(ctx, rule, args))

	case "bool":
		return formatError(ctx, rule, matchBool(ctx, rule, args))

	case "date":
		return formatError(ctx, rule, matchDate(ctx, rule, args))
	}

	return formatError(ctx, rule, fmt.Errorf("invalid variable data type (%s) provided", rule.Type))
}

func (m *Module) matchForce(ctx context.Context, projectID string, rule *config.Rule, args, auth map[string]interface{}) (*model.PostProcess, error) {
	if rule.Clause != nil && rule.Clause.Rule != "" {
		// Match clause with rule!
		_, err := m.matchRule(ctx, projectID, rule.Clause, args, auth)
		if err != nil {
			return nil, nil
		}
	}
	value := rule.Value
	if stringValue, ok := rule.Value.(string); ok {
		loadedValue, err := utils.LoadValue(stringValue, args)
		if err == nil {
			value = loadedValue
		}
	}
	// "res" - add to structure for post processing || "args" - store in args
	if strings.HasPrefix(rule.Field, "res") {
		addToStruct := model.PostProcessAction{Action: "force", Field: rule.Field, Value: value}
		return &model.PostProcess{PostProcessAction: []model.PostProcessAction{addToStruct}}, nil
	} else if strings.HasPrefix(rule.Field, "args") {
		err := utils.StoreValue(ctx, rule.Field, value, args)
		return nil, formatError(ctx, rule, err)
	} else {
		return nil, formatError(ctx, rule, ErrIncorrectRuleFieldType)
	}
}

func (m *Module) matchRemove(ctx context.Context, projectID string, rule *config.Rule, args, auth map[string]interface{}) (*model.PostProcess, error) {
	if rule.Clause != nil && rule.Clause.Rule != "" {
		// Match clause with rule!
		_, err := m.matchRule(ctx, projectID, rule.Clause, args, auth)
		if err != nil {
			return nil, nil
		}
	}
	actions := &model.PostProcess{}
	fields, err := m.getFields(ctx, rule.Fields, args)
	if err != nil {
		return nil, err
	}
	for _, value := range fields {
		field, ok := value.(string)
		if !ok {
			return nil, helpers.Logger.LogError(helpers.GetRequestID(ctx), "Invalid value provided for field (Fields) in security rules where rule is (Remove) array contains a value which is not string", err, map[string]interface{}{})
		}
		// "res" - add field to structure for post processing || "args" - delete field from args
		if strings.HasPrefix(field, "res") {
			addToStruct := model.PostProcessAction{Action: "remove", Field: field, Value: nil}
			actions.PostProcessAction = append(actions.PostProcessAction, addToStruct)
		} else if strings.HasPrefix(field, "args") {
			// Since it depends on the request itself, delete the field from args
			if err := utils.DeleteValue(ctx, field, args); err != nil {
				return nil, formatError(ctx, rule, err)
			}
		} else {
			return nil, formatError(ctx, rule, ErrIncorrectRuleFieldType)
		}
	}
	return actions, nil
}

func (m *Module) matchEncrypt(ctx context.Context, projectID string, rule *config.Rule, args, auth map[string]interface{}) (*model.PostProcess, error) {
	actions := &model.PostProcess{}
	if rule.Clause != nil && rule.Clause.Rule != "" {
		// Match clause with rule!
		_, err := m.matchRule(ctx, projectID, rule.Clause, args, auth)
		if err != nil {
			return actions, nil
		}
	}

	fields, err := m.getFields(ctx, rule.Fields, args)
	if err != nil {
		return nil, err
	}
	for _, value := range fields {
		field, ok := value.(string)
		if !ok {
			return nil, helpers.Logger.LogError(helpers.GetRequestID(ctx), "Invalid value provided for field (Fields) in security rules where rule is (Encrypt) array contains a value which is not string", err, map[string]interface{}{})
		}
		if strings.HasPrefix(field, "res") {
			addToStruct := model.PostProcessAction{Action: "encrypt", Field: field}
			actions.PostProcessAction = append(actions.PostProcessAction, addToStruct)
		} else if strings.HasPrefix(field, "args") {
			loadedValue, err := utils.LoadValue(field, args)
			if err != nil {
				return nil, formatError(ctx, rule, helpers.Logger.LogError(helpers.GetRequestID(ctx), "error loading value in matchEncrypt", err, nil))
			}
			stringValue, ok := loadedValue.(string)
			if !ok {
				return nil, formatError(ctx, rule, fmt.Errorf("Value should be of type string and not %T", loadedValue))
			}
			encryptedValue, err := utils.Encrypt(m.aesKey, stringValue)
			if err != nil {
				return nil, helpers.Logger.LogError(helpers.GetRequestID(ctx), "Unable to encrypt string", err, map[string]interface{}{"valueToEncrypt": stringValue})
			}

			if err = utils.StoreValue(ctx, field, encryptedValue, args); err != nil {
				return nil, formatError(ctx, rule, helpers.Logger.LogError(helpers.GetRequestID(ctx), "error storing value in matchEncrypt", err, nil))
			}
		} else {
			return nil, formatError(ctx, rule, helpers.Logger.LogError(helpers.GetRequestID(ctx), "Invalid field provided for rule encrypt it should either start from res. or args.", fmt.Errorf("invalid field (%s) provided", field), nil))
		}
	}
	return actions, nil
}

func (m *Module) matchDecrypt(ctx context.Context, projectID string, rule *config.Rule, args, auth map[string]interface{}) (*model.PostProcess, error) {
	actions := &model.PostProcess{}
	if rule.Clause != nil && rule.Clause.Rule != "" {
		// Match clause with rule!
		_, err := m.matchRule(ctx, projectID, rule.Clause, args, auth)
		if err != nil {
			return actions, nil
		}
	}

	fields, err := m.getFields(ctx, rule.Fields, args)
	if err != nil {
		return nil, err
	}
	for _, value := range fields {
		field, ok := value.(string)
		if !ok {
			return nil, helpers.Logger.LogError(helpers.GetRequestID(ctx), "Invalid value provided for field (Fields) in security rules where rule is (Decrypt) array contains a value which is not string", err, map[string]interface{}{})
		}
		if strings.HasPrefix(field, "res") {
			addToStruct := model.PostProcessAction{Action: "decrypt", Field: field}
			actions.PostProcessAction = append(actions.PostProcessAction, addToStruct)
		} else if strings.HasPrefix(field, "args") {
			loadedValue, err := utils.LoadValue(field, args)
			if err != nil {
				return nil, formatError(ctx, rule, helpers.Logger.LogError(helpers.GetRequestID(ctx), "error loading value in matchDecrypt", err, nil))
			}
			stringValue, ok := loadedValue.(string)
			if !ok {
				return nil, formatError(ctx, rule, fmt.Errorf("Value should be of type string and not %T", loadedValue))
			}
			decodedValue, err := base64.StdEncoding.DecodeString(stringValue)
			if err != nil {
				return nil, formatError(ctx, rule, err)
			}
			decrypted := make([]byte, len(decodedValue))
			err1 := decryptAESCFB(decrypted, decodedValue, m.aesKey, m.aesKey[:aes.BlockSize])
			if err1 != nil {
				return nil, formatError(ctx, rule, helpers.Logger.LogError(helpers.GetRequestID(ctx), "error decrypting value in matchDecrypt", err, nil))
			}
			er := utils.StoreValue(ctx, field, string(decrypted), args)
			if er != nil {
				return nil, formatError(ctx, rule, helpers.Logger.LogError(helpers.GetRequestID(ctx), "error storing value in matchDecrypt", err, nil))
			}
		} else {
			return nil, formatError(ctx, rule, fmt.Errorf("invalid field (%s) provided", field))
		}
	}
	return actions, nil
}

func decryptAESCFB(dst, src, key, iv []byte) error {
	aesBlockDecrypter, err := aes.NewCipher([]byte(key))
	if err != nil {
		return err
	}
	aesDecrypter := cipher.NewCFBDecrypter(aesBlockDecrypter, iv)
	aesDecrypter.XORKeyStream(dst, src)
	return nil
}

func (m *Module) matchHash(ctx context.Context, projectID string, rule *config.Rule, args, auth map[string]interface{}) (*model.PostProcess, error) {
	actions := &model.PostProcess{}
	if rule.Clause != nil && rule.Clause.Rule != "" {
		// Match clause with rule!
		_, err := m.matchRule(ctx, projectID, rule.Clause, args, auth)
		if err != nil {
			return actions, nil
		}
	}

	fields, err := m.getFields(ctx, rule.Fields, args)
	if err != nil {
		return nil, err
	}
	for _, value := range fields {
		field, ok := value.(string)
		if !ok {
			return nil, helpers.Logger.LogError(helpers.GetRequestID(ctx), "Invalid value provided for field (Fields) in security rules where rule is (Hash) array contains a value which is not string", err, map[string]interface{}{})
		}
		if strings.HasPrefix(field, "res") {
			addToStruct := model.PostProcessAction{Action: "hash", Field: field}
			actions.PostProcessAction = append(actions.PostProcessAction, addToStruct)
		} else if strings.HasPrefix(field, "args") {
			loadedValue, err := utils.LoadValue(field, args)
			if err != nil {
				return nil, formatError(ctx, rule, helpers.Logger.LogError(helpers.GetRequestID(ctx), "error loading value in matchHash", err, nil))
			}
			stringValue, ok := loadedValue.(string)
			if !ok {
				return nil, formatError(ctx, rule, fmt.Errorf("Value should be of type string and not %T", loadedValue))
			}
			hashed := utils.HashString(stringValue)
			er := utils.StoreValue(ctx, field, hashed, args)
			if er != nil {
				return nil, formatError(ctx, rule, helpers.Logger.LogError(helpers.GetRequestID(ctx), "error storing value in matchHash", err, nil))
			}
		} else {
			return nil, formatError(ctx, rule, fmt.Errorf("invalid field (%s) provided", field))
		}
	}
	return actions, nil
}
