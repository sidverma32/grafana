+++
title = "Permissions"
description = "Understand fine-grained access control permissions"
keywords = ["grafana", "fine-grained access-control", "roles", "permissions", "enterprise"]
weight = 115
+++

# Permissions 

Each permission is defined by an action and a scope. When evaluating a fine-grained access decision, consider what specific action a user will be allowed to perform and on what resource.

To grant permissions to a user, create built-in role assignments. For more information, refer to [Built-in role assignments]({{< relref "./roles.md#built-in-role-assignments" >}}).

To learn more about what permissions are used for what resources, refer to [Resources with fine-grained permissions]({{< relref "./_index.md#resources-with-fine-grained-permissions" >}}).

### Action

When evaluating a decision for access, a user is allowed to perform a specific _action_ on a resource if they have a permission with the relevant action assigned to it.

### Scope

_scope_ describes where an action can be performed. For example, you can grant a user access to read a specific user profile, by associating a permission with scope `users:<userId>` to the relevant role.
You can also combine multiple scopes if necessary, by using `/` as a delimiter. 

## Available permissions 

Action | Applicable scopes | Description
--- | --- | ---
roles:list | roles:* | Allows to list available roles without permissions. 
roles:read | roles:* | Allows to read a specific role with it's permissions.
roles:write | permissions:delegate | Allows to create or update a custom role.
roles:delete | permissions:delegate | Allows to delete a custom role.
roles.builtin:list | roles:* | Allows to list built-in role assignments.
roles.builtin:add | permissions:delegate | Allows to create a built-in role assignment.
roles.builtin:remove | permissions:delegate | Allows to delete a built-in role assignment.
reports.admin:write | reports:* | Allows to create or update reports.
reports:delete | reports:* | Allows to delete reports.
reports:read | reports:* | Allows to list all available reports and to get a specific report. 
reports:send | reports:* | Allows to send report email.
reports.settings:write | n/a | Allows to update report settings.
reports.settings:read | n/a | Allows to read report settings.
provisioning:reload | service:access-control | Allows to reload provisioning files after an update.
users:read | global:users:* | Allows to read, search user profiles. 
users:write | global:users:* | Allows to update user profiles.
users.teams:read | global:users:* | Allows to read user teams.
users.authtoken:list | global:users:* | Allows to list auth tokens assigned to users.
users.authtoken:update | global:users:* | Allows to update auth tokens assigned to users.
users.password:update | global:users:* | Allows to update users password.
users:delete | global:users:* | Allows to delete users.
users:create | n/a | Allows to create users.
users:enable | global:users:* | Allows to enable users.
users:disable | global:users:* | Allows to disable users.
users.permissions:update | global:users:* | Allows to update users org level permissions.
users:logout | global:users:* | Allows to enforce logout for users.
users.quotas:list | global:users:* | Allows to list user quotas.
users.quotas:update | global:users:* | Allows to update user quotas.
org.users.read | users:* | Allows to get user profiles within the organization.
org.users.add | users:* | Allows to add users to the organization.
org.users.remove | users:* | Allows to remove users from the organization.
org.users.role:update | users:* | Allows to update users organization role for the assigned organization.
ldap.user:read | n/a | Allows to read LDAP users.
ldap.user:sync | n/a | Allows to sync LDAP users.
ldap.status:read | n/a | Allows to check LDAP status.

## Scope definitions

Scope | Description
--- | --- 
roles:* | Indicates against what roles an action can be performed. For example, `roles:*` assumes any roles, and `roles:randomuid` assumes only a role with UID `randomuid`.  
permissions:delegate | The scope is only applicable for roles associated with the Access Control itself and indicates that you can delegate your permissions only, or a subset of it, by creating a new role or making an assignment. 
reports:* | Indicates against what reports an action can be performed.
service:access-control | Only relevant for provisioning and indicates that the action can be performed only for access control provisioning files.
global:users:* | Indicates that action can be performed against users globally.
users:* | Indicates that an action can be performed against users in organization level.