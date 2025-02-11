document.addEventListener("DOMContentLoaded", function () {
    const emailInput = document.getElementById("email");

    if (emailInput) {
        const email = localStorage.getItem("user_email");
        if (email) {
            emailInput.value = email;
        } else {
            alert("User email not found. Please register first.");
            window.location.href = "Sign Up.html";
        }
    } else {
        console.error("Error: Element with ID 'email' not found in DOM.");
    }

    // ✅ เมื่อกด Submit จะส่งทั้งข้อมูลฟาร์ม + CID ที่อัปโหลดไว้ไปยัง Backend
    const farmerForm = document.getElementById("farmer-form");
    if (farmerForm) {
        farmerForm.addEventListener("submit", async function (event) {
            event.preventDefault();

            const userID = localStorage.getItem("user_id");
            if (!userID) {
                alert("User ID not found. Please register first.");
                return;
            }

            const certificationCID = localStorage.getItem("certification_cid"); // ✅ ดึง `CID` ที่อัปโหลดไว้ก่อนหน้า
            console.log("✅ Certification CID from localStorage:", certificationCID);

            const farmerData = {
                userid: userID,
                company_name: document.getElementById("company_name").value,
                firstname: document.getElementById("firstname").value,
                lastname: document.getElementById("lastname").value,
                email: document.getElementById("email").value,
                address: document.getElementById("address").value,
                address2: document.getElementById("address2").value,
                areacode: document.getElementById("areacode").value,
                phone: document.getElementById("phone").value,
                post: document.getElementById("post").value,
                city: document.getElementById("city").value, 
                location_link: document.getElementById("location_link").value
            };

            try {
                const response = await fetch("http://127.0.0.1:8080/api/v1/farmers", { // ✅ เปลี่ยน API Path
                    method: "POST",
                    headers: { "Content-Type": "application/json" },
                    body: JSON.stringify(farmerData)
                });

                const textResponse = await response.text(); // ✅ รับ response เป็น string
                console.log("Server Response:", textResponse); // ✅ Log response

                if (!response.ok) {
                    alert("Error: " + textResponse);
                    return;
                }

                const data = JSON.parse(textResponse); // ✅ แปลงเป็น JSON ถ้าเป็นไปได้
                alert("Farmer information saved successfully!");

                // ✅ หลังจากบันทึกฟาร์มแล้ว → บันทึก Certification CID ลงตาราง organiccertification
                if (certificationCID) {
                    console.log("📌 Sending certification data to backend...");
                    const certResponse = await fetch("http://127.0.0.1:8080/api/v1/certifications/create", { // ✅ เปลี่ยน API Path
                        method: "POST",
                        headers: { "Content-Type": "application/json" },
                        body: JSON.stringify({
                            farmerid: data.farmer_id, // ✅ ใช้ farmer_id ที่เพิ่งบันทึก
                            certificationtype: "Organic ACT",
                            certificationcid: certificationCID,
                            issued_date: "2025-02-07"
                        })
                    });

                    const certData = await certResponse.json();
                    console.log("✅ Certification Response:", certData);

                    if (certResponse.ok) {
                        alert("Certification saved successfully!");
                    } else {
                        console.error("❌ Failed to save certification:", certData.error);
                        alert("Failed to save certification: " + certData.error);
                    }
                }

                window.location.href = "index.html";

            } catch (error) {
                console.error("Error:", error);
                alert("An error occurred while saving data.");
            }
        });
    } else {
        console.error("Error: Form with ID 'farmer-form' not found in DOM.");
    }
});
