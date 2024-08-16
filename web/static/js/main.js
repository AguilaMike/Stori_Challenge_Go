document.addEventListener('DOMContentLoaded', () => {
    // Elementos del DOM
    const userList = document.getElementById('users');
    const newUserBtn = document.getElementById('new-user-btn');
    const newUserModal = document.getElementById('new-user-modal');
    const newUserForm = document.getElementById('new-user-form');
    const userDetailsModal = document.getElementById('user-details-modal');
    const userDetailsTitle = document.getElementById('user-details-title');
    const transactionsList = document.getElementById('transactions-list');
    const sendEmailBtn = document.getElementById('send-email-btn');
    const uploadForm = document.getElementById('upload-form');
    const fileInput = document.getElementById('transaction-file');
    // url apis
    const apiAccounts = '/api/accounts';
    const apiTransactions = '/api/transactions';

    // Cargar usuarios al iniciar la página
    loadUsers();

    // Event listeners
    newUserBtn.addEventListener('click', () => newUserModal.style.display = 'block');
    newUserForm.addEventListener('submit', createUser);
    sendEmailBtn.addEventListener('click', sendEmail);
    uploadForm.addEventListener('submit', uploadTransactionFile);

    // Cerrar modales al hacer clic fuera de ellos
    window.addEventListener('click', (event) => {
        if (event.target === newUserModal) newUserModal.style.display = 'none';
        if (event.target === userDetailsModal) userDetailsModal.style.display = 'none';
    });

    function loadUsers() {
        fetch(apiAccounts)
            .then(response => response.json())
            .then(users => {
                userList.innerHTML = '';
                if (users && users.length > 0) {
                    users.forEach(user => {
                        const li = document.createElement('li');
                        li.textContent = `${user.nickname} (${user.email})`;
                        li.addEventListener('click', () => showUserDetails(user.id));
                        userList.appendChild(li);
                    });
                } else {
                    const li = document.createElement('li');
                    li.textContent = 'No hay usuarios registrados';
                    userList.appendChild(li);
                }
            })
            .catch(error => console.error('Error loading users:', error));
    }

    function createUser(event) {
        event.preventDefault();
        const formData = new FormData(newUserForm);
        fetch(apiAccounts, {
            method: 'POST',
            body: JSON.stringify(Object.fromEntries(formData)),
            headers: { 'Content-Type': 'application/json' }
        })
        .then(response => response.json())
        .then(() => {
            loadUsers();
            newUserModal.style.display = 'none';
            newUserForm.reset();
            showNotification('Account creada exitosamente', 'success');
        })
        .catch(error => {
            console.error('Error creating user:', error);
            showNotification('Error al crear la cuenta', 'error');
        });
    }

    function showUserDetails(userId) {
        storage.setItem('user_id', userId);
        fetch(`${apiAccounts}/${userId}`)
            .then(response => response.json())
            .then(user => {
                userDetailsTitle.textContent = `Detalles de ${user.nickname}`;
                document.getElementById('send-email-btn').setAttribute('data-user-id', userId);
                return fetch(`${apiTransactions}/summary/${userId}`);
            })
            .then(response => response.json())
            .then(transaction => {
                displayTransactions(transaction);
                userDetailsModal.style.display = 'block';
            })
            .catch(error => console.error('Error loading user details:', error));
    }

    function displayTransactions(transaction) {
        let html = displayHeader(transaction);

        // Obtenemos las claves (yyyy-MM) y las ordenamos
        const sortedKeys = Object.keys(transaction.monthly).sort((a, b) => new Date(b) - new Date(a));

        let detalleHTML = '';
        sortedKeys.forEach(key => {
            const data = transaction.monthly[key];
            detalleHTML += `
                <div class="detail-section-item">
                    <h3>${key} <button class="toggle-details">Mostrar</button></h3>
                    <div class="transaction-details" style="display: none;">
                        <p>Balance: $${data.balance.toFixed(2)}</p>
                        <p>Operaciones: $${data.total_transactions}</p>
                        <p>Promedio Credito: $${data.average_credit.toFixed(2)}</p>
                        <p>Promedio Debito: $${data.average_debit.toFixed(2)}</p>
                        <ul>
                            ${data.transactions.map(t => `<li>$${t.amount} (${new Date(t.input_date).toLocaleDateString()})</li>`).join('')}
                        </ul>
                    </div>
                </div>
            `;
        });
        if (detalleHTML) {
            detalleHTML = `<div class="details-section">${detalleHTML}</div>`;
            html = html.replace('##DETALLES##', detalleHTML);
        } else {
            html = html.replace('##DETALLES##', '');
        }

        transactionsList.innerHTML = html;

        // Add event listeners for toggling details
        document.querySelectorAll('.toggle-details').forEach(button => {
            button.addEventListener('click', function() {
                const details = this.closest('.detail-section-item').querySelector('.transaction-details');
                details.style.display = details.style.display === 'none' ? 'block' : 'none';
                this.textContent = details.style.display === 'none' ? 'Mostrar' : 'Ocultar';
            });
        });

        // Add event listener for toggling credit/debit section
        const creditDebitToggle = document.getElementById('credit-debit-toggle');
        const creditDebitSection = document.querySelector('.credit-debit-section');
        creditDebitToggle.addEventListener('click', function() {
            creditDebitSection.style.display = creditDebitSection.style.display === 'none' ? 'grid' : 'none';
            this.textContent = creditDebitSection.style.display === 'none' ? 'Mostrar Detalles' : 'Ocultar Detalles';
        });
    }

    function displayHeader(transaction) {
        let html = '';
        if (transaction.summary) {
            html += `
            <div class="container-transaction">
                <div class="summary-section">
                    <h3>Resumen Total</h3>
                    <p>Balance: $${transaction.summary.total_balance.toFixed(2)}</p>
                    <p>Operaciones: $${transaction.summary.total_count}</p>
                    <div class="button-container">
                        <button id="credit-debit-toggle">Mostrar Detalles</button>
                    </div>
                </div>
                <div class="credit-debit-section" style="display: none;">
                    <div class="detail-section-item">
                        <h4>Credito</h4>
                        <p>Total: $${transaction.summary.total_credit.toFixed(2)}</p>
                        <p>Promedio: $${transaction.summary.average_credit.toFixed(2)}</p>
                        <p>Operaciones: $${transaction.summary.credit_count}</p>
                    </div>
                    <div class="detail-section-item">
                        <h4>Debito</h4>
                        <p>Total: $${transaction.summary.total_debit.toFixed(2)}</p>
                        <p>Promedio: $${transaction.summary.average_debit.toFixed(2)}</p>
                        <p>Operaciones: $${transaction.summary.debit_count}</p>
                    </div>
                </div>
                ##DETALLES##
            </div>
            `;
        }

        return html;
    }

    function sendEmail() {
        const userId = this.getAttribute('data-user-id');
        fetch(`${apiTransactions}/send-sumamry/${userId}`, { method: 'POST' })
            .then(response => response.json())
            .then(data => showNotification('Resumen enviado por correo', 'success'))
            .catch(error => {
                console.error('Error sending email:', error);
                showNotification('Error al enviar correo', 'error');
            });
    }

    function uploadTransactionFile(event) {
        event.preventDefault();
        const file = fileInput.files[0];
        if (!file) {
            showNotification('Por favor, selecciona un archivo primero.', 'warning');
            return;
        }

        const formData = new FormData();
        formData.append('transactionFile', file);
        formData.append('userID', storage.getItem('user_id'));

        fetch('/upload', {
            method: 'POST',
            body: formData
        })
        .then(response => response.text())
        .then(result => {
            console.log('Upload successful:', result);
            showNotification('Archivo procesado correctamente', 'success');
            showUserDetails(userId);
        })
        .catch(error => {
            console.error('Error uploading file:', error);
            showNotification('Error al procesar el archivo', 'error');
        });
    }

    // LocalStorage
    const storage = {
        setItem: (key, value) => localStorage.setItem(key, JSON.stringify(value)),
        getItem: key => JSON.parse(localStorage.getItem(key))
    };

    // Función para generar un ID único
    function generateUUID() {
        return 'xxxxxxxx-xxxx-4xxx-yxxx-xxxxxxxxxxxx'.replace(/[xy]/g, function(c) {
            var r = Math.random() * 16 | 0, v = c == 'x' ? r : (r & 0x3 | 0x8);
            return v.toString(16);
        });
    }

    // WebSockets
    // Obtener o generar usuario para webSocket
    function getOrCreateUserId() {
        let userId = localStorage.getItem('WSuserId');
        if (!userId) {
            userId = generateUUID();
            localStorage.setItem('WSuserId', userId);
        }
        return userId;
    }

    // Conexión al servidor de WebSockets
    const ws = new WebSocket(`ws://${window.location.host}/ws/?userID=${getOrCreateUserId()}`);

    socket.onmessage = function(event) {
        const data = JSON.parse(event.data);
        if (data.type === 'transaction_update') {
            updateTransactionUI(data.summary);
        }
    };

    socket.onclose = function(event) {
        console.log('WebSocket connection closed:', event);
    };

    function updateTransactionUI(summary) {

    }

    // Función para mostrar notificaciones
    function showNotification(message, type = 'info') {
        const container = document.getElementById('notification-container');
        const notification = document.createElement('div');
        notification.className = `notification ${type}`;
        notification.textContent = message;

        container.appendChild(notification);

        // Animar la entrada de la notificación
        setTimeout(() => {
            notification.style.transform = 'translateX(0)';
            notification.style.opacity = '1';
        }, 100);

        // Remover la notificación después de 5 segundos
        setTimeout(() => {
            notification.style.transform = 'translateX(100%)';
            notification.style.opacity = '0';
            setTimeout(() => {
                container.removeChild(notification);
            }, 300);
        }, 5000);
    }
});
