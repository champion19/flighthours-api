<!DOCTYPE html>
<html lang="es">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>FlightHours - Actualizar Contraseña</title>
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
        .icon {
            width: 80px; height: 80px;
            background: linear-gradient(135deg, #0047AB 0%, #4A90E2 100%);
            border-radius: 50%;
            display: flex; align-items: center; justify-content: center;
            margin: 0 auto 1.5rem;
            font-size: 40px; color: white;
            box-shadow: 0 4px 16px rgba(0, 71, 171, 0.3);
        }
        h1 { color: #2C3E50; font-size: 24px; font-weight: 600; margin-bottom: 1rem; }
        .message { color: #64748B; font-size: 15px; line-height: 1.6; margin-bottom: 1.5rem; }

        /* Form Styles */
        .form-group {
            margin-bottom: 1.25rem;
            text-align: left;
        }
        .form-group label {
            display: block;
            color: #374151;
            font-size: 14px;
            font-weight: 500;
            margin-bottom: 0.5rem;
        }
        .form-group input {
            width: 100%;
            padding: 12px 16px;
            border: 2px solid #E5E7EB;
            border-radius: 10px;
            font-size: 16px;
            transition: border-color 0.2s, box-shadow 0.2s;
        }
        .form-group input:focus {
            outline: none;
            border-color: #0047AB;
            box-shadow: 0 0 0 3px rgba(0, 71, 171, 0.1);
        }
        .password-requirements {
            background: #F3F4F6;
            border-radius: 8px;
            padding: 12px;
            margin-bottom: 1.5rem;
            text-align: left;
        }
        .password-requirements p {
            font-size: 13px;
            color: #6B7280;
            margin-bottom: 4px;
        }
        .password-requirements ul {
            margin: 0;
            padding-left: 20px;
        }
        .password-requirements li {
            font-size: 12px;
            color: #6B7280;
            margin-bottom: 2px;
        }
        .btn {
            display: block;
            width: 100%;
            padding: 14px 24px;
            font-size: 16px;
            font-weight: 600;
            color: white;
            background: linear-gradient(135deg, #0047AB 0%, #4A90E2 100%);
            border: none;
            border-radius: 10px;
            cursor: pointer;
            transition: transform 0.2s, box-shadow 0.2s;
        }
        .btn:hover {
            transform: translateY(-2px);
            box-shadow: 0 4px 16px rgba(0, 71, 171, 0.3);
        }
        .btn:disabled {
            background: #9CA3AF;
            cursor: not-allowed;
            transform: none;
        }
        .error-message {
            background: #FEE2E2;
            color: #DC2626;
            padding: 12px;
            border-radius: 8px;
            margin-bottom: 1rem;
            font-size: 14px;
            display: none;
        }
        .success-card {
            display: none;
        }
        .success-card .icon {
            background: linear-gradient(135deg, #10B981 0%, #059669 100%);
        }
        .success-message {
            color: #059669;
            font-weight: 600;
            font-size: 16px;
            margin-bottom: 1rem;
        }
        .footer { margin-top: 1.5rem; color: #64748B; font-size: 13px; }

        /* Loading state */
        .loading .btn-text { display: none; }
        .loading .btn-loader { display: inline-block; }
        .btn-loader {
            display: none;
            width: 20px;
            height: 20px;
            border: 3px solid rgba(255,255,255,0.3);
            border-top-color: white;
            border-radius: 50%;
            animation: spin 0.8s linear infinite;
        }
        @keyframes spin {
            to { transform: rotate(360deg); }
        }
    </style>
</head>
<body>
    <div class="container">
        <!-- Password Update Form -->
        <div class="card form-card" id="formCard">
            <div class="logo">FlightHours</div>
            <div class="icon">🔐</div>

            <h1>Actualizar Contraseña</h1>

            <p class="message">
                Ingresa tu nueva contraseña para completar la actualización de tu cuenta.
            </p>

            <div class="error-message" id="errorMessage"></div>

            <form id="updatePasswordForm">
                <div class="form-group">
                    <label for="newPassword">Nueva contraseña</label>
                    <input type="password" id="newPassword" name="new_password" required minlength="8" placeholder="Mínimo 8 caracteres">
                </div>

                <div class="form-group">
                    <label for="confirmPassword">Confirmar contraseña</label>
                    <input type="password" id="confirmPassword" name="confirm_password" required minlength="8" placeholder="Repite la contraseña">
                </div>

                <div class="password-requirements">
                    <p><strong>Requisitos de la contraseña:</strong></p>
                    <ul>
                        <li>Mínimo 8 caracteres</li>
                        <li>Se recomienda incluir mayúsculas, minúsculas y números</li>
                    </ul>
                </div>

                <button type="submit" class="btn" id="submitBtn">
                    <span class="btn-text">Actualizar Contraseña</span>
                    <span class="btn-loader"></span>
                </button>
            </form>

            <div class="footer">
                <p>© 2024 FlightHours</p>
            </div>
        </div>

        <!-- Success Card -->
        <div class="card success-card" id="successCard">
            <div class="logo">FlightHours</div>
            <div class="icon">✓</div>

            <h1>¡Contraseña Actualizada!</h1>

            <p class="success-message">
                Tu contraseña ha sido actualizada exitosamente.
            </p>

            <p class="message">
                Ya puedes iniciar sesión con tu nueva contraseña en la aplicación FlightHours.
            </p>

            <p class="message" style="color: #0047AB; font-weight: 600;">
                Puedes cerrar esta ventana.
            </p>

            <div class="footer">
                <p>© 2024 FlightHours</p>
            </div>
        </div>
    </div>

    <script>
        // Get the token from URL
        function getTokenFromUrl() {
            const urlParams = new URLSearchParams(window.location.search);
            // Keycloak includes the key parameter in the action URL
            return urlParams.get('key') || '';
        }

        // Show error message
        function showError(message) {
            const errorEl = document.getElementById('errorMessage');
            errorEl.textContent = message;
            errorEl.style.display = 'block';
        }

        // Hide error message
        function hideError() {
            document.getElementById('errorMessage').style.display = 'none';
        }

        // Set loading state
        function setLoading(isLoading) {
            const btn = document.getElementById('submitBtn');
            if (isLoading) {
                btn.classList.add('loading');
                btn.disabled = true;
            } else {
                btn.classList.remove('loading');
                btn.disabled = false;
            }
        }

        // Show success card
        function showSuccess() {
            document.getElementById('formCard').style.display = 'none';
            document.getElementById('successCard').style.display = 'block';
        }

        // Form submission
        document.getElementById('updatePasswordForm').addEventListener('submit', async function(e) {
            e.preventDefault();
            hideError();

            const newPassword = document.getElementById('newPassword').value;
            const confirmPassword = document.getElementById('confirmPassword').value;
            const token = getTokenFromUrl();

            // Validate passwords match
            if (newPassword !== confirmPassword) {
                showError('Las contraseñas no coinciden');
                return;
            }

            // Validate password length
            if (newPassword.length < 8) {
                showError('La contraseña debe tener al menos 8 caracteres');
                return;
            }

            if (!token) {
                showError('Token de autorización no encontrado. Por favor, usa el enlace del correo electrónico.');
                return;
            }

            setLoading(true);

            try {
                const response = await fetch('http://localhost:8081/flighthours/api/v1/auth/update-password', {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json',
                    },
                    body: JSON.stringify({
                        token: token,
                        new_password: newPassword,
                        confirm_password: confirmPassword
                    })
                });

                const data = await response.json();

                if (response.ok && data.success) {
                    showSuccess();
                } else {
                    const errorMsg = data.mensaje?.contenido || data.message || 'Error al actualizar la contraseña';
                    showError(errorMsg);
                }
            } catch (error) {
                console.error('Error:', error);
                showError('Error de conexión. Por favor, intenta de nuevo.');
            } finally {
                setLoading(false);
            }
        });
    </script>
</body>
</html>
