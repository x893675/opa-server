package api.rbac

# import roles list from data.api.rbac
import data.api.rbac.permissions
import data.api.rbac.roles
import input

# By default, deny requests.
default allow = false

# More than one OR Condition for a variable `allow`

# Allow admins to do anything.
allow {
	user_is_admin
}

# Allow the action if the user is granted permission to perform the action.
allow {
	input.resourceRequest == true

	# check for blacklisted grants for user
	# not user_grant_is_blacklisted

	# Find grants for the user.
	some grant
	user_is_granted[grant]

	# Check if the grant permits the action. And Condition
	is_verb_match(grant.verbs)
	is_apiGroup_match(grant.apiGroups)
	is_resource_match(grant.resources)
	is_resourceName_match(grant.resourceNames)
}

# non resource request match
allow {
	input.resourceRequest == false
}

# user_is_admin is true if...
user_is_admin {
	# for some `i`...
	some i

	# "admin" is the `i`-th element in the user->role mappings for the identified user.
	roles[input.user][i] == "admin"
}

# user_is_granted is a set of grants for the user identified in the request.
# The `grant` will be contained if the set `user_is_granted` for every...
user_is_granted[grant] {
	some i, j

	# `role` assigned an element of the user_roles for this user...
	role := roles[input.user][i]

	# `grant` assigned a single grant from the grants list for 'role'...
	grant := permissions[role][j]
}

# who_are is a set of users who has roles identified in the request.
who_are[user] {
    # for some `user`...
	some user

	# `roleList` assigned a single user roles list...
	roleList := roles[user]

    # Check if the roleList matches the input role
    roleList[_] == input.role
}

# who_all_are is a set of users along with roles list who has roles identified in the request.
who_all_are[{ user: roleList }] {
    # for some `user` and `role`...
	some user, role

	# `roleList` assigned a single user roles list...
	roleList := roles[user]

    # Check if the roleList matches the input role list
    roleList[_] == input.roles[role]
}


# roles_can is a set of roles who has grants identified in the request.
# roles_can[role] {
#     # for some `role` and `permission`...
# 	some role, permission

# 	# `roleList` assigned a single permission role list...
# 	roleList := permissions[role]

#     # `grant` assigned a single grant from all grants list...
#     grant := roleList[permission]

#     # Check if the grant permits the action.
# 	input.grant.action == grant.action
# 	input.grant.type == grant.type
# }

# users_can is a set of roles who has grants identified in the request.
# users_can[user] {
#     # for some `role` and `permission`...
# 	some role, permission

# 	# `roleList` assigned a single permission role list...
# 	roleList := permissions[role]

#     # `grant` assigned a single grant from all grants list...
#     grant := roleList[permission]

#     # Check if the grant permits the action.
# 	input.grant.action == grant.action
# 	input.grant.type == grant.type

# 	 # for some `user`
# 	some user

# 	# `askedRoleList` assigned a single user roles list...
# 	askedRoleList := roles[user]

#     # Check if the askedRoleList matches the selected role from upper iteration
#     askedRoleList[_] == role

# }

is_verb_match(verb) {
	some i
	verb[i] == "*"
}

is_verb_match(verb) {
	some i
	verb[i] == input.verb
}

is_apiGroup_match(apiGroups) {
	some i
	apiGroups[i] == "*"
}

is_apiGroup_match(apiGroups) {
	some i
	apiGroups[i] == input.apiGroup
}

is_resource_match(resources) {
	some i
	resources[i] == "*"
}

is_resource_match(resources) {
	some i
	resources[i] == input.resource
}

is_resourceName_match(resourceNames) {
	count(resourceNames) == 0
}

is_resourceName_match(resourceNames) {
	some i
	resourceNames[i] == input.resourceNames
}

is_nonResourceURL_match(nonResourceURLs) {
	some i
	nonResourceURLs[i] == "*"
}

is_nonResourceURL_match(nonResourceURLs) {
	some i
	nonResourceURLs[i] == input.Path
}

