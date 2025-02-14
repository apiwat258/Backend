document.addEventListener("DOMContentLoaded", function () {
    const token = localStorage.getItem("jwtToken");

    console.log("🔍 Checking Token in `auth.js`:", token); // Debug

    if (!token) {
        console.log("❌ No token found in localStorage.");
        return;
    }

    fetchUserData();
});

// ✅ ดึงข้อมูลผู้ใช้จาก API
async function fetchUserData() {
    const token = localStorage.getItem("jwtToken");

    if (!token) {
        console.log("❌ No token found. Skipping user data fetch.");
        return;
    }

    try {
        const response = await fetch("http://localhost:8080/api/v1/protected/route", {
            method: "GET",
            headers: { "Authorization": `Bearer ${token}` }
        });

        if (!response.ok) {
            console.error("❌ Invalid Token or Session Expired");
            alert("Session expired. Please login again.");
            localStorage.removeItem("jwtToken");
            sessionStorage.removeItem("user_email");
            sessionStorage.removeItem("user_id");
            return;
        }

        const data = await response.json();
        console.log("✅ User Data:", data);

        // ✅ บันทึกข้อมูลลง `localStorage` และ `sessionStorage`
        localStorage.setItem("user_email", data.email);
        localStorage.setItem("user_id", data.user_id);
        sessionStorage.setItem("user_email", data.email);
        sessionStorage.setItem("user_id", data.user_id);

        // ✅ อัปเดต UI บน `index.html`
        if (document.getElementById("userEmail")) {
            document.getElementById("userEmail").textContent = data.email;
        }
        if (document.getElementById("userId")) {
            document.getElementById("userId").textContent = data.user_id;
        }
        if (document.getElementById("accountInfo")) {
            document.getElementById("accountInfo").style.display = "block";
        }
    } catch (error) {
        console.error("❌ Error fetching user data:", error);
    }
}
