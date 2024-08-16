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
        })
        .catch(error => console.error('Error creating user:', error));
    }

    function showUserDetails(userId) {
        storage.setItem('user_id', userId);
        fetch(`${apiAccounts}/${userId}`)
            .then(response => response.json())
            .then(user => {
                userDetailsTitle.textContent = `Detalles de ${user.nickname}`;
                return fetch(`${apiTransactions}/summary/${userId}`);
            })
            .then(response => response.json())
            .then(transactions => {
                displayTransactions(transactions);
                userDetailsModal.style.display = 'block';
            })
            .catch(error => console.error('Error loading user details:', error));
    }

    function displayTransactions(transactions) {
        //const groupedTransactions = groupTransactionsByYearAndMonth(transactions);
        let html = '';
        // for (const [year, months] of Object.entries(groupedTransactions)) {
        //     html += `<h3>${year}</h3>`;
        //     for (const [month, data] of Object.entries(months)) {
        //         html += `
        //             <h4>${month}</h4>
        //             <p>Total: $${data.total.toFixed(2)}</p>
        //             <p>Promedio: $${data.average.toFixed(2)}</p>
        //             <ul>
        //                 ${data.transactions.map(t => `<li>$${t.amount} (${new Date(t.date).toLocaleDateString()})</li>`).join('')}
        //             </ul>
        //         `;
        //     }
        // }
        transactionsList.innerHTML = html;
    }

    function groupTransactionsByYearAndMonth(transactions) {
        return transactions.reduce((acc, t) => {
            const date = new Date(t.date);
            const year = date.getFullYear();
            const month = date.toLocaleString('default', { month: 'long' });
            if (!acc[year]) acc[year] = {};
            if (!acc[year][month]) acc[year][month] = { total: 0, count: 0, transactions: [] };
            acc[year][month].total += t.amount;
            acc[year][month].count++;
            acc[year][month].average = acc[year][month].total / acc[year][month].count;
            acc[year][month].transactions.push(t);
            return acc;
        }, {});
    }

    function sendEmail() {
        const userId = userDetailsTitle.dataset.userId;
        fetch(`/api/accounts/${userId}/send-summary`, { method: 'POST' })
            .then(response => response.json())
            .then(data => alert('Resumen enviado por correo'))
            .catch(error => console.error('Error sending email:', error));
    }

    function uploadTransactionFile(event) {
        event.preventDefault();
        const file = fileInput.files[0];
        if (!file) {
            alert('Por favor, selecciona un archivo primero.');
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
            alert('Archivo procesado correctamente');
            showUserDetails(userId);
        })
        .catch(error => {
            console.error('Error uploading file:', error);
            alert('Error al procesar el archivo');
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


});
