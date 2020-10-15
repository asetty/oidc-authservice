package groups

// XXX
// Maybe the interface can be generalized to take care of getting
// the user ID as well. It could be named something like ProviderFlavor
// or UserProvider. This way getting userid + groups (and maybe additional
// functionality) could be encapsulated into one package.
//
// TODO see how this fits in with refactoring to separate authenticators
// into their own package. Is this groups stuff specific to session authenticator?

type GroupsMethod interface {
	GetGroups(claims map[string]interface{}) ([]string, error)
}
