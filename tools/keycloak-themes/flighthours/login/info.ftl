<!DOCTYPE html>
<html lang="es">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Flighthours - Verificando correo</title>
    <style>
        * { margin: 0; padding: 0; box-sizing: border-box; }
        body {
            font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, Oxygen, Ubuntu, Cantarell, sans-serif;
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            min-height: 100vh;
            display: flex;
            align-items: center;
            justify-content: center;
            padding: 20px;
        }
        .container {

            background: white;
            max-width: 450px;
            width: 100%;
            padding: 3rem 2rem;
            border-radius: 16px;
            box-shadow: 0 20px 60px rgba(0, 0, 0, 0.3);
            text-align: center;
        }
        .logo {
            color: #007BFF;
            font-size: 32px;
            font-weight: 700;
            margin-bottom: 2rem;
        }
        .spinner {
            width: 60px;
            height: 60px;
            border: 5px solid #f3f3f3;
            border-top: 5px solid #007BFF;
            border-radius: 50%;
            animation: spin 0.8s linear infinite;
            margin: 0 auto 2rem;
        }
        .spinner.hidden { display: none; }
        @keyframes spin {
            0% { transform: rotate(0deg); }
            100% { transform: rotate(360deg); }
        }
        h1 {
            color: #1a1a1a;
            font-size: 24px;
            margin-bottom: 1rem;
            font-weight: 600;
        }
        .message {
            color: #666;
            font-size: 15px;
            line-height: 1.6;
            margin-bottom: 2rem;
        }
        .success-icon {
            width: 80px;
            height: 80px;
            background: #28a745;
            border-radius: 50%;
            display: none;
            align-items: center;
            justify-content: center;
            margin: 0 auto 2rem;
        }
        .success-icon.show { display: flex; }
        .success-icon svg { width: 40px; height: 40px; fill: white; }
        .error-icon {
            width: 80px;
            height: 80px;
            background: #dc3545;
            border-radius: 50%;
            display: none;
            align-items: center;
            justify-content: center;
            margin: 0 auto 2rem;
        }
        .error-icon.show { display: flex; }
        .error-icon svg { width: 40px; height: 40px; fill: white; }
        .app-box {
            background: linear-gradient(135deg, #f0f9ff 0%, #e0f2fe 100%);
            border: 2px solid #007BFF;
            border-radius: 12px;
            padding: 1.5rem;
            margin-bottom: 1.5rem;
        }
        .app-box h2 { color: #007BFF; font-size: 18px; margin-bottom: 0.5rem; }
        .app-box p { color: #0369a1; font-size: 14px; margin: 0; }
        .app-icon { font-size: 48px; margin-bottom: 0.5rem; }
        .btn {
            display: inline-block;
            padding: 14px 32px;
            background: linear-gradient(135deg, #007BFF 0%, #0056b3 100%);
            color: white;
            text-decoration: none;
            border-radius: 8px;
            font-weight: 600;
            font-size: 16px;
            border: none;
            cursor: pointer;
            transition: transform 0.2s, box-shadow 0.2s;
            box-shadow: 0 4px 12px rgba(0, 123, 255, 0.3);
        }
        .btn:hover {
            transform: translateY(-2px);
            box-shadow: 0 6px 20px rgba(0, 123, 255, 0.4);
        }
        .footer {
            margin-top: 2rem;
            color: #999;
            font-size: 13px;
        }
        .hidden { display: none !important; }
    </style>
</head>
<body>
    <div class="container">
        <div class="logo">Flighthours</div>

        <!-- Estado: Cargando -->
        <div id="loading-state">
            <div class="spinner"></div>
            <h1>Verificando tu correo...</h1>
            <p class="message">
                Por favor espera un momento. Estamos procesando tu verificación.
            </p>
        </div>

        <!-- Estado: Éxito -->
        <div id="success-state" class="hidden">
            <div class="success-icon show">
                <svg viewBox="0 0 24 24"><path d="M9 16.17L4.83 12l-1.42 1.41L9 19 21 7l-1.41-1.41z"/></svg>
            </div>
            <h1 style="color: #28a745;">¡Correo verificado!</h1>
            <p class="message">
                Tu cuenta ha sido verificada correctamente.
            </p>

            <div class="app-box">
                <div class="app-icon">📱</div>
                <h2>Abre la aplicación MotoGo</h2>
                <p>Ya puedes iniciar sesión con tu correo y contraseña.</p>
            </div>

            <p class="message" style="color: #10b981; font-weight: 600; margin-bottom: 0;">
                Puedes cerrar esta ventana.
            </p>
        </div>

        <!-- Estado: Error -->
        <div id="error-state" class="hidden">
            <div class="error-icon show">
                <svg viewBox="0 0 24 24"><path d="M19 6.41L17.59 5 12 10.59 6.41 5 5 6.41 10.59 12 5 17.59 6.41 19 12 13.41 17.59 19 19 17.59 13.41 12z"/></svg>
            </div>
            <h1 style="color: #dc3545;">Error de verificación</h1>
            <p class="message" id="error-message">
                No pudimos verificar tu correo. El enlace puede haber expirado.
            </p>
            <a href="flighthours://resend-verification" class="btn">Reenviar correo</a>
        </div>

        <!-- Estado: Ya verificado -->
        <div id="already-verified-state" class="hidden">
            <div class="success-icon show">
                <svg viewBox="0 0 24 24"><path d="M9 16.17L4.83 12l-1.42 1.41L9 19 21 7l-1.41-1.41z"/></svg>
            </div>
            <h1 style="color: #17a2b8;">Correo ya verificado</h1>
            <p class="message">
                Tu correo electrónico ya estaba verificado.<br>
                Puedes usar la aplicación normalmente.
            </p>
            <a href="flighthours://home" class="btn">Abrir Flighthours</a>
        </div>

        <div class="footer">
            © ${.now?string('yyyy')} Flighthours
        </div>
    </div>

    <script>
        (function() {
            // Configuración - URL del backend
            const BACKEND_URL = 'http://localhost:8081/flighthours/api/v1/auth/verify-email';

            // Elementos del DOM
            const loadingState = document.getElementById('loading-state');
            const successState = document.getElementById('success-state');
            const errorState = document.getElementById('error-state');
            const alreadyVerifiedState = document.getElementById('already-verified-state');
            const errorMessage = document.getElementById('error-message');

            function showState(state) {
                loadingState.classList.add('hidden');
                successState.classList.add('hidden');
                errorState.classList.add('hidden');
                alreadyVerifiedState.classList.add('hidden');
                state.classList.remove('hidden');
            }

            function getTokenFromUrl() {
                // Obtener el parámetro 'key' de la URL actual
                const urlParams = new URLSearchParams(window.location.search);
                return urlParams.get('key');
            }

            async function verifyEmail() {
                const token = getTokenFromUrl();

                if (!token) {
                    // Si no hay token, mostrar éxito (Keycloak ya procesó)
                    console.log('[FlightHours] No hay token en URL, Keycloak ya procesó la acción');
                    showState(successState);
                    return;
                }

                console.log('[FlightHours] Token encontrado, enviando al backend...');
                console.log('[FlightHours] URL del backend:', BACKEND_URL);

                try {
                    const response = await fetch(BACKEND_URL, {
                        method: 'POST',
                        headers: {
                            'Content-Type': 'application/json',
                        },
                        body: JSON.stringify({ token: token })
                    });

                    const data = await response.json();
                    console.log('[FlightHours] Respuesta del backend:', data);

                    if (response.ok && data.success) {
                        showState(successState);
                    } else if (data.code === 'MOD_KC_EMAIL_ALREADY_VERIFIED_WARN_00001') {
                        showState(alreadyVerifiedState);
                    } else {
                        // Mostrar mensaje de error específico si está disponible
                        if (data.message && data.message.contenido) {
                            errorMessage.textContent = data.message.contenido;
                        } else if (data.message) {
                            errorMessage.textContent = data.message;
                        }
                        showState(errorState);
                    }
                } catch (error) {
                    console.error('[FlightHours] Error de conexión:', error);
                    // Mostrar error de conexión pero NO el éxito
                    errorMessage.textContent = 'Error de conexión con el servidor. Por favor, intenta de nuevo más tarde.';
                    showState(errorState);
                }
            }

            // Ejecutar verificación cuando la página cargue
            if (document.readyState === 'loading') {
                document.addEventListener('DOMContentLoaded', verifyEmail);
            } else {
                verifyEmail();
            }
        })();
    </script>
</body>
</html>
