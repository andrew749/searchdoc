package data_models

/**
* Hold info.plist information from the docset.
 */
type DocsetPlist struct {
	Identifier           string `plist:"CFBundleIdentifier"`
	Name                 string `plist:"CFBundleName"`
	DocsetPlatformFamily string `plist:"DocSetPlatformFamily"`
	isDashDocset         bool   `plist:"isDashDocset"`
}
