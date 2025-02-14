document.addEventListener("DOMContentLoaded", function () {
    const loginForm = document.getElementById("loginForm");

    if (loginForm) {
        loginForm.addEventListener("submit", async function (event) {
            event.preventDefault();

            const email = document.getElementById("member-login-number").value;
            const password = document.getElementById("member-login-password").value;

            const loginData = { email: email, password: password };

            try {
                const response = await fetch("http://localhost:8080/api/v1/auth/login", {
                    method: "POST",
                    headers: { "Content-Type": "application/json" },
                    body: JSON.stringify(loginData)
                });

                const data = await response.json();

                if (response.ok) {
                    alert("✅ Login Successful!");
                    
                    // ✅ บันทึก Token ลง localStorage
                    localStorage.setItem("jwtToken", data.token);
                    console.log("🔍 Stored Token:", localStorage.getItem("jwtToken")); // Debug

                    // ✅ บันทึกอีเมล และ ID ของผู้ใช้
                    localStorage.setItem("user_email", data.email);
                    localStorage.setItem("user_id", data.user_id);

                    // ✅ ไปที่หน้าเลือก Role หรือ Index
                    window.location.href = "Role.html";
                } else {
                    alert("❌ Login Failed: " + (data.error || "Invalid credentials"));
                }
            } catch (error) {
                console.error("❌ Error:", error);
                alert("An error occurred while logging in.");
            }
        });
    }
});
