# Rediseño Aeronáutico - Themes de Keycloak

## 🎨 Nueva Paleta de Colores

### Colores Principales - Aviation Blue
- **Aviation Blue** (`#0047AB`): Azul cobalto profundo - color principal de marca
- **Sky Blue** (`#4A90E2`): Azul cielo brillante - acentos
- **Flight Blue** (`#1E88E5`): Azul vibrante - CTAs y botones

### Colores Neutrales - Professional Grays
- **White** (`#FFFFFF`): Blanco puro
- **Light Gray** (`#F5F7FA`): Gris muy claro - fondos
- **Medium Gray** (`#E1E8ED`): Gris medio - bordes y divisores
- **Text Gray** (`#64748B`): Gris para texto secundario
- **Dark Gray** (`#2C3E50`): Gris oscuro - texto principal

### Colores de Estado
- **Success Green** (`#10B981`): Verde para éxito
- **Warning Amber** (`#F59E0B`): Ámbar para advertencias
- **Error Red** (`#EF4444`): Rojo para errores
- **Info Blue** (`#3B82F6`): Azul para información

## 📝 Archivos Actualizados

### 1. Login Theme CSS (`login/resources/css/login.css`)

**Cambios principales:**
- ✅ Paleta de colores aeronáutica completa
- ✅ Gradiente de fondo: Aviation Blue → Sky Blue
- ✅ Patrón sutil de aviación en el fondo
- ✅ Tipografía Inter (más moderna y profesional)
- ✅ Sombras suaves con tonos azules
- ✅ Botones con gradiente azul
- ✅ Inputs con bordes redondeados (10px)
- ✅ Cards con bordes más redondeados (20px)
- ✅ Animación de fade-in al cargar
- ✅ Icono de avión (✈️) en el logo
- ✅ Soporte para modo de alto contraste
- ✅ Mejoras de accesibilidad

**Características destacadas:**
```css
/* Gradiente de fondo aeronáutico */
background: linear-gradient(135deg, #0047AB 0%, #4A90E2 100%);

/* Patrón sutil de aviación */
body::before {
    background-image: linear-gradient(30deg, rgba(255,255,255,0.03) 12%, ...);
    background-size: 80px 140px;
}

/* Logo con icono de avión */
.logo::before {
    content: '✈️';
}

/* Botones con gradiente */
background: linear-gradient(135deg, #0047AB 0%, #1E88E5 100%);
```

### 2. Login Error Template (`login/error.ftl`)

**Cambios:**
- ✅ Actualizado con paleta aeronáutica
- ✅ Gradiente de fondo azul profundo → azul cielo
- ✅ Logo con icono de avión
- ✅ Tipografía Inter
- ✅ Sombras con tonos azules
- ✅ App-box con colores aeronáuticos
- ✅ Textos en grises profesionales

### 3. Email Template Base (`email/html/template.ftl`)

**Cambios:**
- ✅ Paleta aeronáutica completa
- ✅ Fondo gris claro (#F5F7FA)
- ✅ Borde superior azul aviation (#0047AB)
- ✅ Logo con icono de avión
- ✅ Botones con gradiente azul
- ✅ Sombras suaves azules
- ✅ Info-box y warning-box actualizados
- ✅ Footer con colores profesionales
- ✅ Tipografía Inter

**Características del email:**
```css
/* Logo con avión */
.logo::before {
    content: '✈️ ';
}

/* Botón con gradiente y sombra */
.button {
    background: linear-gradient(135deg, #0047AB 0%, #1E88E5 100%);
    box-shadow: 0 4px 12px rgba(0, 71, 171, 0.2);
}

/* Hover effect */
.button:hover {
    background: linear-gradient(135deg, #003d96 0%, #1976D2 100%);
    box-shadow: 0 6px 16px rgba(0, 71, 171, 0.3);
}
```

## 🎯 Inspiración de Diseño

El diseño está inspirado en:
- ✈️ **Aerolíneas clásicas**: Confianza, cielo y precisión
- 🌐 **Diseño moderno**: Gradientes suaves, sombras sutiles
- 📱 **UI contemporánea**: Bordes redondeados, espaciado generoso
- 🎨 **Paleta profesional**: Azules profundos y grises neutros

## 📊 Comparación Antes/Después

### Antes
- Colores: Púrpura/Violeta (#667eea, #764ba2)
- Azul genérico (#007BFF)
- Tipografía: Arial
- Sombras negras
- Sin iconos

### Después
- Colores: Aviation Blue (#0047AB), Sky Blue (#4A90E2)
- Grises profesionales
- Tipografía: Inter (moderna)
- Sombras azules sutiles
- Icono de avión ✈️
- Patrón de fondo aeronáutico

## 🚀 Instalación

Para aplicar los cambios:

```bash
# 1. Asegúrate de que Docker Desktop esté corriendo
# 2. Ejecuta el script de instalación
bash tools/keycloak-themes/install-theme.sh

# 3. Espera 30 segundos a que Keycloak reinicie
# 4. Refresca el navegador con Cmd+Shift+R
```

## 📱 Responsive Design

Todos los themes son completamente responsive:
- ✅ Desktop (> 600px)
- ✅ Tablet (600px - 480px)
- ✅ Mobile (< 480px)

## ♿ Accesibilidad

- ✅ Contraste WCAG AA compliant
- ✅ Focus visible en todos los elementos interactivos
- ✅ Soporte para modo de alto contraste
- ✅ Tamaños de fuente legibles
- ✅ Espaciado adecuado para touch targets

## 🎨 Elementos Visuales

### Gradientes
- Fondo: Aviation Blue → Sky Blue
- Botones: Aviation Blue → Flight Blue
- Iconos de error: Error Red → Dark Red
- Iconos de éxito: Success Green → Dark Green

### Sombras
- Pequeña: `0 2px 8px rgba(0, 71, 171, 0.08)`
- Media: `0 4px 16px rgba(0, 71, 171, 0.12)`
- Grande: `0 8px 32px rgba(0, 71, 171, 0.16)`

### Bordes Redondeados
- Inputs: 10px
- Botones: 10px
- Cards: 20px (login), 16px (email)
- Alerts: 10px

## 🔄 Próximos Pasos

1. **Iniciar Docker Desktop**
2. **Ejecutar** `bash tools/keycloak-themes/install-theme.sh`
3. **Configurar en Keycloak**:
   - Realm Settings → Themes
   - Login theme: `flighthours`
   - Email theme: `flighthours`
   - Save
4. **Probar** los nuevos diseños

## 💡 Notas

- Los cambios son solo visuales (CSS y HTML)
- No se modificó ninguna funcionalidad
- Compatible con todas las versiones de Keycloak 26.x
- Los emails se verán bien en todos los clientes de correo modernos
- El diseño es profesional y transmite confianza

---

**Diseño creado**: 2025-12-19
**Tema**: Aviation Professional - Classic & Trustworthy
**Colores**: Blue & White/Gray
