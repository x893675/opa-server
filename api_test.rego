#opa test -v api.rego api_test.rego

package api.rbac

import data.api.rbac

# Mock Data
roles = {
	"alice": ["admin"],
	"bob": ["regular"],
}

permissions = {"regular": [
	{
		"verbs": ["GET", "UPDATE"],
		"apiGroups": ["*"],
		"resources": ["namespaces", "clusters"],
		"resourceNames": [],
		"nonResourceURLs": ["/metrics"],
	},
	{
		"verbs": "POST",
		"apiGroups": ["core.io"],
		"resources": ["*"],
		"resourceNames": [],
		"nonResourceURLs": [],
	},
]}

test_admin_allowed {
	allow with input as {"user": "alice"} with rbac.roles as roles
}

test_non_admin_not_allowed {
	not allow with input as {"user": "bob"} with rbac.roles as roles
}

test_grants_allowed {
	allow with input as {"user": "bob", "resourceRequest": true, "verb": "UPDATE", "apiGroup": "*", "resource": "namespaces"} with rbac.roles as roles with rbac.permissions as permissions
}

test_grants_nonResourcesURLs_allowed {
    allow with input as {"user": "bob", "resourceRequest": false,"path":"/metrics"} with rbac.roles as roles with rbac.permissions as permissions
}

test_grants_nonResourcesURLs_not_allowed {
    allow with input as {"user": "bob", "resourceRequest": false,"path":"/healthz"} with rbac.roles as roles with rbac.permissions as permissions
}

test_who_are_list {
    who_are["alice"] with input as {"role": "admin"} with rbac.roles as roles
    who_are["bob"] with input as {"role": "regular"} with rbac.roles as roles
}
