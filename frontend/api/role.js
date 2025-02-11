document.addEventListener("DOMContentLoaded", function () {
    const email = localStorage.getItem("user_email");
    const userID = localStorage.getItem("user_id");

    if (!email || !userID) {
        alert("User data missing. Please register first.");
        window.location.href = "Sign Up.html";  // 🔹 ถ้าไม่มี User ID หรือ Email ให้กลับไปสมัครใหม่
    }

    document.getElementById("userEmail").innerText = email; // 🔹 แสดง Email บนหน้า Role
});

function selectRole(role) {
    const email = localStorage.getItem("user_email");
    const token = localStorage.getItem("auth_token"); // ✅ ใช้ Token (ถ้ามี)

    fetch("http://127.0.0.1:8080/api/v1/auth/update-role", { // ✅ อัปเดต API Path
        method: "POST",
        headers: { 
            "Content-Type": "application/json",
            "Authorization": token ? `Bearer ${token}` : "" // ✅ ใส่ Token (ถ้ามี)
        },
        body: JSON.stringify({ email: email, role: role })
    })
    .then(response => response.json())
    .then(data => {
        if (data.error) {
            alert("Error: " + data.error);
        } else {
            localStorage.setItem("user_role", role);  // 🔹 บันทึก Role ที่เลือก
            alert("Role updated successfully!");
            window.location.href = role + ".html"; // 🔹 เปลี่ยนไปยังหน้าของ Role นั้น ๆ เช่น Farmer.html
        }
    })
    .catch(error => {
        console.error("Error:", error);
    });
}
