package auth

//
// func (m *Module) setPublicKey(pemData string) error {
// 	m.lock.Lock()
// 	defer m.lock.Unlock()
//
// 	block, _ := pem.Decode([]byte(pemData))
// 	if block == nil {
// 		return errors.New("failed to parse PEM block containing the key")
// 	}
//
// 	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
// 	if err != nil {
// 		return err
// 	}
//
// 	key, ok := pub.(*rsa.PublicKey)
//
// 	if !ok {
// 		return errors.New("key type is not a RSA public key")
// 	}
//
// 	// Set the public key
// 	m.config.PublicKey = key
// 	helpers.Logger.LogInfoln("Public key of runner server set successfully")
// 	return nil
// }
//
// // We need to retrieve the public key used by the runner server instance. This needs to be done on a periodic
// // basis since the server may generate new pair of public private keys. Let's call this once a week
// func (m *Module) routineGetPublicKey() {
// 	ticker := time.NewTicker(168 * time.Hour)
// 	for range ticker.C {
// 		m.fetchPublicKey()
// 	}
// }
//
// func (m *Module) fetchPublicKey() (success bool) {
// 	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
// 	defer cancel()
//
// 	req, err := http.NewRequestWithContext(ctx, "GET", "http://api.spaceuptech.com/v1/runner/runner/public-key", nil)
// 	if err != nil {
// 		helpers.Logger.LogError(helpers.GetRequestID(ctx),"Could not fetch runner public key - %s", err.Error())
// 		return false
// 	}
//
// 	res, err := http.DefaultClient.Do(req)
// 	if err != nil {
// 		helpers.Logger.LogError(helpers.GetRequestID(ctx),"Could not fetch runner public key - %s", err.Error())
// 		return false
// 	}
//
// 	publicKey := new(model.PublicKeyPayload)
// 	if err := json.NewDecoder(res.Body).Decode(publicKey); err != nil {
// 		helpers.Logger.LogError(helpers.GetRequestID(ctx),"Could not decode runner public key payload - %s", err.Error())
// 		return false
// 	}
//
// 	if err := m.setPublicKey(publicKey.PemData); err != nil {
// 		helpers.Logger.LogError(helpers.GetRequestID(ctx),"Could not parse runner public key - %s", err.Error())
// 		return false
// 	}
//
// 	return true
// }
