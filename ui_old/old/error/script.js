const urlCur = window.location.href;
const apiUrl = new URL("login", urlCur);
const storageKey = 'token';
console.log(apiUrl);
async function doLogin(username, password) {
    const data = {
        username: username,
        password: password
    };
    try {
        const response = await fetch(apiUrl, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify(data)
        });
        const responseData = await response.json();

        if (responseData.status === 'success') {
            const token = responseData.token;
            console.log('Token:', token);

            // Armazena o token no localStorage
            localStorage.setItem(storageKey, token);

            // Redireciona para a página protegida (substitua pela URL da sua página)
            // window.location.href = 'index.html';
            return true;
        } else {
            alert('Erro no login:', responseData.message);
            
        }
    } catch (error) {
        console.error('Erro:', error);
        alert('Erro ao fazer login. Tente novamente.');
        return false;
    }
}

const loginForm = document.getElementById('loginForm');
const menuItems = document.querySelectorAll('header nav ul li a');

loginForm.addEventListener('submit', (event) => {
    // event.preventDefault();

    const username = document.getElementById('username').value;
    const password = document.getElementById('password').value;

    // Simula a validação de login (substitua com sua lógica real)

    if (doLogin(username, password)) {
        // Habilita os itens do menu
        menuItems.forEach(item => {
            item.style.pointerEvents = 'auto';
            item.style.opacity = 1;
        });

        // Esconde o formulário de login
        loginForm.style.display = 'none';
    } else {
        alert('username or password invalid!');
    }
});

