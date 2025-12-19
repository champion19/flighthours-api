#!/bin/bash

echo "🔍 DIAGNÓSTICO COMPLETO DE KEYCLOAK - REALM FLIGHTHOURS"
echo "========================================================"
echo ""

# Variables
KEYCLOAK_URL="http://localhost:8080"
REALM="flighthours"
CLIENT_ID="emma"
CLIENT_SECRET="M9HfWmIWf6huAnpKPXIGNdDeTfrwcNMt"
ADMIN_USER="admin"
ADMIN_PASS="1997"

echo "1️⃣  Verificando conectividad con Keycloak..."
curl -s -f "$KEYCLOAK_URL" > /dev/null
if [ $? -eq 0 ]; then
    echo "   ✅ Keycloak está en línea"
else
    echo "   ❌ Keycloak NO responde en $KEYCLOAK_URL"
    exit 1
fi
echo ""

echo "2️⃣  Verificando realm '$REALM'..."
REALM_INFO=$(curl -s "$KEYCLOAK_URL/realms/$REALM/.well-known/openid-configuration")
if [ -n "$REALM_INFO" ]; then
    echo "   ✅ Realm '$REALM' existe y está activo"
    echo "   📍 Token endpoint: $(echo $REALM_INFO | grep -o '"token_endpoint":"[^"]*"' | cut -d'"' -f4)"
else
    echo "   ❌ Realm '$REALM' no encontrado"
    exit 1
fi
echo ""

echo "3️⃣  Obteniendo token de admin desde realm MASTER..."
ADMIN_TOKEN=$(curl -s -X POST "$KEYCLOAK_URL/realms/master/protocol/openid-connect/token" \
  -H "Content-Type: application/x-www-form-urlencoded" \
  -d "username=$ADMIN_USER" \
  -d "password=$ADMIN_PASS" \
  -d "grant_type=password" \
  -d "client_id=admin-cli" | grep -o '"access_token":"[^"]*"' | cut -d'"' -f4)

if [ -n "$ADMIN_TOKEN" ]; then
    echo "   ✅ Token de admin obtenido exitosamente"
else
    echo "   ❌ Error obteniendo token de admin"
    echo "   Verifica credenciales: KEYCLOAK_ADMIN=$ADMIN_USER, KEYCLOAK_ADMIN_PASSWORD=$ADMIN_PASS"
    exit 1
fi
echo ""

echo "4️⃣  Verificando cliente '$CLIENT_ID' en realm '$REALM'..."
CLIENT_INFO=$(curl -s -X GET "$KEYCLOAK_URL/admin/realms/$REALM/clients" \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -H "Content-Type: application/json" | grep -o "\"clientId\":\"$CLIENT_ID\"")

if [ -n "$CLIENT_INFO" ]; then
    echo "   ✅ Cliente '$CLIENT_ID' existe en realm '$REALM'"
else
    echo "   ❌ Cliente '$CLIENT_ID' NO encontrado en realm '$REALM'"
fi
echo ""

echo "5️⃣  Verificando roles en realm '$REALM'..."
ROLES=$(curl -s -X GET "$KEYCLOAK_URL/admin/realms/$REALM/roles" \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -H "Content-Type: application/json")

if echo "$ROLES" | grep -q '"name":"pilot"'; then
    echo "   ✅ Rol 'pilot' existe"
else
    echo "   ⚠️  Rol 'pilot' NO encontrado"
    echo "   👉 Crear rol 'pilot' en Keycloak Admin Console"
fi

# Listar todos los roles
echo "   📋 Roles disponibles:"
echo "$ROLES" | grep -o '"name":"[^"]*"' | cut -d'"' -f4 | sed 's/^/      - /'
echo ""

echo "6️⃣  Verificando políticas de contraseña..."
REALM_CONFIG=$(curl -s -X GET "$KEYCLOAK_URL/admin/realms/$REALM" \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -H "Content-Type: application/json")

PASSWORD_POLICY=$(echo "$REALM_CONFIG" | grep -o '"passwordPolicy":"[^"]*"' | cut -d'"' -f4)
if [ -n "$PASSWORD_POLICY" ]; then
    echo "   ✅ Política de contraseña configurada:"
    echo "      $PASSWORD_POLICY"
else
    echo "   ⚠️  Sin política de contraseña (permitirá cualquier contraseña)"
fi
echo ""

echo "7️⃣  Verificando configuración de SMTP..."
SMTP_SERVER=$(echo "$REALM_CONFIG" | grep -o '"smtpServer":')
if [ -n "$SMTP_SERVER" ]; then
    echo "   ✅ SMTP configurado"
else
    echo "   ⚠️  SMTP no configurado (emails de verificación fallarán)"
fi
echo ""

echo "8️⃣  Verificando acciones requeridas..."
REQUIRED_ACTIONS=$(curl -s -X GET "$KEYCLOAK_URL/admin/realms/$REALM/authentication/required-actions" \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -H "Content-Type: application/json")

echo "   📋 Acciones requeridas activas:"
echo "$REQUIRED_ACTIONS" | grep -o '"alias":"[^"]*".*"enabled":true' | cut -d'"' -f4 | sed 's/^/      - /'
echo ""

echo "9️⃣  Verificando cliente '$CLIENT_ID' - Configuración detallada..."
CLIENT_ID_UUID=$(curl -s -X GET "$KEYCLOAK_URL/admin/realms/$REALM/clients" \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -H "Content-Type: application/json" | grep -B2 "\"clientId\":\"$CLIENT_ID\"" | grep -o '"id":"[^"]*"' | head -1 | cut -d'"' -f4)

if [ -n "$CLIENT_ID_UUID" ]; then
    CLIENT_DETAILS=$(curl -s -X GET "$KEYCLOAK_URL/admin/realms/$REALM/clients/$CLIENT_ID_UUID" \
      -H "Authorization: Bearer $ADMIN_TOKEN" \
      -H "Content-Type: application/json")

    echo "   📋 Configuración del cliente:"
    echo "      - ID: $CLIENT_ID_UUID"
    echo "      - Direct Access Grants: $(echo $CLIENT_DETAILS | grep -o '"directAccessGrantsEnabled":[^,]*' | cut -d':' -f2)"
    echo "      - Standard Flow: $(echo $CLIENT_DETAILS | grep -o '"standardFlowEnabled":[^,]*' | cut -d':' -f2)"
    echo "      - Service Accounts: $(echo $CLIENT_DETAILS | grep -o '"serviceAccountsEnabled":[^,]*' | cut -d':' -f2)"
fi
echo ""

echo "✅ DIAGNÓSTICO COMPLETADO"
echo "========================================================"
echo ""
echo "📝 RESUMEN:"
echo "   - Realm: $REALM"
echo "   - Cliente: $CLIENT_ID"
echo "   - Admin: $ADMIN_USER (desde realm master)"
echo ""
echo "💡 SIGUIENTE PASO:"
echo "   Ejecuta tu aplicación Go y registra un usuario pilot"
echo "   El código ya está corregido para usar realm 'master' en LoginAdmin"
