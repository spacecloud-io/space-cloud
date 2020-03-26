package auth

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/spaceuptech/space-cloud/gateway/config"
	"github.com/spaceuptech/space-cloud/gateway/model"
	"github.com/spaceuptech/space-cloud/gateway/utils"
)

func (m *Module) matchRule(ctx context.Context, project string, rule *config.Rule, args, auth map[string]interface{}) (*model.PostProcess, error) {
	if project != m.project {
		return nil, errors.New("invalid project details provided")
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
		return nil, errors.New("the operation being performed is denied")

	case "match":
		return nil, match(rule, args)

	case "and":
		return m.matchAnd(ctx, project, rule, args, auth)

	case "or":
		return m.matchOr(ctx, project, rule, args, auth)

	case "webhook":
		return nil, m.matchFunc(ctx, rule, m.makeHTTPRequest, args)

	case "query":
		return m.matchQuery(ctx, project, rule, m.crud, args)

	case "force":
		return m.matchForce(ctx, project, rule, args, auth)

	case "remove":
		return m.matchRemove(ctx, project, rule, args, auth)

	case "encrypt":
		return m.matchEncrypt(rule, args)

	case "decrypt":
		return m.matchDecrypt(rule, args)

	case "hash":
		return matchHash(rule, args)

	default:
		return nil, ErrIncorrectMatch
	}
}

func (m *Module) matchFunc(ctx context.Context, rule *config.Rule, MakeHTTPRequest utils.MakeHTTPRequest, args map[string]interface{}) error {
	obj := args["args"].(map[string]interface{})
	token := obj["token"].(string)

	scToken, err := m.GetSCAccessToken()
	if err != nil {
		return err
	}

	var result interface{}
	return MakeHTTPRequest(ctx, "POST", rule.URL, token, scToken, obj, &result)
}

func (m *Module) matchQuery(ctx context.Context, project string, rule *config.Rule, crud model.CrudAuthInterface, args map[string]interface{}) (*model.PostProcess, error) {
	// Adjust the find object to load any variables referenced from state
	rule.Find = utils.Adjust(rule.Find, args).(map[string]interface{})

	// Create a new read request
	req := &model.ReadRequest{Find: rule.Find, Operation: utils.All}

	// Execute the read request
	data, err := crud.Read(ctx, rule.DB, project, rule.Col, req)
	if err != nil {
		return nil, err
	}
	args["result"] = data
	return m.matchRule(ctx, project, rule.Clause, args, nil)
}

func (m *Module) matchAnd(ctx context.Context, projectID string, rule *config.Rule, args, auth map[string]interface{}) (*model.PostProcess, error) {
	completeAction := &model.PostProcess{}
	for _, r := range rule.Clauses {
		postProcess, err := m.matchRule(ctx, projectID, r, args, auth)
		// if err is not nil then return error without checking the other clauses.
		if err != nil {
			return &model.PostProcess{}, err
		}
		if postProcess != nil {
			completeAction.PostProcessAction = append(completeAction.PostProcessAction, postProcess.PostProcessAction...)
		}
	}
	return completeAction, nil
}

func (m *Module) matchOr(ctx context.Context, projectID string, rule *config.Rule, args, auth map[string]interface{}) (*model.PostProcess, error) {
	// append all parameters returned by all clauses! and then return mainStruct
	for _, r := range rule.Clauses {
		postProcess, err := m.matchRule(ctx, projectID, r, args, auth)
		if err == nil {
			// if condition is satisfied -> exit the function
			return postProcess, nil
		}
	}
	// if condition is not satisfied -> return empty model.PostProcess and error
	return nil, ErrIncorrectMatch
}

func match(rule *config.Rule, args map[string]interface{}) error {
	switch rule.Type {
	case "string":
		return matchString(rule, args)

	case "number":
		return matchNumber(rule, args)

	case "bool":
		return matchBool(rule, args)
	}

	return ErrIncorrectMatch
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
		err := utils.StoreValue(rule.Field, value, args)
		return nil, err
	} else {
		return nil, ErrIncorrectRuleFieldType
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
	for _, field := range rule.Fields {
		// "res" - add field to structure for post processing || "args" - delete field from args
		if strings.HasPrefix(field, "res") {
			addToStruct := model.PostProcessAction{Action: "remove", Field: field, Value: nil}
			actions.PostProcessAction = append(actions.PostProcessAction, addToStruct)
		} else if strings.HasPrefix(field, "args") {
			// Since it depends on the request itself, delete the field from args
			if err := utils.DeleteValue(field, args); err != nil {
				return nil, err
			}
		} else {
			return nil, ErrIncorrectRuleFieldType
		}
	}
	return actions, nil
}

