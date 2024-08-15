document.addEventListener('DOMContentLoaded', () => {
    const userList = document.getElementById('users');
    const newUserBtn = document.getElementById('new-user-btn');
    const newUserModal = document.getElementById('new-user-modal');
    const newUserForm = document.getElementById('new-user-form');
    const userDetailsModal = document.getElementById('user-details-modal');
    const userDetailsTitle = document.getElementById('user-details-title');
    const transactionsList = document.getElementById('transactions-list');
    const sendEmailBtn = document.getElementById('send-email-btn');
    const uploadForm = document.getElementById('upload-form');

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
        fetch('/users')
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
        fetch('/users/create', {
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
        fetch(`/users/detail/${userId}`)
            .then(response => response.json())
            .then(user => {
                userDetailsTitle.textContent = `Detalles de ${user.nickname}`;
                return fetch(`/api/transactions?account_id=${userId}`);
            })
            .then(response => response.json())
            .then(transactions => {
                displayTransactions(transactions);
                userDetailsModal.style.display = 'block';
            })
            .catch(error => console.error('Error loading user details:', error));
    }

    function displayTransactions(transactions) {
        const groupedTransactions = groupTransactionsByYearAndMonth(transactions);
        let html = '';
        for (const [year, months] of Object.entries(groupedTransactions)) {
            html += `<h3>${year}</h3>`;
            for (const [month, data] of Object.entries(months)) {
                html += `
                    <h4>${month}</h4>
                    <p>Total: $${data.total.toFixed(2)}</p>
                    <p>Promedio: $${data.average.toFixed(2)}</p>
                    <ul>
                        ${data.transactions.map(t => `<li>$${t.amount} (${new Date(t.date).toLocaleDateString()})</li>`).join('')}
                    </ul>
                `;
            }
        }
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
        const userId = userDetailsTitle.dataset.userId;
        const formData = new FormData(uploadForm);
        formData.append('user_id', userId);
        fetch('/upload', {
            method: 'POST',
            body: formData
        })
        .then(response => response.json())
        .then(data => {
            alert('Archivo procesado correctamente');
            showUserDetails(userId);
        })
        .catch(error => console.error('Error uploading file:', error));
    }
});
