#!/bin/bash

# Script de inicio rápido para el stack de observabilidad
# Autor: Equipo Flighthours
# Fecha: 2025-12-06

set -e

# Colores
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
CYAN='\033[0;36m'
NC='\033[0m'

print_header() {
    echo ""
    echo -e "${CYAN}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
    echo -e "${CYAN}  $1${NC}"
    echo -e "${CYAN}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
    echo ""
}

print_status() {
    case $1 in
        "success") echo -e "${GREEN}✅ $2${NC}" ;;
        "error") echo -e "${RED}❌ $2${NC}" ;;
        "warning") echo -e "${YELLOW}⚠️  $2${NC}" ;;
        "info") echo -e "${BLUE}ℹ️  $2${NC}" ;;
    esac
}

print_header "🚀 Flighthours API - Inicio Rápido de Observabilidad"

echo -e "${BLUE}Este script configurará automáticamente:${NC}"
echo "  • Prometheus (métricas)"
echo "  • Loki (agregación de logs)"
echo "  • Promtail (captura de logs)"
echo "  • Grafana (visualización)"
echo ""
read -p "¿Continuar? (s/n) " -n 1 -r
echo ""

if [[ ! $REPLY =~ ^[Ss]$ ]]; then
    echo "Cancelado."
    exit 0
fi

# Cambiar al directorio del proyecto
cd "$(dirname "$0")/.."

# 1. Detener servicios existentes (si hay)
print_header "1. Limpiando servicios anteriores"

if docker-compose -f docker-compose.grafana.yml ps | grep -q Up; then
    print_status "info" "Deteniendo servicios existentes..."
    docker-compose -f docker-compose.grafana.yml down
    print_status "success" "Servicios detenidos"
else
    print_status "info" "No hay servicios corriendo"
fi

# 2. Crear directorio de logs
print_header "2. Preparando directorio de logs"

LOG_DIR="/tmp/flighthours-logs"
if [ ! -d "$LOG_DIR" ]; then
    mkdir -p "$LOG_DIR"
    print_status "success" "Directorio de logs creado: $LOG_DIR"
else
    print_status "info" "Directorio de logs ya existe: $LOG_DIR"
fi

# 3. Iniciar stack de observabilidad
print_header "3. Iniciando stack de observabilidad"

print_status "info" "Iniciando Prometheus, Loki, Promtail y Grafana..."
docker-compose -f docker-compose.grafana.yml up -d

# Esperar a que los servicios estén listos
echo ""
print_status "info" "Esperando a que los servicios estén listos..."
sleep 5

# Verificar servicios
SERVICES=("flighthours-prometheus" "flighthours-loki" "flighthours-promtail" "flighthours-grafana")
ALL_UP=true

for service in "${SERVICES[@]}"; do
    if docker ps | grep -q "$service"; then
        print_status "success" "$service está corriendo"
    else
        print_status "error" "$service NO está corriendo"
        ALL_UP=false
    fi
done

if [ "$ALL_UP" = false ]; then
    print_status "error" "Algunos servicios no iniciaron correctamente"
    echo ""
    print_status "info" "Revisa los logs con:"
    echo "  docker-compose -f docker-compose.grafana.yml logs"
    exit 1
fi

# 4. Verificar conectividad
print_header "4. Verificando conectividad"

sleep 3

# Prometheus
if curl -s http://localhost:9090/-/healthy > /dev/null 2>&1; then
    print_status "success" "Prometheus respondiendo en http://localhost:9090"
else
    print_status "warning" "Prometheus aún no está listo (puede tomar unos segundos)"
fi

# Loki
if curl -s http://localhost:3100/ready > /dev/null 2>&1; then
    print_status "success" "Loki respondiendo en http://localhost:3100"
else
    print_status "warning" "Loki aún no está listo (puede tomar unos segundos)"
fi

# Grafana
if curl -s http://localhost:3000/api/health > /dev/null 2>&1; then
    print_status "success" "Grafana respondiendo en http://localhost:3000"
else
    print_status "warning" "Grafana aún no está listo (puede tomar unos segundos)"
fi

# 5. Instrucciones finales
print_header "✅ Stack de Observabilidad Iniciado"

echo -e "${GREEN}URLs de acceso:${NC}"
echo -e "  • Prometheus:  ${CYAN}http://localhost:9090${NC}"
echo -e "  • Loki:        ${CYAN}http://localhost:3100${NC}"
echo -e "  • Grafana:     ${CYAN}http://localhost:3000${NC} ${YELLOW}(admin/admin)${NC}"
echo ""
echo -e "${BLUE}Próximos pasos:${NC}"
echo ""
echo -e "${YELLOW}1. Iniciar el backend con logging:${NC}"
echo "   ./scripts/run-backend-with-logging.sh"
echo ""
echo -e "${YELLOW}2. O iniciar el backend de forma normal:${NC}"
echo "   go run cmd/main.go"
echo ""
echo -e "${YELLOW}3. Generar tráfico:${NC}"
echo "   curl http://localhost:8085/health"
echo ""
echo -e "${YELLOW}4. Abrir Grafana:${NC}"
echo "   open http://localhost:3000"
echo "   (Login: admin/admin)"
echo ""
echo -e "${YELLOW}5. Ver dashboards:${NC}"
echo "   • Flighthours API - Metrics Overview"
echo "   • Flighthours API - Logs"
echo ""
echo -e "${BLUE}Comandos útiles:${NC}"
echo ""
echo "• Verificar estado completo:"
echo "  ./scripts/verify-observability.sh"
echo ""
echo "• Ver logs de servicios:"
echo "  docker-compose -f docker-compose.grafana.yml logs -f"
echo ""
echo "• Detener servicios:"
echo "  docker-compose -f docker-compose.grafana.yml down"
echo ""
echo "• Reiniciar servicios:"
echo "  docker-compose -f docker-compose.grafana.yml restart"
echo ""

print_status "success" "¡Todo listo! 🎉"