func (m *Module) matchEncrypt(rule *config.Rule, args map[string]interface{}) (*model.PostProcess, error) {
	actions := &model.PostProcess{}
	for _, field := range rule.Fields {
		if strings.HasPrefix(field, "res") {
			addToStruct := model.PostProcessAction{Action: "encrypt", Field: field}
			actions.PostProcessAction = append(actions.PostProcessAction, addToStruct)
		} else if strings.HasPrefix(field, "args") {
			loadedValue, err := utils.LoadValue(field, args)
			if err != nil {
				logrus.Errorln("error loading value in matchEncrypt: ", err)
				return nil, err
			}
			stringValue, ok := loadedValue.(string)
			if !ok {
				return nil, fmt.Errorf("Value should be of type string and not %T", loadedValue)
			}
			encrypted := make([]byte, len(stringValue))
			err1 := encryptAESCFB(encrypted, []byte(stringValue), m.aesKey, m.aesKey[:aes.BlockSize])
			if err1 != nil {
				logrus.Errorln("error encrypting value in matchEncrypt: ", err1)
				return nil, err1
			}
			er := utils.StoreValue(field, base64.StdEncoding.EncodeToString(encrypted), args)
			if er != nil {
				logrus.Errorln("error storing value in matchEncrypt: ", er)
				return nil, er
			}
		} else {
			return nil, fmt.Errorf("invalid field (%s) provided", field)
		}
	}
	return actions, nil
}

func (m *Module) matchDecrypt(rule *config.Rule, args map[string]interface{}) (*model.PostProcess, error) {
	actions := &model.PostProcess{}
	for _, field := range rule.Fields {
		if strings.HasPrefix(field, "res") {
			addToStruct := model.PostProcessAction{Action: "decrypt", Field: field}
			actions.PostProcessAction = append(actions.PostProcessAction, addToStruct)
		} else if strings.HasPrefix(field, "args") {
			loadedValue, err := utils.LoadValue(field, args)
			if err != nil {
				logrus.Errorln("error loading value in matchDecrypt: ", err)
				return nil, err
			}
			stringValue, ok := loadedValue.(string)
			if !ok {
				return nil, fmt.Errorf("Value should be of type string and not %T", loadedValue)
			}
			decodedValue, err := base64.StdEncoding.DecodeString(stringValue)
			if err != nil {
				return nil, err
			}
			decrypted := make([]byte, len(decodedValue))
			err1 := decryptAESCFB(decrypted, decodedValue, m.aesKey, m.aesKey[:aes.BlockSize])
			if err1 != nil {
				logrus.Errorln("error decrypting value in matchDecrypt: ", err1)
				return nil, err1
			}
			er := utils.StoreValue(field, string(decrypted), args)
			if er != nil {
				logrus.Errorln("error storing value in matchDecrypt: ", er)
				return nil, er
			}
		} else {
			return nil, fmt.Errorf("invalid field (%s) provided", field)
		}
	}
	return actions, nil
}

func encryptAESCFB(dst, src, key, iv []byte) error {
	aesBlockEncrypter, err := aes.NewCipher([]byte(key))
	if err != nil {
		return err
	}
	aesEncrypter := cipher.NewCFBEncrypter(aesBlockEncrypter, iv)
	aesEncrypter.XORKeyStream(dst, src)
	return nil
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

func matchHash(rule *config.Rule, args map[string]interface{}) (*model.PostProcess, error) {
	actions := &model.PostProcess{}
	for _, field := range rule.Fields {
		if strings.HasPrefix(field, "res") {
			addToStruct := model.PostProcessAction{Action: "hash", Field: field}
			actions.PostProcessAction = append(actions.PostProcessAction, addToStruct)
		} else if strings.HasPrefix(field, "args") {
			loadedValue, err := utils.LoadValue(field, args)
			if err != nil {
				logrus.Errorln("error loading value in matchHash: ", err)
				return nil, err
			}
			stringValue, ok := loadedValue.(string)
			if !ok {
				return nil, fmt.Errorf("Value should be of type string and not %T", loadedValue)
			}
			h := sha256.New()
			_, _ = h.Write([]byte(stringValue))
			hashed := hex.EncodeToString(h.Sum(nil))
			er := utils.StoreValue(field, hashed, args)
			if er != nil {
				logrus.Errorln("error storing value in matchHash: ", er)
				return nil, er
			}
		} else {
			return nil, fmt.Errorf("invalid field (%s) provided", field)
		}
	}
	return actions, nil
}
