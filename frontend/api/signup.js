document.getElementById('signup-form').addEventListener('submit', async function (event) {
    event.preventDefault();

    const email = document.getElementById('signup-email').value;
    const password = document.getElementById('signup-password').value;
    const confirmPassword = document.getElementById('signup-confirm-password').value;

    if (password !== confirmPassword) {
        alert('Passwords do not match!');
        return;
    }

    try {
        const response = await fetch('http://127.0.0.1:8080/api/v1/auth/register', { // ✅ อัปเดต Endpoint
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ email, password }),
        });

        const data = await response.json();
        
        if (response.ok) {
            // ✅ บันทึกข้อมูลจาก Backend
            if (data.user_id) {
                localStorage.setItem("user_id", data.user_id);
            }
            if (data.token) { // ✅ หากมี token ให้เก็บไว้
                localStorage.setItem("auth_token", data.token);
            }
            localStorage.setItem('user_email', email); 

            alert("Signup successful! Redirecting to Role Selection...");
            window.location.href = 'Role.html';
        } else {
            alert(data.message || 'Registration failed!');
        }
    } catch (error) {
        console.error('Error:', error);
        alert('An error occurred. Please try again later.');
    }
});
