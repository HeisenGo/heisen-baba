#!/bin/sh

# Get an admin access token
ADMIN_TOKEN=$(curl -X POST "http://keycloak:8080/realms/master/protocol/openid-connect/token" \
  -H "Content-Type: application/x-www-form-urlencoded" \
  -d "client_id=admin-cli&grant_type=password&username=admin&password=admin" \
  | jq -r .access_token)

# Check if token was obtained
if [ -z "$ADMIN_TOKEN" ]; then
  echo "Failed to obtain admin token"
  exit 1
fi

# Create a new realm
curl -X POST "http://keycloak:8080/admin/realms" \
  -H "Authorization: Bearer ${ADMIN_TOKEN}" \
  -H "Content-Type: application/json" \
  -d '{
    "id": "my-realm",
    "realm": "my-realm",
    "enabled": true
  }'

# Create a new client
curl -X POST "http://keycloak:8080/admin/realms/my-realm/clients" \
  -H "Authorization: Bearer ${ADMIN_TOKEN}" \
  -H "Content-Type: application/json" \
  -d '{
    "clientId": "my-client-id",
    "enabled": true,
    "protocol": "openid-connect",
    "publicClient": false,
    "secret": "my-client-secret",
    "redirectUris": ["http://oauth2-proxy:4180/oauth2/callback"]
  }'

# Create a new role
curl -X POST "http://keycloak:8080/admin/realms/my-realm/roles" \
  -H "Authorization: Bearer ${ADMIN_TOKEN}" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "my-role"
  }'

echo "\nScript finished"