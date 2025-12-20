<!DOCTYPE html>
<html lang="es">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Flighthours - Redirigiendo...</title>
    <style>
        * { margin: 0; padding: 0; box-sizing: border-box; }
        body {
            font-family: 'Inter', -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
            background: linear-gradient(135deg, #0047AB 0%, #4A90E2 100%);
            background-attachment: fixed;
            min-height: 100vh;
            display: flex;
            align-items: center;
            justify-content: center;
            padding: 20px;
        }
        .container { max-width: 480px; width: 100%; }
        .card {
            background: white;
            padding: 3rem 2.5rem;
            border-radius: 20px;
            box-shadow: 0 8px 32px rgba(0, 71, 171, 0.16);
            text-align: center;
        }
        .logo {
            color: #0047AB;
            font-size: 32px;
            font-weight: 700;
            margin-bottom: 2rem;
            display: flex;
            align-items: center;
            justify-content: center;
            gap: 0.5rem;
        }
        .logo::before { content: '✈️'; font-size: 28px; }
        .spinner {
            width: 50px;
            height: 50px;
            border: 4px solid #E1E8ED;
            border-top: 4px solid #0047AB;
            border-radius: 50%;
            animation: spin 1s linear infinite;
            margin: 0 auto 1.5rem;
        }
        @keyframes spin {
            0% { transform: rotate(0deg); }
            100% { transform: rotate(360deg); }
        }
        h1 { color: #2C3E50; font-size: 22px; font-weight: 600; margin-bottom: 1rem; }
        .message { color: #64748B; font-size: 15px; line-height: 1.6; margin-bottom: 1.5rem; }
        .button {
            display: inline-block;
            padding: 14px 28px;
            background: linear-gradient(135deg, #0047AB 0%, #1E88E5 100%);
            color: white;
            text-decoration: none;
            border: none;
            border-radius: 10px;
            font-weight: 600;
            cursor: pointer;
            transition: all 0.3s ease;
            box-shadow: 0 4px 12px rgba(0, 71, 171, 0.3);
        }
        .button:hover {
            background: linear-gradient(135deg, #003d96 0%, #1976D2 100%);
            transform: translateY(-2px);
            box-shadow: 0 6px 16px rgba(0, 71, 171, 0.4);
        }
        .footer { margin-top: 1.5rem; color: #64748B; font-size: 13px; }
    </style>
</head>
<body>
    <div class="container">
        <div class="card">
            <div class="logo">Flighthours</div>
            <div class="spinner"></div>

            <h1>Procesando tu verificación...</h1>

            <p class="message">
                Estamos verificando tu correo electrónico. Serás redirigido automáticamente.
            </p>

            <!-- Formulario que se auto-submits para ejecutar la acción -->
            <form id="kc-action-form" action="${url.loginAction}" method="post">
                <noscript>
                    <p class="message">
                        <button type="submit" class="button">Haz clic aquí para continuar</button>
                    </p>
                </noscript>
            </form>

            <div class="footer">
                <p>© ${.now?string('yyyy')} Flighthours</p>
            </div>
        </div>
    </div>

    <script>
        // Auto-submit el formulario inmediatamente cuando la página carga
        document.addEventListener('DOMContentLoaded', function() {
            document.getElementById('kc-action-form').submit();
        });
    </script>
</body>
</html>
