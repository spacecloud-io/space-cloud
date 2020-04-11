package main

// func getplugin(latestversion string) (*cobra.Command, error) {
// 	mod := fmt.Sprintf("%s/cmd_%s.so", getSpaceCLIDirectory(), latestversion)
// 	plug, err := plugin.Open(mod)
// 	if err != nil {
// 		return nil, err
// 	}
// 	commands, err := plug.Lookup("GetRootCommand")
// 	if err != nil {
// 		return nil, err
// 	}
//
// 	rootCmd := commands.(func() *cobra.Command)()
// 	return rootCmd, nil
//
// }
