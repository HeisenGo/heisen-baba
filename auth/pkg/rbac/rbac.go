package rbac

import (
    "github.com/mikespook/gorbac"
)

var rbac = gorbac.New()

func InitRBAC() {
    // Define roles
    ownerRole := gorbac.NewStdRole("owner")
    employeeRole := gorbac.NewStdRole("employee")
    techTeamRole := gorbac.NewStdRole("tech_team")

    // Define permissions
    viewCompanyPermission := gorbac.NewStdPermission("view_company")
    editCompanyPermission := gorbac.NewStdPermission("edit_company")
    viewTripPermission := gorbac.NewStdPermission("view_trip")
    editTripPermission := gorbac.NewStdPermission("edit_trip")

    // Assign permissions to roles
    ownerRole.Assign(viewCompanyPermission)
    ownerRole.Assign(editCompanyPermission)
    ownerRole.Assign(viewTripPermission)
    ownerRole.Assign(editTripPermission)

    employeeRole.Assign(viewCompanyPermission)
    employeeRole.Assign(viewTripPermission)

    techTeamRole.Assign(viewCompanyPermission)
    techTeamRole.Assign(viewTripPermission)
    techTeamRole.Assign(editTripPermission)

    // Add roles to RBAC
    rbac.Add(ownerRole)
    rbac.Add(employeeRole)
    rbac.Add(techTeamRole)
}

func IsGranted(roles []string, permission string) bool {
    for _, role := range roles {
        if rbac.IsGranted(role, gorbac.NewStdPermission(permission), nil) {
            return true
        }
    }
    return false
}